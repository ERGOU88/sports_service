package event

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/dao"
	"sports_service/server/global/rdskey"
	"sports_service/server/global/app/log"
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
