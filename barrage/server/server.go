package main

import (
	"sports_service/server/global/consts"
	pbBarrage "sports_service/server/proto/barrage"
	"sports_service/server/tools/nsq"
	"fmt"
	"time"
	"github.com/golang/protobuf/proto"
)

func main () {
	nsq.ConnectNsqProduct("127.0.0.1:4150")
	// 正在看视频
	FakeWatchingVideo()
	// 正在发弹幕
	FakeBarrage()
}

func FakeWatchingVideo() {
	watch := pbBarrage.ReqWatchVideo{
		Xid: "1",
		VideoId: "59",
		Uid: "1",
	}

	bts, err := proto.Marshal(&watch)
	if err != nil {
		fmt.Println("\n proto marshal err:%s", err)
	}

	msg := &pbBarrage.Message{
		MsgType: pbBarrage.MessageType_TYPE_WATCH_VIDEO,
		MsgId: "1",
		Body: bts,
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("\n proto marshal data err:%s", err)
	}

	if err := nsq.NsqProducer.Publish(consts.MESSAGE_TOPIC, data); err != nil {
		fmt.Println("\n nsq publish err:%s", err)
	}
}

func FakeBarrage() {
	barrage := &pbBarrage.BarrageMessage{
		Xid:"1",
		Barrage: &pbBarrage.BarrageInfo{
			Uid: "1",
			Content: "弹幕001",
			VideoId: "59",
			CurDuration: 30,
			SendTime: time.Now().Unix(),
		},
	}

	bts, err := proto.Marshal(barrage)
	if err != nil {
		fmt.Println("\n proto marshal err:%s", err)
	}

	msg := &pbBarrage.Message{
		MsgType: pbBarrage.MessageType_TYPE_BARRAGE,
		MsgId: "1",
		Body: bts,
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("\n proto marshal data err:%s", err)
	}

	if err := nsq.NsqProducer.Publish(consts.MESSAGE_TOPIC, data); err != nil {
		fmt.Println("\n nsq publish err:%s", err)
	}
}
