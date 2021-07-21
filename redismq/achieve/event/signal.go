package event

import (
	"os"
	"os/signal"
	"syscall"
)

var closing bool

func InitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			closing = true
			return
		case syscall.SIGHUP:
			reload()
		default:
			return
		}
	}
}

func reload() {
}
