package models

type Complaint struct {
	ComplaintType int    `json:"complaint_type" xorm:"not null default 0 comment('0 举报其他 1 举报视频 2 举报帖子') TINYINT(2)"`
	ComposeId     int64  `json:"compose_id" xorm:"not null comment('举报的作品id（视频/帖子id）') BIGINT(20)"`
	Content       string `json:"content" xorm:"comment('回复内容') MEDIUMTEXT"`
	Cover         string `json:"cover" xorm:"not null default '' comment('图片地址  逗号分隔') VARCHAR(512)"`
	CreateAt      int    `json:"create_at" xorm:"not null INT(11)"`
	Id            int    `json:"id" xorm:"not null pk autoincr comment('主键') INT(11)"`
	IsDispose     int    `json:"is_dispose" xorm:"default 1 comment('是否受理 1未受理 2受理') TINYINT(1)"`
	Reason        string `json:"reason" xorm:"comment('原因') MEDIUMTEXT"`
	ToUid         string `json:"to_uid" xorm:"not null comment('被举报人') VARCHAR(60)"`
	UpdateAt      int    `json:"update_at" xorm:"not null INT(11)"`
	UserId        string `json:"user_id" xorm:"not null comment('举报人') VARCHAR(60)"`
}
