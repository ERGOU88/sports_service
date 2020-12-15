package protocol

// 事件
type Event struct {
  Uid       string          `json:"uid"`       // 用户id
  EventType int32           `json:"eventType"` // 事件类型  1.预约咨询师 订单超时 30分钟 2.订单付款提示 15分钟
  Ts        int64           `json:"ts"`        // 时间
  Data      []byte          `json:"data"`
}

// 事件数据
type Data struct {
  Cover       string    `json:"cover"`      // 视频封面
  NickName    string    `json:"nick_name"`  // 点赞人昵称
  Content     string    `json:"content"`    // 点赞内容
}

