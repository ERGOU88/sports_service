package event

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/rdskey"
	"sports_service/server/nsqlx/protocol"
	"sports_service/server/util"
	"time"
)

// 事件消息
func PushEventMsg(msg []byte) {
	log.Log.Infof("event_trace: 事件推送, msg:%s", string(msg))
	//body := newEvent(userId, nickname, cover, content, eventType)
	conn := dao.RedisPool().Get()

	num, err := redis.Int(conn.Do("LPUSH", rdskey.MSG_EVENT_KEY, msg))
	if err != nil || num != 1 {
		log.Log.Infof("event_trace: msg push fail, err:%s", err)
	}
}

func NewEvent(userId, nickname, cover, content string, eventType int32) []byte {
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

