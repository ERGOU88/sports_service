package event

import (
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/nsqlx/protocol"
	"sports_service/tools/amqp"
	"sports_service/util"
	"time"
)

// 事件消息
func PushEventMsg(amqpDsn, userId, nickname, cover, content string, eventType int32) {
	log.Log.Errorf("event_trace: 事件推送，eventType:%d", eventType)
	eventPublish(amqpDsn, userId, nickname, cover, content, eventType)
}

// 事件
func eventPublish(amqpDsn, userId, nickname, cover, content string, eventType int32) {
	// 建立会话
	session, err := amqp.NewSession(amqpDsn)
	if err != nil {
		log.Log.Errorf("amqp_trace: new session fail, err:%s", err)
		return
	}
	defer session.Close()

	// 生产者
	producer, err := amqp.NewProducer(session, consts.EVENT_EXCHANGE_NAME, consts.EXCHANGE_DIRECT)
	if err != nil {
		log.Log.Errorf("amqp_trace: new producer fail, err:%s", err)
		return
	}
	defer producer.Close()

	body := newEvent(userId, nickname, cover, content, eventType)
	if err = producer.Publish(
		consts.EVENT_ROUTING_KEY,
		"application/json",
		string(body),
	); err != nil {
		log.Log.Errorf("amqp_trace: publish fail, err:%s", err)
	}
}

func newEvent(userId, nickname, cover, content string, eventType int32) []byte {
	event := new(protocol.Event)
	event.UserId = userId
	event.EventType = eventType
	event.Ts = time.Now().Unix()

	data := new(protocol.Data)
	data.NickName = nickname
	data.Cover = cover
	data.Content = content

	msg, _ := util.JsonFast.Marshal(data)
	event.Data = msg
	b, err := util.JsonFast.Marshal(event)
	if err != nil {
		log.Log.Errorf("event_trace: marshal err:%s", err)
	}

	return b
}
