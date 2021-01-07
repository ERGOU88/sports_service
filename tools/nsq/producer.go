package nsq

import (
	"fmt"
	nsq "github.com/nsqio/go-nsq"
)

var NsqProducer *nsq.Producer

func ConnectNsqProduct(nsqAddr string) {
	fmt.Println("init nsq product start...")
	var err error
	NsqProducer, err = nsq.NewProducer(nsqAddr, nsq.NewConfig())
	if err != nil {
		panic(fmt.Sprintf("nsq connect error:%s", err))
		return
	}

	fmt.Println("init nsq  product end...")
}
