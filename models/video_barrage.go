package models

type VideoBarrage struct {
	BarrageType      int    `json:"barrage_type" xorm:"not null default 0 comment('预留字段') TINYINT(2)"`
	Color            string `json:"color" xorm:"not null default '' comment('弹幕字体颜色') VARCHAR(100)"`
	Content          string `json:"content" xorm:"not null default '' comment('弹幕内容') VARCHAR(512)"`
	Font             string `json:"font" xorm:"not null default '' comment('弹幕字体') VARCHAR(100)"`
	Id               int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Location         int    `json:"location" xorm:"not null default 0 comment('弹幕位置') TINYINT(2)"`
	SendTime         int64  `json:"send_time" xorm:"not null default 0 comment('弹幕发送时间') BIGINT(20)"`
	UserId           string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	VideoCurDuration int    `json:"video_cur_duration" xorm:"not null comment('视频当前时长节点（单位：秒）') INT(8)"`
	VideoId          int64  `json:"video_id" xorm:"not null comment('视频id') index BIGINT(20)"`
}
