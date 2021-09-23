package models

type Videos struct {
	VideoId       int64  `json:"video_id" xorm:"not null pk autoincr comment('视频id') BIGINT(20)"`
	Title         string `json:"title" xorm:"comment('视频标题') MEDIUMTEXT"`
	Describe      string `json:"describe" xorm:"comment('视频描述') MEDIUMTEXT"`
	Cover         string `json:"cover" xorm:"not null default '' comment('视频封面') VARCHAR(521)"`
	VideoAddr     string `json:"video_addr" xorm:"not null default '' comment('视频地址') VARCHAR(521)"`
	UserId        string `json:"user_id" xorm:"not null default '' comment('用户id') index VARCHAR(60)"`
	UserType      int    `json:"user_type" xorm:"not null default 0 comment('添加用户类型（0：管理员[sys_user]，1：用户[user]）') TINYINT(1)"`
	Sortorder     int    `json:"sortorder" xorm:"not null default 1 comment('排序') INT(11)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('审核状态（0：审核中，1：审核通过 2：审核不通过 3：逻辑删除）') TINYINT(1)"`
	IsRecommend   int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop         int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	VideoDuration int    `json:"video_duration" xorm:"not null default 0 comment('视频时长（单位：毫秒）') INT(8)"`
	RecContent    string `json:"rec_content" xorm:"comment('推荐理由') MEDIUMTEXT"`
	TopContent    string `json:"top_content" xorm:"comment('置顶理由') MEDIUMTEXT"`
	VideoWidth    int64  `json:"video_width" xorm:"not null default 0 comment('视频宽') BIGINT(20)"`
	VideoHeight   int64  `json:"video_height" xorm:"not null default 0 comment('视频高') BIGINT(20)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
	FileId        int64  `json:"file_id" xorm:"not null default 0 comment('腾讯云文件id') BIGINT(20)"`
	Size          int64  `json:"size" xorm:"not null default 0 comment('视频大小（字节数）') BIGINT(20)"`
	PlayInfo      string `json:"play_info" xorm:"not null comment('视频转码数据') MEDIUMTEXT"`
	PubType       int    `json:"pub_type" xorm:"not null default 1 comment('1 首页发布视频 2 社区发布视频 3 管理后台发布至推荐板块') TINYINT(1)"`
	Subarea       int    `json:"subarea" xorm:"not null default 0 comment('视频所属分区') INT(11)"`
	Album         int64  `json:"album" xorm:"not null default 0 comment('视频所属专辑') BIGINT(20)"`
	AiStatus      int    `json:"ai_status" xorm:"not null default 0 comment('AI审核状态（0：未使用AI审核，1：AI审核通过 2：AI审核不通过 3：AI建议复审') TINYINT(1)"`
	SectionId     int    `json:"section_id" xorm:"not null default 0 comment('视频所属推荐板块id') INT(11)"`
}
