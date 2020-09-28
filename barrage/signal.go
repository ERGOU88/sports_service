package main

import (
	"fmt"
	"os"
	"os/signal"
	"sports_service/server/tools/nsq"
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
			// 停止websocket服务（清理用户链接）
			StopWebsocket()
			fmt.Printf("start stop service success\n")
			// 停止消费
			nsq.Stop()
			return
		default:
			return
		}
	}
}
