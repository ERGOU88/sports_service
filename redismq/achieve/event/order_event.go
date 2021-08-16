package event

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/dao"
	"sports_service/server/global/rdskey"
	producer "sports_service/server/redismq/event"
	"time"
	"sports_service/server/global/app/log"
)

func LoopPopOrderEvent() {
	for !closing {
		conn := dao.RedisPool().Get()
		values, err := redis.Values(conn.Do("BRPOP", rdskey.MSG_ORDER_EVENT_KEY, 0))
		conn.Close()
		if err != nil {
			log.Log.Errorf("redisMq_trace: loopPop event fail, err:%s", err)
			// 防止出现错误时 频繁刷日志
			time.Sleep(time.Second)
			continue
		}

		if len(values) < 2 {
			log.Log.Errorf("redisMq_trace: invalid values, len:%d, values:%+v", len(values), values)
		}


		bts, ok := values[1].([]byte)
		if !ok {
			log.Log.Errorf("redisMq_trace: value[1] unSupport type")
			continue
		}

		if err := EventConsumer(bts); err != nil {
			log.Log.Errorf("redisMq_trace: event consumer fail, err:%s, msg:%s", err, string(bts))
			// 重新投递消息
			producer.PushEventMsg(bts)
		}

	}
}
