package models

type VenueProductInfo struct {
	Id                int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	ProductCode       string `json:"product_code" xorm:"default '' comment('商品码') index VARCHAR(255)"`
	ProductName       string `json:"product_name" xorm:"not null comment('商品名称') VARCHAR(128)"`
	ProductType       int    `json:"product_type" xorm:"not null comment('2101 临时卡 2201 次卡 2311 月卡 2321 季卡 2331 半年卡 2341 年卡 4001 储值卡') INT(8)"`
	ProductCategory   int    `json:"product_category" xorm:"not null default 0 comment('商品类别 1000 预约类 2000 卡类 5000 实物类') INT(8)"`
	ProductChannel    int    `json:"product_channel" xorm:"not null default 0 comment('商品渠道 1 线上 2 线下') TINYINT(4)"`
	RealAmount        int    `json:"real_amount" xorm:"not null comment('真实价格（单位：分）') INT(11)"`
	CurAmount         int    `json:"cur_amount" xorm:"not null comment('当前价格 (包含真实价格、 折扣价格（单位：分）') INT(11)"`
	DiscountRate      int    `json:"discount_rate" xorm:"not null default 0 comment('折扣率') INT(11)"`
	DiscountAmount    int    `json:"discount_amount" xorm:"not null default 0 comment('优惠的金额') INT(11)"`
	EffectiveDuration int64  `json:"effective_duration" xorm:"not null default 0 comment('可用时长（秒）') BIGINT(11)"`
	ExpireDuration    int64  `json:"expire_duration" xorm:"not null default 0 comment('过期时长（秒）') BIGINT(11)"`
	Icon              string `json:"icon" xorm:"not null default '' comment('商品icon') VARCHAR(1000)"`
	Image             string `json:"image" xorm:"not null default '' comment('商品图片') VARCHAR(1000)"`
	Describe          string `json:"describe" xorm:"not null default '' comment('商品介绍') VARCHAR(1000)"`
	Instructions      string `json:"instructions" xorm:"not null default '' comment('购买须知') VARCHAR(1000)"`
	Status            int    `json:"status" xorm:"default 0 comment('状态0上架，1下架') TINYINT(4)"`
	VenueId           int64  `json:"venue_id" xorm:"not null comment('场馆id') BIGINT(20)"`
	CreateAt          int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt          int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
