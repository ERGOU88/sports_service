package amqp

import (
  "sports_service/server/global/consts"
  "testing"
)

func TestProducer(t *testing.T) {
  session, err := NewSession("amqp://admin:admin@127.0.0.1:5672")
  if err != nil {
    t.Logf("new session fail, err:%s", err)
    return
  }
  defer session.Close()
  // 生产者
  producer, err := NewProducer(session, consts.EVENT_EXCHANGE_NAME, "direct")
  if err != nil {
    t.Logf("new producer fail, err:%s", err)
    return
  }
  defer producer.Close()


  producer.Publish(consts.EVENT_ROUTING_KEY, "application/json", "hello world~")
  // 延时消息 2000毫秒
  producer.DeferredPublish("delayMsg", "application/json", "delay msg~", "5000")
}

func TestConsumer(t *testing.T) {
  session, err := NewSession("amqp://admin:admin@127.0.0.1:5672")
  if err != nil {
    t.Logf("new session fail, err:%s", err)
    return
  }
  defer session.Close()

  consumer, err := NewConsumer(session, consts.EVENT_QUEUE, "delayMsg", consts.EVENT_EXCHANGE_NAME, "direct", consts.EVENT_ROUTING_KEY)
  if err != nil {
    t.Logf("new consumer fail, err:%s", err)
    return
  }
  defer consumer.Close()

  events, err := consumer.Consume()
  if err != nil {
    t.Logf("consume queue fail, err:%s", err)
    return
  }

  for dataBody := range events {
    t.Log(string(dataBody.Body))
    dataBody.Ack(false)
  }

}
