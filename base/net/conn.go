package net

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	LogSend = false
	LogRecv = false
)

type Message interface{}

type EventFunc func(Message, error)

type Connection struct {
	c          *net.TCPConn
	sch        chan Message
	sendClosed int32

	fe EventFunc

	w    sync.WaitGroup
	conf *ConnectionConfig
}

type ConnectionConfig struct {
	ReadTimeoutS   int
	SendChanLen    int
	MaxRecvMsgLen  int
	SendChanWaitMS int
}

func defaultConnectionConfig(conf *ConnectionConfig) *ConnectionConfig {
	if conf == nil {
		conf = &ConnectionConfig{}
	}

	if conf.ReadTimeoutS == 0 {
		conf.ReadTimeoutS = 30
	}
	if conf.SendChanLen == 0 {
		conf.SendChanLen = 1000
	}
	if conf.MaxRecvMsgLen <= 0 {
		conf.MaxRecvMsgLen = 1024 * 2
	}
	if conf.SendChanWaitMS == 0 {
		conf.SendChanWaitMS = 10
	}
	return conf
}

/*

client:
	send
	closesend
	recv
	recv eof
	close

server:
	recv
	recv eof
	send
	close

*/

type MessageReadWriter interface {
	ReadMsgWithLimit(io.Reader, int) (Message, error)
	WriteMsg(io.Writer, Message) error
}

func tryRecv(ch chan Message) (ms []Message, closed bool) {
	for {
		select {
		case m := <-ch:
			if m != nil {
				ms = append(ms, m)
			} else {
				return ms, true
			}
		default:
			return ms, false
		}
	}
}

func NewConnection(c1 *net.TCPConn, fe EventFunc, io MessageReadWriter, conf *ConnectionConfig) *Connection {
	conf = defaultConnectionConfig(conf)
	c := &Connection{c: c1, sch: make(chan Message, conf.SendChanLen), fe: fe, conf: conf}
	c.w.Add(2)

	go func() {
		defer func() {
			c.w.Done()
		}()
		for {
			if c.conf.ReadTimeoutS > 0 {
				c.c.SetDeadline(time.Now().Add(time.Second * time.Duration(c.conf.ReadTimeoutS)))
			}
			m, e := io.ReadMsgWithLimit(c.c, c.conf.MaxRecvMsgLen)
			if e != nil {
				fe(m, e)
				return
			}
			if LogRecv {
				log.Println(c.c, "recv", m)
			}
			fe(m, e)
		}
	}()

	go func() {
		defer func() {
			c.w.Done()
		}()

		w := bufio.NewWriter(c.c)
		for {
			ms, closed := tryRecv(c.sch)
			if len(ms) == 0 && c.conf.SendChanWaitMS > 0 {
				time.Sleep(time.Millisecond * time.Duration(c.conf.SendChanWaitMS))
			}

			for _, m := range ms {
				io.WriteMsg(w, m)
			}
			e := w.Flush()
			if e != nil {
				atomic.StoreInt32(&c.sendClosed, 1)
				c.fe(nil, e)
				return
			}
			if closed {
				c.c.CloseWrite()
				return
			}
		}

	}()

	return c
}

func (c *Connection) Send(m Message) {
	if atomic.LoadInt32(&c.sendClosed) == 0 {
		select {
		case c.sch <- m:
			if LogSend {
				log.Println(c.c, "sent", m)
			}
		default:
		}
	}
}

func (c *Connection) CloseSend() {
	if atomic.CompareAndSwapInt32(&c.sendClosed, 0, 1) {
		c.sch <- nil
	}
}

func (c *Connection) Close() {
	c.CloseSend()
	c.c.Close()
	c.w.Wait()
}
