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
}
