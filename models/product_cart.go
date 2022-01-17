package models

type ProductCart struct {
	Id        int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	UserId    string `json:"user_id" xorm:"not null default '0' comment('用户id') index VARCHAR(60)"`
	SkuId     int    `json:"sku_id" xorm:"not null default 0 comment('sku id') INT(11)"`
	Count     int    `json:"count" xorm:"not null default 0 comment('数量') INT(10)"`
	IsCheck   int    `json:"is_check" xorm:"not null default 0 comment('0选中 1未选中') TINYINT(1)"`
	ProductId int    `json:"product_id" xorm:"not null default 0 comment('商品id') INT(11)"`
	CreateAt  int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 有效 1 无效') TINYINT(1)"`
}
