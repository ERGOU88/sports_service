package models

type VenueUserCardbag struct {
	Id       int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	VenueId  int    `json:"venue_id" xorm:"comment('场馆ID') INT(11)"`
	UserId   string `json:"user_id" xorm:"comment('用户ID') VARCHAR(255)"`
	QrCode   string `json:"qr_code" xorm:"comment('ID卡号') VARCHAR(255)"`
	AdminId  int    `json:"admin_id" xorm:"default 0 comment('管理员ID') INT(11)"`
	Status   int    `json:"status" xorm:"default 0 comment('可用状态0可用、1禁用') TINYINT(4)"`
	CreateAt int    `json:"create_at" xorm:"default 0 comment('创建时间') INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"default 0 comment('修改时间') INT(11)"`
}
