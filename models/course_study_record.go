package models

type CourseStudyRecord struct {
	Id       int64  `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	CourseId int64  `json:"course_id" xorm:"not null comment('课程id') BIGINT(20)"`
	CreateAt int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
	Status   int    `json:"status" xorm:"not null default 0 comment('0 未购买 1 购买过') TINYINT(1)"`
}
