package models

type FpvContestInfo struct {
	Id          int    `json:"id" xorm:"not null pk autoincr comment('赛事id') INT(11)"`
	ContestName string `json:"contest_name" xorm:"not null comment('赛事名称 例如：x-fly fpv大赛s1') VARCHAR(128)"`
	Organizer   string `json:"organizer" xorm:"not null default '' comment('举报方') VARCHAR(255)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('赛事状态 0 未开始 1 已开始 2 已结束') TINYINT(2)"`
	StartTm     int    `json:"start_tm" xorm:"not null default 0 comment('赛事开始时间') INT(11)"`
	EndTm       int    `json:"end_tm" xorm:"not null default 0 comment('赛事结束时间') INT(11)"`
	SignUpNum   int    `json:"sign_up_num" xorm:"not null default 0 comment('可报名人数 0 表示没有限制') INT(8)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
