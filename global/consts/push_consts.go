package consts

type MessageType int32

const (
  // 系统类推送 -1
  MSG_TYPE_SYSTEM_NOTIFY                      MessageType = -1

  // 活动类推送消息
  MSG_TYPE_ACTIVITY_NOTIFY                    MessageType = 10000
  // 视频点赞推送
  MSG_TYPE_VIDEO_LIKE_NOTIFY                  MessageType = 20000
  // 评论/回复点赞推送
  MSG_TYPE_COMMENT_LIKE_NOTIFY                MessageType = 20001
  // 收藏视频推送
  MSG_TYPE_VIDEO_COLLECT_NOTIFY               MessageType = 30000
  // 关注推送
  MSG_TYPE_FOCUS_NOTIFY                       MessageType = 40000
  // 关注的用户发布新视频推送
  MSG_TYPE_FOCUS_USER_PUBLISH_NOTIFY          MessageType = 40001
  // 视频评论推送
  MSG_TYPE_VIDEO_COMMENT_NOTIFY               MessageType = 50000
  // 视频回复推送
  MSG_TYPE_VIDEO_REPLY_NOTIFY                 MessageType = 50001

)

const (
  // android: notification 通知栏推送  message 自定义推送
  ANDROID_PUSH_TYPE_NOTIFICATION  = "notification"
  ANDROID_PUSH_TYPE_CUSTOM        = "message"
)
