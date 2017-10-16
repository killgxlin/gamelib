package net

import (
	"log"
	"net"
	"sync"
)

type AcceptEvent func(*net.TCPConn, error)

type Acceptor struct {
	l *net.TCPListener

	fe AcceptEvent

	w sync.WaitGroup
}

func (a *Acceptor) Addr() *net.TCPAddr {
	return a.l.Addr().(*net.TCPAddr)
}

func NewAcceptor2(l *net.TCPListener, fe AcceptEvent) (*Acceptor, error) {
	a := &Acceptor{fe: fe}

	a.l = l
	a.w.Add(1)
	go func() {
		defer a.w.Done()
		for {
			c, e := a.l.AcceptTCP()
			log.Println("acceptok", c)
			if e != nil {
				if e, ok := e.(net.Error); ok && e.Temporary() {
					continue
				}
				a.fe(c, e)
				return
			}

			a.fe(c, nil)
		}
	}()

	return a, nil
}

func NewAcceptor(addr string, fe AcceptEvent) (*Acceptor, error) {
	a := &Acceptor{fe: fe}

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal(e)
	}

	a.l = l.(*net.TCPListener)
	a.w.Add(1)
	go func() {
		defer a.w.Done()
		for {
			c, e := a.l.AcceptTCP()
			log.Println("acceptok", c)
			if e != nil {
				if e, ok := e.(net.Error); ok && e.Temporary() {
					continue
				}
				a.fe(c, e)
				return
			}

			a.fe(c, nil)
		}
	}()

	return a, nil
}

func (a *Acceptor) Close() {
	a.l.Close()
	a.w.Wait()
}
