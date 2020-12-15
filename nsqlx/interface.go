package nsqlx

import (
  "sports_service/server/nsqlx/achieve/event"
  "sports_service/server/tools/nsq"
)

// 初始化消费者
func InitNsqConsumer() {
 nsq.HandleConsumer(event.EventConsumer, 3, "event")
}
