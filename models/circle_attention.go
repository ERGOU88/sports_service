package models

type CircleAttention struct {
	CircleId int    `json:"circle_id" xorm:"not null comment('圈子id') index INT(11)"`
	CreateAt int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Id       int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	Status   int    `json:"status" xorm:"not null default 1 comment('1表示关注 2表示取消关注') TINYINT(1)"`
	UserId   string `json:"user_id" xorm:"not null comment('关注圈子的用户id') index VARCHAR(60)"`
}
