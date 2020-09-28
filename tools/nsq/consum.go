package nsq

import (
	"github.com/nsqio/go-nsq"
)

type ConsumerFunc func(channel string) *nsq.Consumer

var Consumers []*nsq.Consumer

func HandleConsumer(handle ConsumerFunc, chNums int, channel string) {
	for i := 1; i <= chNums; i++ {
		Consumers = append(Consumers, handle(channel))
	}
}

func Stop() {
	for _, conn := range Consumers {
		conn.Stop()
	}
}
