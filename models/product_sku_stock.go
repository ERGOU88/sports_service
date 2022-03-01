package models

type ProductSkuStock struct {
	SkuId        int `json:"sku_id" xorm:"not null pk comment('sku id') INT(11)"`
	Stock        int `json:"stock" xorm:"not null comment('库存数量') INT(11)"`
	MaxBuy       int `json:"max_buy" xorm:"not null default 0 comment('限购 0 表示无限制') INT(11)"`
	MinBuy       int `json:"min_buy" xorm:"not null default 0 comment('起购数') INT(11)"`
	CreateAt     int `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt     int `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	ProductId    int `json:"product_id" xorm:"not null comment('商品id') INT(11)"`
	PurchasedNum int `json:"purchased_num" xorm:"not null default 0 comment('已购买数量 [冻结库存]') INT(11)"`
}
