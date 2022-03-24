package models

type CourseVideos struct {
	Id            int64  `json:"id" xorm:"pk autoincr comment('自增文件id') BIGINT(20)"`
	CourseId      int64  `json:"course_id" xorm:"not null comment('课程id') index BIGINT(20)"`
	Title         string `json:"title" xorm:"not null comment('课程视频标题') VARCHAR(521)"`
	IsFree        int    `json:"is_free" xorm:"not null default 0 comment('0 收费 1 免费') TINYINT(1)"`
	VipIsFree     int    `json:"vip_is_free" xorm:"not null default 0 comment('会员是否免费 0 收费 1 免费') TINYINT(1)"`
	Cover         string `json:"cover" xorm:"not null default '' comment('视频封面') VARCHAR(521)"`
	VideoAddr     string `json:"video_addr" xorm:"not null default '' comment('视频地址') VARCHAR(521)"`
	FileOrder     int    `json:"file_order" xorm:"not null default 1 comment('目录文件顺序 正序') INT(11)"`
	VideoDuration int    `json:"video_duration" xorm:"not null default 0 comment('视频时长（单位：毫秒）') INT(8)"`
	VideoWidth    int64  `json:"video_width" xorm:"not null default 0 comment('视频宽') BIGINT(20)"`
	VideoHeight   int64  `json:"video_height" xorm:"not null default 0 comment('视频高') BIGINT(20)"`
	TxFileId      string `json:"tx_file_id" xorm:"not null default '0' comment('腾讯云文件id') VARCHAR(128)"`
	Size          int64  `json:"size" xorm:"not null default 0 comment('视频大小（字节数）') BIGINT(20)"`
	PlayInfo      string `json:"play_info" xorm:"comment('视频转码信息') MEDIUMTEXT"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('0 正常 1 待发布 2 剔除') TINYINT(1)"`
}
