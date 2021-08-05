package models

type VenueRechargeConfig struct {
	Id             int    `json:"id" xorm:"not null pk autoincr comment('主键') INT(11)"`
	ProductName    string `json:"product_name" xorm:"not null comment('商品名称') VARCHAR(128)"`
	Price          int    `json:"price" xorm:"not null comment('商品价格（软妹币：分）') INT(11)"`
	ReceiveAmount  int    `json:"receive_amount" xorm:"not null comment('充值可得金额(分)') INT(11)"`
	PlatformType   int    `json:"platform_type" xorm:"not null default 0 comment('0 android端 1 iOS端') TINYINT(1)"`
	IsRecommend    int    `json:"is_recommend" xorm:"not null default 0 comment('是否推荐 0 不推荐 1 推荐') TINYINT(1)"`
	IosProductName string `json:"ios_product_name" xorm:"not null default '' comment('iOS商品名称') VARCHAR(100)"`
	CreateAt       int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
