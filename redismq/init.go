package redismq

import (
	"sports_service/server/redismq/achieve/event"
)

func InitRedisMq() {
	event.LoopPopStatEvent()
	//go event.LoopPopOrderEvent()
	event.LoopPopTopEvent()
	event.InitSignal()
}
