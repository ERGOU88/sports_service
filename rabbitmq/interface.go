package rabbitmq

import (
  "sports_service/server/rabbitmq/achieve/event"
)

// 初始化消费者
func InitRabbitmqConsumer() {
  go event.ConnectEventConsumer()
}
