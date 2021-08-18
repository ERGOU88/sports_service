package event

import (
	"os"
	"os/signal"
	"syscall"
)

var closing bool

func InitSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-sigChan
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
