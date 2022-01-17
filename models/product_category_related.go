package models

type ProductCategoryRelated struct {
	Id           int    `json:"id" xorm:"not null pk autoincr comment('自增id') INT(11)"`
	ProductId    int    `json:"product_id" xorm:"not null comment('商品id') INT(11)"`
	CategoryId   int    `json:"category_id" xorm:"not null comment('分类id') INT(11)"`
	CategoryName string `json:"category_name" xorm:"comment('分类名') VARCHAR(521)"`
	CreateAt     int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Status       int    `json:"status" xorm:"not null default 0 comment('0 正常 1 删除') TINYINT(1)"`
}
