package models

type CourseStatistic struct {
	Id       int64 `json:"id" xorm:"pk autoincr comment('课程id') index BIGINT(20)"`
	CourseId int64 `json:"course_id" xorm:"not null comment('课程id') BIGINT(20)"`
	StudyNum int   `json:"study_num" xorm:"not null default 0 comment('学习人数') INT(11)"`
	CreateAt int   `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt int   `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
