package models

type VenueMaterialObjectRecord struct {
	Id           int64  `json:"id" xorm:"pk autoincr comment('id') BIGINT(20)"`
	UserId       string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	ProductType  int    `json:"product_type" xorm:"not null comment('5101 实物类商品') INT(8)"`
	PayOrderId   string `json:"pay_order_id" xorm:"not null default '' comment('订单号') VARCHAR(150)"`
	PurchasedNum int    `json:"purchased_num" xorm:"not null comment('购买的数量') INT(10)"`
	ProductImg   string `json:"product_img" xorm:"not null default '' comment('商品图片') VARCHAR(1000)"`
	ProductName  string `json:"product_name" xorm:"not null comment('商品名称') VARCHAR(128)"`
	Describe     string `json:"describe" xorm:"not null default '' comment('商品介绍') VARCHAR(1000)"`
	VenueId      int64  `json:"venue_id" xorm:"not null comment('场馆id') BIGINT(20)"`
	CreateAt     int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt     int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
