package models

type FpvContestIntegralRules struct {
	Id              int64  `json:"id" xorm:"pk autoincr comment('规则id') BIGINT(20)"`
	ContestId       int64  `json:"contest_id" xorm:"not null comment('所属赛事id') index BIGINT(20)"`
	ScheduleId      int64  `json:"schedule_id" xorm:"not null comment('所属赛程id') BIGINT(20)"`
	Integral        int    `json:"integral" xorm:"not null default 0 comment('可得积分数') INT(11)"`
	RuleOrder       int    `json:"rule_order" xorm:"not null default 0 comment('积分规则校验顺序') TINYINT(2)"`
	MinimumRank     int    `json:"minimum_rank" xorm:"not null comment('最小排名 包含 例： 第1 可得10分 则该字段为 1') INT(8)"`
	MaximumRank     int    `json:"maximum_rank" xorm:"not null comment('最大排名 0 表示不校验  不包含 例：第1 可得10分 则该字段为 2') INT(8)"`
	RuleDescription string `json:"rule_description" xorm:"not null default '' comment('积分规则说明 例：1-100名可得10积分') VARCHAR(512)"`
	Status          int    `json:"status" xorm:"not null default 0 comment('0表示正常 1表示废弃') TINYINT(1)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
