package redismq

import (
	"sports_service/server/redismq/achieve/event"
)

func InitRedisMq() {
	event.LoopPopStatEvent()
	event.InitSignal()
}
