package models

type YcoinTask struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	Name        string `json:"name" xorm:"not null comment('任务名称') VARCHAR(128)"`
	Explain     string `json:"explain" xorm:"not null default '' comment('说明') VARCHAR(256)"`
	TaskType    int    `json:"task_type" xorm:"not null comment('任务类型 1.登陆应用') TINYINT(1)"`
	Ycoin       int    `json:"ycoin" xorm:"not null comment('每次获取游币数') INT(11)"`
	Count       int    `json:"count" xorm:"not null comment('限制次数 0没有限制') INT(6)"`
	PeriodLimit int    `json:"period_limit" xorm:"not null comment('限购周期 0永久 1.1天 2.1周 3.1月 4.1年 ') TINYINT(1)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('状态，0为关闭任务 1为开启任务') TINYINT(1)"`
	TaskIcon    string `json:"task_icon" xorm:"not null default '' comment('任务图标') VARCHAR(256)"`
	Describe    string `json:"describe" xorm:"not null default '' comment('描述') VARCHAR(256)"`
	CreateAt    int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
