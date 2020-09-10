package models

type Feedback struct {
	Id       int64  `json:"id" xorm:"pk comment('自增主键') BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	Phone    string `json:"phone" xorm:"comment('手机号码') VARCHAR(200)"`
	Text     string `json:"text" xorm:"comment('反馈的内容') MEDIUMTEXT"`
	Contact  string `json:"contact" xorm:"comment('联系方式') VARCHAR(200)"`
	Status   int    `json:"status" xorm:"not null default 0 comment('状态 0未回复 1已回复') TINYINT(1)"`
	CreateAt int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	Content  string `json:"content" xorm:"comment('回复内容') MEDIUMTEXT"`
	Pics     string `json:"pics" xorm:"comment('上传的图片，多张逗号分隔') VARCHAR(512)"`
}
