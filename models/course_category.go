package models

type CourseCategory struct {
	CourseId string `json:"course_id" xorm:"not null pk comment('课程id') VARCHAR(60)"`
	CateId   int    `json:"cate_id" xorm:"not null pk comment('分类id') INT(11)"`
	Name     string `json:"name" xorm:"not null default '' comment('分类名') VARCHAR(64)"`
	CreateAt int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
}
