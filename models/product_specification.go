package models

type ProductSpecification struct {
	CategoryId     int64  `json:"category_id" xorm:"not null pk comment('规格模板所属商品分类id') BIGINT(20)"`
	Specifications string `json:"specifications" xorm:"not null default '' comment('规格参数模板，json格式') VARCHAR(3000)"`
	CreateAt       int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt       int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
