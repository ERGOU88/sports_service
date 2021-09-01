package models

type VenueUserCardbagTemp struct {
	Id      int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	VenueId int    `json:"venue_id" xorm:"comment('场馆ID') INT(11)"`
	UserId  string `json:"user_id" xorm:"comment('用户ID') index VARCHAR(255)"`
	Type    int    `json:"type" xorm:"comment('2101 临时卡
2201 次卡') INT(11)"`
	Duration int `json:"duration" xorm:"comment('可以时长') INT(11)"`
	DtStart  int `json:"dt_start" xorm:"comment('开始时间') INT(11)"`
	DtEnd    int `json:"dt_end" xorm:"comment('结束时间') INT(11)"`
	Status   int `json:"status" xorm:"comment('状态0正常,1禁用') INT(11)"`
	CreateAt int `json:"create_at" xorm:"INT(11)"`
	UpdateAt int `json:"update_at" xorm:"INT(11)"`
}
