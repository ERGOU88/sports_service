package models

type VenueProductInfo struct {
	Id                int    `json:"id" xorm:"not null pk autoincr comment('主键') INT(11)"`
	ProductName       string `json:"product_name" xorm:"not null comment('商品名称') VARCHAR(128)"`
	Price             int    `json:"price" xorm:"not null comment('商品价格（分）') INT(11)"`
	ProductType       int    `json:"product_type" xorm:"not null comment('1 次卡 2 购买月卡 3 购买季卡 4 购买年卡 5 体验券 ') TINYINT(2)"`
	EffectiveDuration int    `json:"effective_duration" xorm:"not null default 0 comment('有效时长') INT(11)"`
	RealAmount        int    `json:"real_amount" xorm:"not null comment('真实价格（单位：分）') INT(11)"`
	CurAmount         int    `json:"cur_amount" xorm:"not null comment('当前价格 (包含真实价格、 折扣价格（单位：分）') INT(11)"`
	DiscountRate      int    `json:"discount_rate" xorm:"not null default 0 comment('折扣率') INT(11)"`
	DiscountAmount    int    `json:"discount_amount" xorm:"not null default 0 comment('优惠的金额') INT(11)"`
	CreateAt          int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt          int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
