package rabbitmq

import (
	"sports_service/global/app/log"
	"sports_service/rabbitmq/achieve/event"
)

// 初始化消费者
func InitRabbitmqConsumer() {
	if err := event.ConnectEventConsumer(); err != nil {
		log.Log.Errorf("amqp_trace: connect event consumer fail, err:%s", err)
		panic(err)
	}

	log.Log.Debug("setup rabbitmq success")
}
