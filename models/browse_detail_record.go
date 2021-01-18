package models

type BrowseDetailRecord struct {
	Id            int64  `json:"id" xorm:"pk autoincr comment('id') BIGINT(20)"`
	VideoId       int64  `json:"video_id" xorm:"not null comment('视频id') index BIGINT(20)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	StudyDuration int    `json:"study_duration" xorm:"not null comment('当前播放的时长（秒）') INT(8)"`
	TotalDuration int    `json:"total_duration" xorm:"not null comment('当前视频总时长（秒）') INT(8)"`
	CurProgress   string `json:"cur_progress" xorm:"not null default 0.00 comment('当前播放的进度') DECIMAL(5,2)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
