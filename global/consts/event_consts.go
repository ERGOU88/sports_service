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

// EventType 事件类型
// 0  系统类
// 1  活动类
// 2  视频点赞
// 3  评论/回复点赞
// 4  收藏视频
// 5  关注用户
// 6  关注的用户发布新视频
// 7  视频评论
// 8  视频回复
// 9  在视频评论/回复中@
// 10 帖子点赞
// 11 帖子评论点赞
// 12 关注的人发布的新帖子
// 13 帖子评论
// 14 帖子回复
// 15 在帖子评论/回复中@
// 16 在发布帖子时@
const (
  SYSTEM_MSG                      = iota
  ACTIVITY_MSG
  VIDEO_LIKE_MSG
  VIDEO_COMMENT_LIKE_MSG
  COLLECT_VIDEO_MSG
  FOCUS_USER_MSG
  FOCUS_USER_PUBLISH_VIDEO_MSG
  VIDEO_COMMENT_MSG
  VIDEO_REPLY_MSG
  VIDEO_COMMENT_AT_MSG
  POST_LIKE_MSG
  POST_COMMENT_LIKE_MSG
  FOCUS_USER_PUBLISH_POST_MSG
  POST_COMMENT_MSG
  POST_REPLY_MSG
  POST_COMMENT_AT_MSG
  POST_PUBLISH_AT_MSG
  INFORMATION_COMMENT_MSG
  INFORMATION_REPLY_MSG
)



const (
  // 预约场馆订单超时事件
  ORDER_EVENT_VENUE_TIME_OUT           = 1
  // 预约私教订单超时
  ORDER_EVENT_COACH_TIME_OUT           = 2
  // 预约大课订单超时
  ORDER_EVENT_COURSE_TIME_OUT          = 3
)
