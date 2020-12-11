package models

type SystemMessage struct {
	ExpireTime    int    `json:"expire_time" xorm:"not null comment('过期时间') INT(11)"`
	Extra         string `json:"extra" xorm:"not null default ' ' comment('附件内容 例如：奖励') VARCHAR(2056)"`
	ReceiveId     string `json:"receive_id" xorm:"not null comment('接收者id') index VARCHAR(60)"`
	SendDefault   int    `json:"send_default" xorm:"not null comment('1时发送所有用户，0时则不采用') TINYINT(2)"`
	SendId        int    `json:"send_id" xorm:"not null comment('发送者ID（后台用户id）') INT(11)"`
	SendTime      int    `json:"send_time" xorm:"not null comment('发送时间') INT(11)"`
	SendType      int    `json:"send_type" xorm:"not null default 0 comment('0.默认为系统通知 1.收到@通知 2.收到点赞通知 3.收到收藏通知 4.收到分享通知 5.收到评论/回复通知 6.特殊业务的系统奖励通知 7.活动通知') index TINYINT(1)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('0 未读 1 已读  默认未读') TINYINT(1)"`
	SystemContent string `json:"system_content" xorm:"not null default '' comment('通知内容') VARCHAR(1000)"`
	SystemId      int64  `json:"system_id" xorm:"not null pk autoincr comment('系统通知ID') BIGINT(20)"`
	SystemTopic   string `json:"system_topic" xorm:"not null default '' comment('通知标题') VARCHAR(500)"`
}
