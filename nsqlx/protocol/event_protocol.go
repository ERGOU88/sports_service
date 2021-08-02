package protocol


// 事件
// EventType 事件类型
// 0 系统类
// 1 活动类
// 2 视频点赞
// 3 评论/回复点赞
// 4 收藏视频
// 5 关注用户
// 6 关注的用户发布新视频
// 7 视频评论
// 8 视频回复
type Event struct {
  UserId    string `json:"user_id"`   // 用户id
  EventType int32  `json:"eventType"`
  Ts        int64  `json:"ts"`        // 时间
  Data      []byte `json:"data"`
}

// 事件数据
type Data struct {
  Cover       string    `json:"cover"`      // 封面
  NickName    string    `json:"nick_name"`  // 昵称
  Content     string    `json:"content"`    // 内容
  ComposeId   string    `json:"compose_id"` // 作品id 视频/帖子
}

