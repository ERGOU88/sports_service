package consts

type MessageType int32

const (
  // 系统类推送 -1
  MSG_TYPE_SYSTEM_NOTIFY                      MessageType = -1

  // 活动类推送消息
  MSG_TYPE_ACTIVITY_NOTIFY                    MessageType = 10000
)

const (
  // android: notification 通知栏推送  message 自定义推送
  ANDROID_PUSH_TYPE_NOTIFICATION  = "notification"
  ANDROID_PUSH_TYPE_CUSTOM        = "message"
)
