package models

type VenueProductInfo struct {
	Id                int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	ProductName       string `json:"product_name" xorm:"not null comment('商品名称') VARCHAR(128)"`
	ProductType       int    `json:"product_type" xorm:"not null comment('类型 2001 月卡 2002 季卡 2003 年卡 2004 体验券 3001 储值卡') INT(8)"`
	RealAmount        int    `json:"real_amount" xorm:"not null comment('真实价格（单位：分）') INT(11)"`
	CurAmount         int    `json:"cur_amount" xorm:"not null comment('当前价格 (包含真实价格、 折扣价格（单位：分）') INT(11)"`
	DiscountRate      int    `json:"discount_rate" xorm:"not null default 0 comment('折扣率') INT(11)"`
	DiscountAmount    int    `json:"discount_amount" xorm:"not null default 0 comment('优惠的金额') INT(11)"`
	VenueId           int64  `json:"venue_id" xorm:"not null comment('场馆id') BIGINT(20)"`
	CreateAt          int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt          int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
	EffectiveDuration int    `json:"effective_duration" xorm:"not null default 0 comment('可用时长（秒）') INT(11)"`
	ExpireDuration    int    `json:"expire_duration" xorm:"not null default 0 comment('过期时长（秒）') INT(11)"`
	Icon              string `json:"icon" xorm:"not null default '' comment('商品icon') VARCHAR(500)"`
	Image             string `json:"image" xorm:"not null default '' comment('商品图片') VARCHAR(1000)"`
	Describe          string `json:"describe" xorm:"not null default '' comment('商品介绍') VARCHAR(1000)"`
	Title             string `json:"title" xorm:"not null default '' comment('简介') VARCHAR(300)"`
	InstanceType      int    `json:"instance_type" xorm:"default 1 comment('实例类型，1: 体验卡；2: 线下食品') TINYINT(4)"`
}
