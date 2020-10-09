package client

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	pbBarrage "sports_service/server/proto/barrage"
	"sync/atomic"
	"time"
)

var addr = flag.String("addr", "127.0.0.1:15001", "http service address")

var urls string = ""

var client, receive int64 = 0, 0

func ClientConn() {
	//flag.Parse()
	log.SetFlags(0)
	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:15001", Path: "/ws"}
	urls = u.String()

	log.Printf("connecting to %s", urls)

	go FakeClient()

	//for {
	//	r := atomic.SwapInt64(&receive, 0)
	//	fmt.Println(client, r)
	//	time.Sleep(time.Second)
	//}
	//
	//select {}
}

func FakeClient() {
	c, _, err := websocket.DefaultDialer.Dial(urls, nil)
	if err != nil {
		fmt.Println("dial error", err)
		log.Fatal("dial:", err)
		return
	}

	defer c.Close()
	req := pbBarrage.ReqConnMessage{
		AppId:  "mj4mQaop",
		Secret: "adf4OisG",
		Timestamp: "1588888888",
		Version: "1.1.0",
		Sign: "72fbf038825db3acfd89506d9dc91c42",
	}

	bts, err := proto.Marshal(&req)
	if err != nil {
		fmt.Println("marshal err:", err)
	}

	xx := pbBarrage.Message{
		MsgId: "1",
		MsgType: pbBarrage.MessageType_TYPE_CONN,
		Body: bts,
	}

	fmt.Println(xx)

	data, err := proto.Marshal(&xx)
	if err != nil {
		fmt.Println("proto marshal err:", err)
	}

	if err := c.WriteMessage(websocket.BinaryMessage, data); err != nil {
		fmt.Println("\nfailed")
	}

	go func() {
		heartbeat := pbBarrage.ReqHeartBeatMessage{
			Xid: "1",
		}

		hb, err := proto.Marshal(&heartbeat)
		if err != nil {
			fmt.Println("error")
		}

		for {
			time.Sleep(100 * time.Second)
			xx2 := pbBarrage.Message{
				MsgId: "1",
				MsgType: pbBarrage.MessageType_TYPE_HEART_BEAT,
				Body: hb,
			}

			fmt.Println(xx2)
			data, err := proto.Marshal(&xx2)
			if err != nil {
				fmt.Println("\nproto marshal err:", err)
			}

			err = c.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				fmt.Println("\nfailed heartbeat")
			} else {
				fmt.Println("\nheartbeat success")
			}
		}
	}()

	for {
		msgType, readMsg, err := c.ReadMessage()
		if err != nil {
			log.Println("\nread:", err)
			atomic.AddInt64(&client, -1)
			return
		}

		msg := pbBarrage.Message{}
		if err := proto.Unmarshal(readMsg, &msg); err != nil {
			fmt.Printf("\nproto unamrshal msg err:%s", err)
		}

		fmt.Printf("\nmsgType:%d", msgType)

		if msg.MsgType == pbBarrage.MessageType_TYPE_CONN_RES {
			atomic.AddInt64(&client, 1)
		}

		atomic.AddInt64(&receive, 1)

		switch msg.MsgType {
		case pbBarrage.MessageType_TYPE_BARRAGE:
			barrage := &pbBarrage.BarrageMessage{}
			if err := proto.Unmarshal(msg.Body, barrage); err != nil {
				fmt.Println("\nproto unmarshal err:", err)
				continue
			}

			log.Printf("\nrecv: %v", barrage)
		}
	}

}
