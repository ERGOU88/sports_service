package consts

const (
  // 事件消息
  EVENT_TOPIC = "event-topic"
)

const (
  EVENT_EXCHANGE_NAME = "event-exchange"
  EVENT_QUEUE         = "event-queue"
  EVENT_ROUTING_KEY   = "event-routing"
)

// direct交换器，匹配规则为：如果路由键匹配，消息就被投送到相关的队列
// fanout交换器中没有路由键的概念，会把消息发送到所有绑定在此交换器上面的队列中。
// topic交换器采用模糊匹配路由键的原则进行转发消息到队列中
const (
  EXCHANGE_DIRECT = "direct"
  EXCHANGE_FANOUT = "fanout"
  EXCHANGE_TOPIC  = "topic"
)

const (
  SYSTEM_MSG                      = iota
  ACTIVITY_MSG
  VIDEO_LIKE_MSG
  VIDEO_COMMENT_LIKE_MSG
  COLLECT_VIDEO_MSG
  FOCUS_USER_MSG
  FOCUS_USER_PUBLISH_MSG
  VIDEO_COMMENT_MSG
  VIDEO_REPLY_MSG
  POST_LIKE_MSG
  POST_COMMENT_LIKE_MSG
  POST_COMMENT_MSG
  POST_REPLY_MSG
)
