package consts

type MessageType int32

const (
  // 系统类推送 -1
  MSG_TYPE_SYSTEM_NOTIFY                      MessageType = -1

  // 活动类推送消息
  MSG_TYPE_ACTIVITY_NOTIFY                    MessageType = 10000
  // 视频点赞推送
  MSG_TYPE_VIDEO_LIKE_NOTIFY                  MessageType = 20000
  // 视频评论/回复点赞推送
  MSG_TYPE_VIDEO_COMMENT_LIKE_NOTIFY          MessageType = 20001
  // 帖子点赞推送
  MSG_TYPE_POST_LIKE_NOTIFY                   MessageType = 20002
  // 帖子评论/回复点赞推送
  MSG_TYPE_POST_COMMENT_LIKE_NOTIFY           MessageType = 20003
  // 收藏视频推送
  MSG_TYPE_VIDEO_COLLECT_NOTIFY               MessageType = 30000
  // 关注推送
  MSG_TYPE_FOCUS_NOTIFY                       MessageType = 40000
  // 关注的用户发布新视频推送
  MSG_TYPE_FOCUS_PUBLISH_VIDEO_NOTIFY         MessageType = 40001
  // 关注的用户发布新帖子推送
  MSG_TYPE_FOCUS_PUBLISH_POST_NOTIFY          MessageType = 40002
  // 视频评论推送
  MSG_TYPE_VIDEO_COMMENT_NOTIFY               MessageType = 50000
  // 视频回复推送
  MSG_TYPE_VIDEO_REPLY_NOTIFY                 MessageType = 50001
  // 帖子评论推送
  MSG_TYPE_POST_COMMENT_NOTIFY                MessageType = 50002
  // 帖子回复推送
  MSG_TYPE_POST_REPLY_NOTIFY                  MessageType = 50003
  // 视频评论/回复中 @
  MSG_TYPE_VIDEO_COMMENT_AT_NOTIFY            MessageType = 50004
  // 帖子评论/回复中 @
  MSG_TYPE_POST_COMMENT_AT_NOTIFY             MessageType = 50005
  // 发布帖子内容中 @
  MSG_TYPE_POST_PUBLISH_AT_NOTIFY             MessageType = 50006

)

var NotifyDoc = map[MessageType]string{
  MSG_TYPE_SYSTEM_NOTIFY :             "系统类推送",
  MSG_TYPE_ACTIVITY_NOTIFY:            "活动类推送消息",
  MSG_TYPE_VIDEO_LIKE_NOTIFY :         "视频点赞推送",
  MSG_TYPE_VIDEO_COMMENT_LIKE_NOTIFY:  "视频评论/回复点赞推送",
  MSG_TYPE_POST_LIKE_NOTIFY:           "帖子点赞推送",
  MSG_TYPE_POST_COMMENT_LIKE_NOTIFY:   "帖子评论/回复点赞推送",
  MSG_TYPE_VIDEO_COLLECT_NOTIFY :      "收藏视频推送",
  MSG_TYPE_FOCUS_NOTIFY:               "关注推送",
  MSG_TYPE_FOCUS_PUBLISH_VIDEO_NOTIFY: "关注的用户发布新视频推送",
  MSG_TYPE_FOCUS_PUBLISH_POST_NOTIFY:  "关注的用户发布新帖子推送",
  MSG_TYPE_VIDEO_COMMENT_NOTIFY:       "视频评论推送",
  MSG_TYPE_VIDEO_REPLY_NOTIFY:         "视频回复推送",
  MSG_TYPE_POST_COMMENT_NOTIFY:        "帖子评论推送",
  MSG_TYPE_POST_REPLY_NOTIFY:          "帖子回复推送",
  MSG_TYPE_VIDEO_COMMENT_AT_NOTIFY:    "视频评论/回复中@",
  MSG_TYPE_POST_COMMENT_AT_NOTIFY:     "帖子评论/回复中@",
  MSG_TYPE_POST_PUBLISH_AT_NOTIFY:     "帖子发布时 内容中@",
}

const (
  // android: notification 通知栏推送  message 自定义推送
  ANDROID_PUSH_TYPE_NOTIFICATION  = "notification"
  ANDROID_PUSH_TYPE_CUSTOM        = "message"
)
