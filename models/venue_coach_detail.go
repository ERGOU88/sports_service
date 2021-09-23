package models

type VenueCoachDetail struct {
	Id               int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	Title            string `json:"title" xorm:"comment('教练职称') MEDIUMTEXT"`
	Name             string `json:"name" xorm:"not null default '' comment('私教名称') VARCHAR(60)"`
	Address          string `json:"address" xorm:"not null comment('私教地点') VARCHAR(128)"`
	Designation      string `json:"designation" xorm:"not null default '' comment('认证称号') VARCHAR(60)"`
	Describe         string `json:"describe" xorm:"comment('描述') MEDIUMTEXT"`
	AreasOfExpertise string `json:"areas_of_expertise" xorm:"not null default '' comment('擅长领域') VARCHAR(512)"`
	Cover            string `json:"cover" xorm:"not null default '' comment('封面 ') VARCHAR(256)"`
	Avatar           string `json:"avatar" xorm:"not null default '' comment('头像') VARCHAR(256)"`
	CreateAt         int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt         int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Status           int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	VenueId          int64  `json:"venue_id" xorm:"not null default 0 comment('场馆id') BIGINT(20)"`
}
