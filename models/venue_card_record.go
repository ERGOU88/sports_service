package models

type VenueCardRecord struct {
	Id             int64  `json:"id" xorm:"pk autoincr comment('id') BIGINT(20)"`
	UserId         string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	UseUserId      string `json:"use_user_id" xorm:"default '' comment('使用者UserID') VARCHAR(60)"`
	ProductType    int    `json:"product_type" xorm:"not null comment('2101 临时卡 2201 次卡 2311 月卡 2321 季卡 2331 半年卡 2341 年卡') INT(8)"`
	PayOrderId     string `json:"pay_order_id" xorm:"not null default '' comment('订单号') VARCHAR(150)"`
	PurchasedNum   int    `json:"purchased_num" xorm:"not null comment('购买的数量') INT(10)"`
	Status         int    `json:"status" xorm:"not null default 0 comment('0 不可用 1 可用') TINYINT(1)"`
	SingleDuration int    `json:"single_duration" xorm:"not null default 0 comment('单个时长') INT(11)"`
	ExpireDuration int    `json:"expire_duration" xorm:"not null default 0 comment('过期时长[单个商品]') INT(11)"`
	Duration       int    `json:"duration" xorm:"not null default 0 comment('购买总时长') INT(11)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
