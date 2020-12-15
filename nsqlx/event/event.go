package event

import (
  "sports_service/server/nsqlx/protocol"
  "sports_service/server/tools/nsq"
  "sports_service/server/util"
  "time"
  "sports_service/server/global/app/log"
  "sports_service/server/global/consts"
)

// 事件消息
func PushEventMsg(userId, nickname, cover, content string, eventType int32) {
  log.Log.Errorf("event_trace: 事件推送，eventType:%d", eventType)
  eventNSQPub(userId, nickname, cover, content, eventType)
}

// 事件
func eventNSQPub(userId, nickname, cover, content string, eventType int32) {
  body := newEvent(userId, nickname, cover, content, eventType)
  if err := nsq.NsqProducer.Publish(consts.EVENT_TOPIC, body); err != nil {
    log.Log.Errorf("event_trace: publish event err:%s, uid:%s", err, userId)
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

  msg , _ := util.JsonFast.Marshal(data)
  event.Data = msg
  b, err := util.JsonFast.Marshal(event)
  if err != nil {
    log.Log.Errorf("event_trace: marshal err:%s", err)
  }

  return b
}

