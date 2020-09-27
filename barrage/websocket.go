package main

import (
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"os"
)

const (
	READ_BUF_SIZE  = 2048
	WRITE_BUF_SIZE = 2045
	READ_DEAD_LINE = 120
)

const (
	MAX_MSG_PER_CONN = 1024
)

var running = false

var upgrader = websocket.Upgrader{
	ReadBufferSize: READ_BUF_SIZE,
	WriteBufferSize: WRITE_BUF_SIZE,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartWebsocket(bindAddr string) {
	r := mux.NewRouter()
	r.HandleFunc("/ws", WebSocketHandler)
	http.Handle("/", r)

	running = true

	go func() {
		if err := http.ListenAndServe(bindAddr, nil); err != nil {
			fmt.Println("listen fail", err)
			os.Exit(1)
		}
	}()

}

func StopWebsocket() {
	running = false
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if !running {
		http.Error(w, "system rebooting...", 406)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil || ws == nil {
		fmt.Printf("upgrader.Upgrade() failed: %v raw header: %v", err, r.Header)
		return
	}

	// todo: 处理socket连接
}


