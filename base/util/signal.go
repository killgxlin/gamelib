package util

import (
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func WaitSigInt() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c
}

func WaitSigHup() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP)
	<-c
}

func WaitSigUsr1() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1)
	<-c
}
