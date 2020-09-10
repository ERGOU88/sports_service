package models

type PostCommentParentChild struct {
	Id       int64 `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	ParentId int64 `json:"parent_id" xorm:"not null comment('评论父id') index BIGINT(20)"`
	ChildId  int64 `json:"child_id" xorm:"not null comment('评论子id') index BIGINT(20)"`
	CreateAt int   `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
}
