package models

type VenuePaymentChannel struct {
	Id         int    `json:"id" xorm:"not null pk autoincr INT(10)"`
	Title      string `json:"title" xorm:"comment('名称') VARCHAR(255)"`
	Icon       string `json:"icon" xorm:"comment('icon') VARCHAR(255)"`
	Background string `json:"background" xorm:"VARCHAR(255)"`
	Status     int    `json:"status" xorm:"comment('状态 0 可用，1禁用') TINYINT(4)"`
	Sort       int    `json:"sort" xorm:"comment('排序') TINYINT(4)"`
	Desc       string `json:"desc" xorm:"VARCHAR(255)"`
	Rate       int    `json:"rate" xorm:"INT(255)"`
	AppId      string `json:"app_id" xorm:"VARCHAR(255)"`
	AppKey     string `json:"app_key" xorm:"VARCHAR(255)"`
	AppSecret  string `json:"app_secret" xorm:"VARCHAR(255)"`
	PublicKey  string `json:"public_key" xorm:"TEXT"`
	PrivateKey string `json:"private_key" xorm:"TEXT"`
	CreateAt   int    `json:"create_at" xorm:"INT(11)"`
	UpdateAt   int    `json:"update_at" xorm:"INT(11)"`
	Identifier string `json:"identifier" xorm:"comment('支付渠道标识') VARCHAR(255)"`
}
