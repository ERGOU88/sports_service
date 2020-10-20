package models

type YcoinReceiveRecord struct {
	CreateAt int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	Explain  string `json:"explain" xorm:"not null default '' comment('任务描述') VARCHAR(256)"`
	Id       int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	Name     string `json:"name" xorm:"not null comment('任务名称') VARCHAR(128)"`
	Status   int    `json:"status" xorm:"not null default 0 comment('状态，0为展示 1为不展示') TINYINT(1)"`
	TaskType int    `json:"task_type" xorm:"not null comment('任务类型 1.登陆应用') TINYINT(1)"`
	UserId   string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	Ycoin    int    `json:"ycoin" xorm:"not null comment('获取游币数') INT(11)"`
}
