package models

type VenueUserCardbag struct {
	Id       int    `json:"id" xorm:"INT(11)"`
	UserId   string `json:"user_id" xorm:"comment('用户ID') VARCHAR(255)"`
	QrCode   string `json:"qr_code" xorm:"comment('ID卡号') VARCHAR(255)"`
	CardType int    `json:"card_type" xorm:"comment('卡类型:1、年卡、2半年卡、3季卡、4月卡、5、此卡') TINYINT(4)"`
	AdminId  int    `json:"admin_id" xorm:"default 0 comment('管理员ID') INT(11)"`
	Status   int    `json:"status" xorm:"default 0 comment('可用状态0可用、1禁用') TINYINT(4)"`
	CreateAt int    `json:"create_at" xorm:"default 0 comment('创建时间') INT(11)"`
	Duration int    `json:"duration" xorm:"default 0 comment('可用时长') INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"default 0 comment('修改时间') INT(11)"`
	DtStart  int    `json:"dt_start" xorm:"default 0 comment('启用时间') INT(11)"`
	DtEnd    int    `json:"dt_end" xorm:"default 0 comment('过期时间') INT(11)"`
}
