package models

type PostingInfo struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('帖子id') BIGINT(20)"`
	SectionId   int    `json:"section_id" xorm:"not null comment('版块id') INT(11)"`
	Title       string `json:"title" xorm:"comment('帖子标题') MEDIUMTEXT"`
	Describe    string `json:"describe" xorm:"comment('帖子描述') MEDIUMTEXT"`
	Content     string `json:"content" xorm:"comment('帖子内容 图片列表/json结构 例如转发的视频 完整结构') MEDIUMTEXT"`
	VideoId     int64  `json:"video_id" xorm:"not null comment('关联的视频id') BIGINT(20)"`
	UserId      string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	PostingType int    `json:"posting_type" xorm:"not null default 0 comment('帖子类型  0 纯文本 1 图文 2 视频 + 文字 ') TINYINT(1)"`
	ContentType int    `json:"content_type" xorm:"not null default 0 comment('内容类型 0 发布 1 转发视频 2 转发帖子') TINYINT(1)"`
	Sortorder   int    `json:"sortorder" xorm:"not null default 0 comment('排序权重') INT(11)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('状态 审核状态（0：审核中，1：审核通过 2：审核不通过 3：逻辑删除）') TINYINT(1)"`
	IsRecommend int    `json:"is_recommend" xorm:"not null default 0 comment('推荐（0：不推荐；1：推荐）') TINYINT(1)"`
	IsTop       int    `json:"is_top" xorm:"not null default 0 comment('置顶（0：不置顶；1：置顶；）') TINYINT(1)"`
	IsCream     int    `json:"is_cream" xorm:"not null default 0 comment('是否精华帖（0: 不是 1: 是）') TINYINT(1)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 comment('修改时间') INT(11)"`
}
