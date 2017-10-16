package util

import (
	"errors"
	"fmt"
	"log"
	stdnet "net"
	"strconv"
)

func Listen(string, ip string, start, end int) (*stdnet.TCPListener, error) {
	for i := start; i < end; i++ {
		l, e := stdnet.Listen("tcp", ip+":"+strconv.Itoa(i))
		if e != nil {
			if e, ok := e.(*stdnet.OpError); ok {
				log.Printf("trace listen port:%v err:%v", i, e)
				continue
			}
		}
		log.Printf("trace listen port:%v", i)
		return l.(*stdnet.TCPListener), nil
	}
	return nil, errors.New("has not port")
}

func FindFreeAddr(net string, ip string, start, end int) (string, error) {
	l, e := Listen(net, ip, start, end)
	if e != nil {
		return "", e
	}
	addr := l.Addr().String()
	e = l.Close()
	if e != nil {
		return "", e
	}
	return addr, nil
}
func FindLanAddr(net string, start, end int) (string, error) {
	ip := GetLanIp()
	return FindFreeAddr(net, ip, start, end)
}

func GetLanAddr(port int) string {
	return fmt.Sprintf("%v:%v", GetLanIp(), port)
}

func GetLanIp() string {
	as, e := stdnet.InterfaceAddrs()
	if e != nil {
		log.Fatal(e)
	}
	for _, a := range as {
		ta := a.(*stdnet.IPNet)
		if ta == nil || len(ta.IP.To4()) != stdnet.IPv4len || ta.IP.IsLoopback() {
			continue
		}
		return ta.IP.String()
	}
	return "0.0.0.0"
}
