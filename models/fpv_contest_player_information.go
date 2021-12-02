package models

type FpvContestPlayerInformation struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	Name      string `json:"name" xorm:"not null comment('选手名称 例如：pdd') VARCHAR(60)"`
	Photo     string `json:"photo" xorm:"not null default '' comment('选手照片') VARCHAR(512)"`
	Country   string `json:"country" xorm:"not null default '' comment('国家') VARCHAR(128)"`
	Province  string `json:"province" xorm:"not null default '' comment('省份') VARCHAR(128)"`
	City      string `json:"city" xorm:"not null default '' comment('城市') VARCHAR(128)"`
	Age       int    `json:"age" xorm:"not null default 0 comment('年龄') INT(3)"`
	Hobby     string `json:"hobby" xorm:"not null default '' comment('爱好') VARCHAR(255)"`
	ContestId int    `json:"contest_id" xorm:"not null default 0 comment('参加的赛事id') index INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	IdCard    string `json:"id_card" xorm:"not null default '' comment('证件号码') VARCHAR(255)"`
	IdType    int    `json:"id_type" xorm:"not null default 0 comment('1 身份证 2 居住证 3 护照 4 港澳') TINYINT(2)"`
	Gender    int    `json:"gender" xorm:"not null default 0 comment('0 未知 1 男 2 女') TINYINT(1)"`
	Born      string `json:"born" xorm:"not null default '' comment('出生年月日') VARCHAR(128)"`
	MobileNum string `json:"mobile_num" xorm:"not null default '' comment('手机号码') VARCHAR(60)"`
}
