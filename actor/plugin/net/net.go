package net

import (
	"log"
	bnet "gamelib/base/net"
	"gamelib/base/util"
	snet "net"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// Acceptor ------------------------------------------------------------------
type AcceptorEvent struct {
	C *snet.TCPConn
	E error
}

type Disconnected struct {
}

type AcceptorAddrEvent struct {
	A *snet.TCPAddr
}

func (ae *AcceptorEvent) GetPort() int {
	return ae.C.LocalAddr().(*snet.TCPAddr).Port
}

type Acceptor struct {
	addr  string
	lower int
	upper int

	l       *bnet.Acceptor
	connNum int
}

func (a *Acceptor) OnStart(ctx actor.Context) {
	l, e := bnet.NewAcceptor(
		a.addr,
		func(c *snet.TCPConn, e error) {
			ctx.Self().Tell(&AcceptorEvent{C: c, E: e})
		})
	util.PanicOnErr(e)
	a.l = l
	ctx.Self().Tell(&AcceptorAddrEvent{l.Addr()})
}
func (a *Acceptor) OnOtherMessage(ctx actor.Context, msg interface{}) {
	switch m := msg.(type) {
	case *actor.Stopping, *actor.Restart:
		if a.l != nil {
			a.l.Close()
			a.l = nil
		}
	case *actor.Terminated:
		a.connNum--
		log.Println("conn==", a.connNum)
		if a.connNum <= a.lower && a.l == nil {
			log.Println("conn== start accept", a.connNum, a.lower, a.upper)
			l, e := bnet.NewAcceptor(
				a.addr,
				func(c *snet.TCPConn, e error) {

					ctx.Self().Tell(&AcceptorEvent{C: c, E: e})
				})
			util.PanicOnErr(e)
			a.l = l
			ctx.Self().Tell(&AcceptorAddrEvent{l.Addr()})
		}
	case *AcceptorEvent:
		if m.C != nil {
			a.connNum++
			log.Println("conn==", a.connNum)
			if a.connNum >= a.upper && a.l != nil {
				log.Println("conn== stop accept", a.connNum, a.lower, a.upper)
				a.l.Close()
				a.l = nil
			}
		}
	}
}

func NewAcceptor(addr string, upper, lower int) *Acceptor {
	a := &Acceptor{
		addr:  addr,
		lower: lower,
		upper: upper,
	}
	return a
}

// Conn -----------------------------------------------------------------
type ConnectionEvent struct {
	E error
}

type op int

const (
	send      op = 1
	closesend op = 2
	close     op = 3
)

type ConnOp struct {
	op
	m bnet.Message
}

func SendMsg(ctx actor.Context, m bnet.Message) {
	ctx.Self().Tell(&ConnOp{op: send, m: m})
}

func CloseSend(ctx actor.Context) {
	ctx.Self().Tell(&ConnOp{op: closesend})
}

func Close(ctx actor.Context) {
	ctx.Self().Tell(&ConnOp{op: close})
}

type Connection struct {
	bc *bnet.Connection

	c           *snet.TCPConn
	io          bnet.MessageReadWriter
	readTimeOut int
}

func (c *Connection) OnStart(ctx actor.Context) {
	c.bc = bnet.NewConnection(
		c.c,
		func(m bnet.Message, e error) {
			if m != nil {
				ctx.Self().Tell(m)
			}
			if e != nil {
				ctx.Self().Tell(&ConnectionEvent{E: e})
			}
		},
		c.io,
		&bnet.ConnectionConfig{ReadTimeoutS: c.readTimeOut},
	)
}
func (c *Connection) OnOtherMessage(ctx actor.Context, m interface{}) {
	switch m := ctx.Message().(type) {
	case *actor.Stopping, *actor.Restart:
		c.bc.Close()
	case *ConnOp:
		switch m.op {
		case send:
			c.bc.Send(m.m)
		case closesend:
			c.bc.CloseSend()
		case close:
			c.bc.Close()
		}
	}
}

func NewConnection(c *snet.TCPConn, io bnet.MessageReadWriter, logsend, logrecv bool, readTimeOut int) *Connection {
	conn := &Connection{
		c:           c,
		io:          io,
		readTimeOut: readTimeOut,
	}
	return conn
}

func MakeAcceptor2(l *snet.TCPListener, upper, lower int) func(actor.ActorFunc) actor.ActorFunc {
	return func(next actor.ActorFunc) actor.ActorFunc {
		var (
			addr = l.Addr().String()
			a    *bnet.Acceptor
			e    error

			connNum int
		)
		return func(ctx actor.Context) {
			switch m := ctx.Message().(type) {
			case *actor.Started:
				a, e = bnet.NewAcceptor2(
					l,
					func(c *snet.TCPConn, e error) {

						ctx.Self().Tell(&AcceptorEvent{C: c, E: e})
					})
				util.PanicOnErr(e)
				ctx.Self().Tell(&AcceptorAddrEvent{a.Addr()})
			case *actor.Stopping, *actor.Restart:
				if a != nil {
					a.Close()
					a = nil
				}
			case *actor.Terminated:
				connNum--
				log.Println("conn==", connNum)
				if connNum <= lower && a == nil {
					log.Println("conn== start accept", connNum, lower, upper)
					a, e = bnet.NewAcceptor(
						addr,
						func(c *snet.TCPConn, e error) {

							ctx.Self().Tell(&AcceptorEvent{C: c, E: e})
						})
					util.PanicOnErr(e)
					ctx.Self().Tell(&AcceptorAddrEvent{a.Addr()})
				}
			case *AcceptorEvent:
				if m.C != nil {
					connNum++
					log.Println("conn==", connNum)
					if connNum >= upper && a != nil {
						log.Println("conn== stop accept", connNum, lower, upper)
						a.Close()
						a = nil
					}
				}
			}
			next(ctx)
		}
	}
}

// warning
func MakeAcceptor(addr string, upper, lower int) func(actor.ActorFunc) actor.ActorFunc {
	return func(next actor.ActorFunc) actor.ActorFunc {
		var (
			a *bnet.Acceptor
			e error

			connNum int
		)
		return func(ctx actor.Context) {
			switch m := ctx.Message().(type) {
			case *actor.Started:
				a, e = bnet.NewAcceptor(
					addr,
					func(c *snet.TCPConn, e error) {

						ctx.Self().Tell(&AcceptorEvent{C: c, E: e})
					})
				util.PanicOnErr(e)
				ctx.Self().Tell(&AcceptorAddrEvent{a.Addr()})
			case *actor.Stopping, *actor.Restart:
				if a != nil {
					a.Close()
					a = nil
				}
			case *actor.Terminated:
				connNum--
				log.Println("conn==", connNum)
				if connNum <= lower && a == nil {
					log.Println("conn== start accept", connNum, lower, upper)
					a, e = bnet.NewAcceptor(
						addr,
						func(c *snet.TCPConn, e error) {

							ctx.Self().Tell(&AcceptorEvent{C: c, E: e})
						})
					util.PanicOnErr(e)
					ctx.Self().Tell(&AcceptorAddrEvent{a.Addr()})
				}
			case *AcceptorEvent:
				if m.C != nil {
					connNum++
					log.Println("conn==", connNum)
					if connNum >= upper && a != nil {
						log.Println("conn== stop accept", connNum, lower, upper)
						a.Close()
						a = nil
					}
				}
			}
			next(ctx)
		}
	}
}

func MakeConnection(c *snet.TCPConn, io bnet.MessageReadWriter, logsend, logrecv bool, readTimeOut int) func(actor.ActorFunc) actor.ActorFunc {
	return func(next actor.ActorFunc) actor.ActorFunc {
		var (
			conn *bnet.Connection
		)
		return func(ctx actor.Context) {
			switch m := ctx.Message().(type) {
			case *actor.Started:
				conn = bnet.NewConnection(
					c,
					func(m bnet.Message, e error) {
						if m != nil {
							ctx.Self().Tell(m)
						}
						if e != nil {
							ctx.Self().Tell(&ConnectionEvent{E: e})
						}
					},
					io,
					&bnet.ConnectionConfig{ReadTimeoutS: readTimeOut},
				)

			case *actor.Stopping, *actor.Restart:
				conn.Close()
			case *ConnOp:
				switch m.op {
				case send:
					conn.Send(m.m)
				case closesend:
					conn.CloseSend()
				case close:
					conn.Close()
				}
			}
			next(ctx)
		}
	}
}
