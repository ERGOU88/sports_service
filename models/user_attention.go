package models

type UserAttention struct {
	AttentionUid string `json:"attention_uid" xorm:"not null comment('关注的用户id') index VARCHAR(60)"`
	CreateAt     int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Id           int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	Status       int    `json:"status" xorm:"not null default 1 comment('1表示关注 0表示取消关注') TINYINT(1)"`
	UserId       string `json:"user_id" xorm:"not null comment('被关注的用户id') index VARCHAR(60)"`
}
