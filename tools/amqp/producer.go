package amqp

import (
  "github.com/streadway/amqp"
)

// Producer 生产者
type Producer struct {
  ExchangeName string
  ExchangeType string
  Session      *Session
  Ch           *amqp.Channel
}

// NewProducer 得到生产者对象
// exchangeType:
// fanout:把所有发送到该Exchange的消息投递到所有与它绑定的队列中。
// direct:把消息投递到那些binding key与routing key完全匹配的队列中。
// topic:将消息路由到binding key与routing key模式匹配的队列中。
// direct交换器，匹配规则为：如果路由键匹配，消息就被投送到相关的队列
// fanout交换器中没有路由键的概念，会把消息发送到所有绑定在此交换器上面的队列中。
// topic交换器采用模糊匹配路由键的原则进行转发消息到队列中
func NewProducer(session *Session, exchangeName, exchangeType string) (*Producer, error) {
  producer := &Producer{
    Session:      session,
    ExchangeName: exchangeName,
    ExchangeType: exchangeType,
  }
  channel, err := producer.Session.Conn.Channel()
  if err != nil {
    return nil, err
  }
  producer.Ch = channel

  err = producer.Ch.ExchangeDeclare(
    exchangeName, // name of the exchange
    exchangeType, // type
    true,         // durable
    false,        // delete when complete
    false,        // internal
    false,        // noWait
    nil,
  )
  if err != nil {
    return nil, err
  }

  return producer, nil
}

// Publish  消息投递
func (p *Producer) Publish(routingKey, contentType, body string) error {
  err := p.Ch.Publish(
    p.ExchangeName,
    routingKey,
    false, false,
    amqp.Publishing{
      ContentType:  contentType,
      Body:         []byte(body),
      DeliveryMode: amqp.Transient,
    })
  return err
}

// Close  关闭通道
func (p *Producer) Close() error {
  return p.Ch.Close()
}
