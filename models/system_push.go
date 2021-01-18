package models

type SystemPush struct {
	PushId      int64  `json:"push_id" xorm:"not null pk autoincr comment('系统推送ID') BIGINT(20)"`
	SendId      int    `json:"send_id" xorm:"not null comment('发送者ID（后台用户id）') INT(11)"`
	ReceiveId   string `json:"receive_id" xorm:"not null comment('接收者id') index VARCHAR(60)"`
	SendDefault int    `json:"send_default" xorm:"not null comment('1时发送所有用户，0时则不采用') TINYINT(2)"`
	PushTopic   string `json:"push_topic" xorm:"not null comment('推送标题') VARCHAR(60)"`
	PushContent string `json:"push_content" xorm:"not null comment('推送内容') VARCHAR(255)"`
	SendTime    int    `json:"send_time" xorm:"not null comment('发送时间') INT(11)"`
	SendType    int    `json:"send_type" xorm:"not null default 0 comment('0.默认为系统推送 1.收到@推送 2.收到点赞推送 3.收到收藏推送 4.收到分享推送 5.收到评论/回复推送 6.特殊业务的系统奖励推送 7.活动推送') index TINYINT(1)"`
	Extra       string `json:"extra" xorm:"not null default ' ' comment('透传数据 ') VARCHAR(2056)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('0 成功 1 失败') TINYINT(1)"`
}
