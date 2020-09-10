package models

type PostInfo struct {
	PostId        int64  `json:"post_id" xorm:"not null pk comment('帖子id') BIGINT(20)"`
	Title         string `json:"title" xorm:"comment('帖子标题') MEDIUMTEXT"`
	Describe      string `json:"describe" xorm:"comment('帖子描述') MEDIUMTEXT"`
	Cover         string `json:"cover" xorm:"not null default '' comment('封面') VARCHAR(521)"`
	VideoAddr     string `json:"video_addr" xorm:"not null default '' comment('视频地址') VARCHAR(521)"`
	CircleId      int    `json:"circle_id" xorm:"not null comment('圈子id') INT(11)"`
	SectionId     int    `json:"section_id" xorm:"not null comment('板块id') INT(11)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	UserType      int    `json:"user_type" xorm:"not null default 0 comment('添加用户类型（0：管理员[sys_user]，1：用户[user]）') TINYINT(1)"`
	Sortorder     int    `json:"sortorder" xorm:"not null default 0 comment('权重') INT(11)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('状态（0：展示 1：隐藏 2：删除 预留状态）') TINYINT(1)"`
	IsRecommend   int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop         int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	VideoDuration int    `json:"video_duration" xorm:"not null default 0 comment('视频时长（单位：秒）') INT(8)"`
	RecContent    string `json:"rec_content" xorm:"comment('推荐理由') MEDIUMTEXT"`
	TopContent    string `json:"top_content" xorm:"comment('置顶理由') MEDIUMTEXT"`
	VideoWidth    int64  `json:"video_width" xorm:"not null default 0 comment('视频宽') BIGINT(20)"`
	VideoHeight   int64  `json:"video_height" xorm:"not null default 0 comment('视频高') BIGINT(20)"`
}
