package models

type PostingInfo struct {
	Id            int64  `json:"id" xorm:"pk comment('帖子id') BIGINT(20)"`
	CommunityId   int    `json:"community_id" xorm:"not null comment('帖子所属社区id') index INT(11)"`
	Title         string `json:"title" xorm:"comment('帖子标题') MEDIUMTEXT"`
	Describe      string `json:"describe" xorm:"comment('帖子描述') MEDIUMTEXT"`
	Content       string `json:"content" xorm:"comment('帖子内容[富文本]') MEDIUMTEXT"`
	Cover         string `json:"cover" xorm:"not null default '' comment('封面') VARCHAR(521)"`
	VideoAddr     string `json:"video_addr" xorm:"not null default '' comment('视频地址') VARCHAR(521)"`
	VideoDuration int    `json:"video_duration" xorm:"not null default 0 comment('视频时长（单位：毫秒）') INT(8)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	PostingType   int    `json:"posting_type" xorm:"not null default 0 comment('帖子类型  0 纯文本 1 图文 2 视频 + 文字 ') TINYINT(1)"`
	ContentType   int    `json:"content_type" xorm:"not null default 0 comment('内容类型 0 发布 1 转发') TINYINT(1)"`
	ImagesAddr    string `json:"images_addr" xorm:"not null default '' comment('图片地址  支持多张') VARCHAR(521)"`
	Sortorder     int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('状态 审核状态（0：审核中，1：审核通过 2：审核不通过 3：逻辑删除）') TINYINT(1)"`
	IsRecommend   int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop         int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
