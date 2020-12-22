package models

type VideosExamine struct {
	Cover         string `json:"cover" xorm:"not null default '' comment('视频封面') VARCHAR(521)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	Describe      string `json:"describe" xorm:"comment('视频描述') MEDIUMTEXT"`
	IsRecommend   int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop         int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	LabelId       string `json:"label_id" xorm:"default '' comment('标签id') index VARCHAR(521)"`
	LabelName     string `json:"label_name" xorm:"comment('标签') index MEDIUMTEXT"`
	Manager       int    `json:"manager" xorm:"not null default 0 comment('后台操作用户') INT(11)"`
	RecContent    string `json:"rec_content" xorm:"comment('推荐理由') MEDIUMTEXT"`
	Sortorder     int    `json:"sortorder" xorm:"not null default 1 comment('排序') INT(11)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('审核状态（0：无操作，1：审核通过 2：审核不通过 3：删除）') TINYINT(1)"`
	Title         string `json:"title" xorm:"comment('视频标题') MEDIUMTEXT"`
	TopContent    string `json:"top_content" xorm:"comment('置顶理由') MEDIUMTEXT"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	UserType      int    `json:"user_type" xorm:"not null default 0 comment('添加用户类型（0：管理员[sys_user]，1：用户[user]）') TINYINT(1)"`
	VideoAddr     string `json:"video_addr" xorm:"not null default '' comment('视频地址') VARCHAR(521)"`
	VideoDuration int    `json:"video_duration" xorm:"not null default 0 comment('视频时长（单位：秒）') INT(8)"`
	VideoHeight   int64  `json:"video_height" xorm:"not null default 0 comment('视频高') BIGINT(20)"`
	VideoId       int64  `json:"video_id" xorm:"not null pk autoincr comment('视频id') BIGINT(20)"`
	VideoWidth    int64  `json:"video_width" xorm:"not null default 0 comment('视频宽') BIGINT(20)"`
}
