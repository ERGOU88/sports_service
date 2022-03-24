package models

type UserStudyRecord struct {
	Id            int64  `json:"id" xorm:"pk autoincr comment('自增id') BIGINT(20)"`
	FileId        int64  `json:"file_id" xorm:"not null comment('课程视频文件id') index BIGINT(20)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	StudyDuration int    `json:"study_duration" xorm:"not null comment('学习当前视频时长(真实停留时长)') INT(8)"`
	TotalDuration int    `json:"total_duration" xorm:"not null comment('当前课程视频总时长') INT(8)"`
	CurProgress   string `json:"cur_progress" xorm:"not null default 0.00 comment('当前视频进度') DECIMAL(5,2)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	CourseId      int64  `json:"course_id" xorm:"not null comment('课程id') BIGINT(20)"`
	PlayDuration  int    `json:"play_duration" xorm:"not null comment('视频已观看的时长（用户可能直接拉取进度条）') INT(8)"`
}
