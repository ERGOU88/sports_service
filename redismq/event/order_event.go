package event

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/dao"
	"sports_service/global/app/log"
	"sports_service/global/rdskey"
	"sports_service/redismq/protocol"
	"sports_service/util"
	"time"
)

// 订单事件推送
func PushOrderEventMsg(msg []byte) {
	log.Log.Infof("event_trace: 订单事件推送, msg:%s", string(msg))
	conn := dao.RedisPool().Get()
	defer conn.Close()

	if _, err := redis.Int(conn.Do("LPUSH", rdskey.MSG_ORDER_EVENT_KEY, msg)); err != nil {
		log.Log.Infof("event_trace: msg push fail, err:%s", err)
	}

	return
}

// 订单超时事件
func NewOrderEvent(userId, orderId string, processTm int64, eventType int32) []byte {
	event := new(protocol.Event)
	event.EventType = eventType
	event.Ts = time.Now().Unix()
	event.UserId = userId

	data := new(protocol.OrderData)
	data.OrderId = orderId
	data.ProcessTm = processTm

	msg, _ := util.JsonFast.Marshal(data)
	event.Data = msg
	b, err := util.JsonFast.Marshal(event)
	if err != nil {
		log.Log.Errorf("event_trace: marshal err:%s", err)
	}

	return b
}
