package amqp

import (
  "fmt"
  "time"
  "github.com/streadway/amqp"
)

// Consumer 消费者对象
type Consumer struct {
  QueueName    string
  ExchangeName string
  ExchangeType string
  RoutingKey   string
  Ch           *amqp.Channel
  Session      *Session
  Queue        amqp.Queue
}

// NewConsumer 得到消费者对象
func NewConsumer(session *Session, queueName, delayQueue, exchangeName, exchangeType, routingKey string) (*Consumer, error) {
  consumer := &Consumer{
    Session:      session,
    QueueName:    queueName,
    ExchangeName: exchangeName,
    ExchangeType: exchangeType,
    RoutingKey:   routingKey,
  }

  channel, err := consumer.Session.Conn.Channel()
  if err != nil {
    return nil, err
  }
  consumer.Ch = channel

  // 确保rabbitMQ一个一个发送消息
  if err = consumer.Ch.Qos(200, 0, true); err != nil {
    return nil, err
  }

  err = consumer.Ch.ExchangeDeclare(
    exchangeName, // name of the exchange
    exchangeType, // type
    true,          // durable
    false,      // delete when complete
    false,        // internal
    false,        // noWait
    nil,
  )
  if err != nil {
    return nil, err
  }
  consumer.Queue, err = consumer.Ch.QueueDeclare(
    queueName, // name
    true,      // durable  持久性的,如果事前已经声明了该队列，不能重复声明
    false,     // delete when unused
    false,     // exclusive 如果是真，连接一断开，队列删除
    false,     // no-wait
    nil,       // arguments
  )
  if err != nil {
    return nil, err
  }

  // 延时消息队列名
  if delayQueue != "" {
    // 延时消息队列
    if _, err = consumer.Ch.QueueDeclare(
      delayQueue, // name
      true,       // durable  持久性的,如果事前已经声明了该队列，不能重复声明
      false,      // delete when unused
      false,      // exclusive 如果是真，连接一断开，队列删除
      false,      // no-wait
      amqp.Table{
        // 将过期的消息发送到指定的 exchange 中
        "x-dead-letter-exchange": exchangeName,
        // 将过期的消息发送到自定的 route当中
        "x-dead-letter-routing-key": routingKey,
      }, // arguments
    ); err != nil {
      return nil, err
    }
  }

  // 队列和交换机绑定，即是队列订阅了发到这个交换机的消息
  err = consumer.Ch.QueueBind(
    queueName,
    routingKey,
    exchangeName,
    false,
    nil,
  )
  if err != nil {
    return nil, err
  }

  return consumer, nil
}

// Consume 消费队列
func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
  msgChan, err := c.Ch.Consume(
    c.QueueName, // queue
    "",          // consumer
    false,       // auto-ack   设置为真自动确认消息
    false,       // exclusive
    false,       // no-local
    false,       // no-wait
    nil,         // args)
  )
  if err != nil {
    return nil, err
  }

  return msgChan, nil
}

// ReConnect 重连
func (c *Consumer) ReConnect() {
  tikcer := time.NewTicker(time.Second * 3)
  for range tikcer.C {
    var err error
    c.Session.Conn, err = amqp.Dial(c.Session.DSN)
    if err != nil {
      fmt.Print(err)
      continue
    }
    c.Ch, err = c.Session.Conn.Channel()
    if err != nil {
      fmt.Print(err)
      continue
    }
    break
  }
}

// Close 关闭consumer
func (c *Consumer) Close() {
  c.Ch.Close()
}
