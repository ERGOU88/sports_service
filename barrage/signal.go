package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// InitSignal register signals handler.
func InitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			fmt.Printf("start stop service\n")
			StopWebsocket()
			fmt.Printf("start stop service success\n")
			return
		default:
			return
		}
	}
}
