package event

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/dao"
	"sports_service/server/global/rdskey"
	"sports_service/server/global/app/log"
	"sports_service/server/redismq/protocol"
	"sports_service/server/util"
	"time"
)

// 订单事件推送
func PushOrderEventMsg(msg []byte) {
	log.Log.Infof("event_trace: 订单事件推送, msg:%s", string(msg))
	conn := dao.RedisPool().Get()

	num, err := redis.Int(conn.Do("LPUSH", rdskey.MSG_ORDER_EVENT_KEY, msg))
	if err != nil || num != 1 {
		log.Log.Infof("event_trace: msg push fail, err:%s", err)
	}
}

func NewOrderEvent(toUserId, composeId, nickname, cover, content string, eventType int32) []byte {
	event := new(protocol.Event)
	event.UserId = toUserId
	event.EventType = eventType
	event.Ts = time.Now().Unix()

	data := new(protocol.PushData)
	data.NickName = nickname
	data.Cover = cover
	data.Content = content
	data.ComposeId = composeId

	msg , _ := util.JsonFast.Marshal(data)
	event.Data = msg
	b, err := util.JsonFast.Marshal(event)
	if err != nil {
		log.Log.Errorf("event_trace: marshal err:%s", err)
	}

	return b
}