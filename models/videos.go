package models

type Videos struct {
	VideoId       int64  `json:"video_id" xorm:"not null pk comment('视频id') BIGINT(20)"`
	Title         string `json:"title" xorm:"comment('视频标题') MEDIUMTEXT"`
	Describe      string `json:"describe" xorm:"comment('视频描述') MEDIUMTEXT"`
	Cover         string `json:"cover" xorm:"default '' comment('视频封面') VARCHAR(521)"`
	VideoAddr     string `json:"video_addr" xorm:"default '' comment('视频地址') VARCHAR(521)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	Sortorder     int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('状态（0：展示，1：隐藏 ）') TINYINT(1)"`
	IsRecommend   int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop         int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶）') TINYINT(1)"`
	VideoDuration int    `json:"video_duration" xorm:"not null default 0 comment('视频时长（单位：秒）') INT(8)"`
	RecContent    string `json:"rec_content" xorm:"comment('推荐理由') MEDIUMTEXT"`
	TopContent    string `json:"top_content" xorm:"comment('置顶理由') MEDIUMTEXT"`
	VideoWidth    int64  `json:"video_width" xorm:"not null default 0 comment('视频宽') BIGINT(20)"`
	VideoHeight   int64  `json:"video_height" xorm:"not null default 0 comment('视频高') BIGINT(20)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
