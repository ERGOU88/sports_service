package models

type VenueRefundRules struct {
	Id              int    `json:"id" xorm:"not null pk autoincr comment('自增ID') INT(10)"`
	RuleOrder       int    `json:"rule_order" xorm:"not null default 0 comment('规则校验顺序') TINYINT(2)"`
	RuleName        string `json:"rule_name" xorm:"not null comment('规则名称') VARCHAR(255)"`
	RuleMinDuration int    `json:"rule_min_duration" xorm:"not null default 0 comment('规则校验最小时长') INT(11)"`
	RuleMaxDuration int    `json:"rule_max_duration" xorm:"not null default 0 comment('规则校验最大时长') INT(11)"`
	FeeRate         int    `json:"fee_rate" xorm:"not null default 0 comment('手续费比例 例如：1.55% 则 值为155 [乘以100存储]') INT(5)"`
	MinimumCharge   int    `json:"minimum_charge" xorm:"not null default 0 comment('最低收取手续费金额 单位[分]') INT(10)"`
	Status          int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	FeeRateCn       string `json:"fee_rate_cn" xorm:"not null default '' comment('手续费比例 中文 例如：票价5%') VARCHAR(255)"`
}
