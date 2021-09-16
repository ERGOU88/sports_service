package models

type VideoLiveReplay struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	UserId      string `json:"user_id" xorm:"not null default '' comment('用户id') VARCHAR(60)"`
	LiveId      int64  `json:"live_id" xorm:"not null default 0 comment('关联直播（video_live)ID ') BIGINT(20)"`
	HistoryAddr string `json:"history_addr" xorm:"not null default '' comment('回放地址') VARCHAR(512)"`
	Labeltype   int    `json:"labeltype" xorm:"not null default 0 comment('0 下架 1上架') TINYINT(1)"`
	IsDel       int    `json:"is_del" xorm:"not null default 0 comment('0 正常  1 删除') TINYINT(1)"`
	Title       string `json:"title" xorm:"not null default '' comment('标题') VARCHAR(255)"`
	PlayNum     int64  `json:"play_num" xorm:"not null default 0 comment('播放次数') BIGINT(20)"`
	IsTranscode int    `json:"is_transcode" xorm:"not null default 0 comment('是否转码,0-未申请，1-已申请') TINYINT(2)"`
	TaskId      string `json:"task_id" xorm:"not null default '' comment('回放转码任务ID') VARCHAR(64)"`
	Duration    int    `json:"duration" xorm:"not null default 0 comment('时长（秒）') INT(11)"`
	Size        int64  `json:"size" xorm:"not null default 0 comment('视频大小（字节数）') BIGINT(20)"`
	PlayInfo    string `json:"play_info" xorm:"not null comment('转码数据') MEDIUMTEXT"`
	Cover       string `json:"cover" xorm:"not null default '' comment('回放封面') VARCHAR(512)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 INT(11)"`
	Describe    string `json:"describe" xorm:"not null default '' comment('描述') VARCHAR(512)"`
	FileId      string `json:"file_id" xorm:"not null default '' comment('腾讯云文件id') VARCHAR(64)"`
}
