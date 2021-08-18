package redismq

import (
	"sports_service/server/redismq/achieve/event"
)

func InitRedisMq() {
	go event.LoopPopStatEvent()
	go event.LoopPopOrderEvent()
	event.InitSignal()
}
