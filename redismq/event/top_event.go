package event

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/rdskey"
	"sports_service/server/redismq/protocol"
	"sports_service/server/util"
	"time"
)

// 作品置顶事件消息
func PushTopEventMsg(msg []byte) {
	log.Log.Infof("event_trace: 置顶事件, msg:%s", string(msg))
	conn := dao.RedisPool().Get()
	defer conn.Close()

	if _, err := redis.Int(conn.Do("LPUSH", rdskey.MSG_TOP_EVENT_KEY, msg)); err != nil  {
		log.Log.Infof("event_trace: msg push fail, err:%s", err)
	}

	return
}

func NewTopEvent(userId, id string, eventType int32) []byte {
	event := new(protocol.Event)
	event.UserId = userId
	event.EventType = eventType
	event.Ts = time.Now().Unix()

	data := &protocol.WorkInfo{
		Id: id,
	}

	msg , _ := util.JsonFast.Marshal(data)
	event.Data = msg
	b, err := util.JsonFast.Marshal(event)
	if err != nil {
		log.Log.Errorf("event_trace: marshal err:%s", err)
	}

	return b
}

