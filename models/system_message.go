package models

type SystemMessage struct {
	SystemId      int64  `json:"system_id" xorm:"not null pk autoincr comment('系统通知ID') BIGINT(20)"`
	SendId        string `json:"send_id" xorm:"not null default '' comment('发送者ID（后台用户）') VARCHAR(60)"`
	ReceiveId     string `json:"receive_id" xorm:"not null default '' comment('接收者id') index VARCHAR(60)"`
	Cover         string `json:"cover" xorm:"not null default '' comment('消息封面') VARCHAR(256)"`
	SendDefault   int    `json:"send_default" xorm:"not null default 0 comment('1时发送所有用户，0时则不采用') TINYINT(2)"`
	SystemTopic   string `json:"system_topic" xorm:"not null comment('通知标题') MEDIUMTEXT"`
	SystemContent string `json:"system_content" xorm:"not null comment('通知内容') MEDIUMTEXT"`
	SendTime      int    `json:"send_time" xorm:"not null comment('发送时间') INT(11)"`
	ExpireTime    int    `json:"expire_time" xorm:"not null default 0 comment('过期时间') INT(11)"`
	SendType      int    `json:"send_type" xorm:"not null default 0 comment('0.默认为系统通知 1 活动通知  2 待支付订单延时提示消息（15分钟 用户端） 3. 待咨询订单通知(开始前1天 用户端及咨询师端) 4.待咨询订单通知(开始前1小时 用户端及咨询师端) 5. 咨询师写评估报告通知（结束后1小时 咨询师端）6. 咨询师写评估报告通知（结束后24小时 咨询师端）') index TINYINT(1)"`
	Extra         string `json:"extra" xorm:"not null default ' ' comment('附件内容 例如：奖励') VARCHAR(1024)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('0 未读 1 已读  默认未读') TINYINT(1)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	AndroidTaskId string `json:"android_task_id" xorm:"not null default '' comment('友盟任务id[android端]') VARCHAR(256)"`
	IosTaskId     string `json:"ios_task_id" xorm:"not null default '' comment('友盟任务id[ios端]') VARCHAR(256)"`
	UmengPlatform int    `json:"umeng_platform" xorm:"not null default 0 comment('0 所有 1 android 2 ios') TINYINT(1)"`
	SendStatus    int    `json:"send_status" xorm:"not null default 0 comment('0 已发送 1 未发送 2 已撤回 3 已删除') TINYINT(1)"`
}
