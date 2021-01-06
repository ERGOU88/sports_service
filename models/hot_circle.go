package models

type HotCircle struct {
	CreateAt    int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	HotCircleId string `json:"hot_circle_id" xorm:"not null comment('热门圈子id 多个用逗号分隔 例如：3,6,11,21') VARCHAR(128)"`
	Id          int    `json:"id" xorm:"not null pk autoincr comment('自增主键') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
