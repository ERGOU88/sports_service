package main

import (
	"github.com/golang/protobuf/proto"
	"sports_service/server/global/consts"
	"sports_service/server/global/rdskey"
	pbBarrage "sports_service/server/proto/barrage"
	"sports_service/server/tools/nsq"
	nsqsvc "github.com/nsqio/go-nsq"
	"sports_service/server/barrage/config"
	"fmt"
	"time"
	"errors"
	"sports_service/server/dao"
)

var msgCh = make(chan *pbBarrage.Message, config.Global.MaxMsgCacheLen)

// 初始化消费者
func InitNsqConsumer() {
	nsq.HandleConsumer(MessageConsumer, 3, "message")
}

func connectConsumer(channel string) (*nsqsvc.Consumer, error) {
	nsqConfig := nsqsvc.NewConfig()
	fmt.Printf("初始化 topic: %v,channel:%v", consts.MESSAGE_TOPIC, channel)
	consumer, err := nsqsvc.NewConsumer(consts.MESSAGE_TOPIC, channel, nsqConfig)
	if err != nil {
		fmt.Printf("new consumer err:%v", err)
		return consumer, err
	}

	return consumer, err
}

func MessageConsumer(channel string) (consumer *nsqsvc.Consumer) {
	consumer, err := connectConsumer(channel)
	if err != nil {
		panic(fmt.Sprintf("consumer conn err:%s", err))
	}

	consumer.AddHandler(nsqsvc.HandlerFunc(NsqHandlerMsg))

	err = consumer.ConnectToNSQD(config.Global.NsqAddr)
	if err != nil {
		fmt.Printf("ConnectToNSQD err:%s", err)
	}

	return
}

func NsqHandlerMsg(msg *nsqsvc.Message) error {
	event := new(pbBarrage.Message)
	if err := proto.Unmarshal(msg.Body, event); err != nil {
		fmt.Printf("message_event: proto unmarshal err:%s", err)
		return err
	}

	if err := handle(event); err != nil {
		msg.RequeueWithoutBackoff(time.Second * 3)
		fmt.Printf("handle err:%s", err)
		return err
	}

	return nil
}

func handle(event *pbBarrage.Message) error {
	if len(msgCh) < config.Global.MaxMsgCacheLen {
		msgCh <- event
		return nil
	}

	return errors.New("msg channel full")
}

// 读取消息
func ReadChanelMessage() {
	var msg *pbBarrage.Message
	for {
		msg = <- msgCh
		if msg == nil {
			continue
		}
		switch msg.MsgType {
		// 文本消息
		case pbBarrage.MessageType_TYPE_TEXT:
		// 视频弹幕消息(主动推)
		case pbBarrage.MessageType_TYPE_BARRAGE:
			fmt.Printf("\nbarrage msg.body:%s", string(msg.Body))
			barrage := pbBarrage.BarrageMessage{}
			if err := proto.Unmarshal(msg.Body, &barrage); err != nil {
				fmt.Printf("\nbarrage msg, proto unmarshal err:%s", err)
				continue
			}
			// 获取正在观看该视频的用户们
			list, err := getWatchUserByVideoId(barrage.Barrage.VideoId)
			if err != nil || len(list) == 0 {
				fmt.Printf("\nbarrage msg, get watch user err:%s, videoId:%s", err, barrage.Barrage.VideoId)
				continue
			}

			for _, xid := range list {
				// 推送弹幕消息
				PushMessage(xid, msg.Body, pbBarrage.MessageType_TYPE_BARRAGE)
			}

		// 广播消息
		case pbBarrage.MessageType_TYPE_BROADCAST:
			// todo:


		// 观看视频消息
		case pbBarrage.MessageType_TYPE_WATCH_VIDEO:
			watchMsg := &pbBarrage.ReqWatchVideo{}
			if err := proto.Unmarshal(msg.Body, watchMsg); err != nil {
				fmt.Printf("\nwatch video msg, proto unmarshal err:%s", err)
				continue
			}

			fmt.Printf("\nwatchMsg:%v", watchMsg)

			if _, err := recordWatchUserByVideoId(watchMsg.VideoId, watchMsg.Xid); err != nil {
				fmt.Printf("\nrecord watch user err:%s, videoId:%s", err, watchMsg.VideoId)
			}

		// 退出观看消息
		case pbBarrage.MessageType_TYPE_EXIT_VIDEO:
			exitMsg := &pbBarrage.ReqExitVideo{}
			if err := proto.Unmarshal(msg.Body, exitMsg); err != nil {
				fmt.Printf("exit watch video msg, proto unmarshal err:%s", err)
				continue
			}

			if _, err := delWatchUserByVideoId(exitMsg.VideoId, exitMsg.Xid); err != nil {
				fmt.Printf("del watch user err:%s, videoId:%s", err, exitMsg.VideoId)
			}

		}

	}
}

// 获取正在观看该视频的用户们
func getWatchUserByVideoId(videoId string) ([]string, error) {
	rds := dao.NewRedisDao()
	return rds.SMEMBERS(rdskey.MakeKey(rdskey.USER_WATCHING_VIDEO, videoId))
}

// 记录正在观看该视频的用户们
func recordWatchUserByVideoId(videoId, xid string) (int, error) {
	rds := dao.NewRedisDao()
	return rds.SADD(rdskey.MakeKey(rdskey.USER_WATCHING_VIDEO, videoId), xid)
}

// 删除退出观看视频的用户
func delWatchUserByVideoId(videoId, xid string) (int, error) {
	rds := dao.NewRedisDao()
	return rds.SREM(rdskey.MakeKey(rdskey.USER_WATCHING_VIDEO, videoId), xid)
}

