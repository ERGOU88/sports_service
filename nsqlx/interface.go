package nsqlx

import (
	"sports_service/nsqlx/achieve/event"
	"sports_service/tools/nsq"
)

// 初始化消费者
func InitNsqConsumer() {
	nsq.HandleConsumer(event.EventConsumer, 3, "event")
}
