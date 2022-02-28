package models

type ProductService struct {
	ServiceId int64 `json:"service_id" xorm:"not null pk comment('服务id') BIGINT(20)"`
	ProductId int   `json:"product_id" xorm:"not null pk comment('商品id') INT(11)"`
	CreateAt  int   `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
}
