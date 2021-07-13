package models

type TencentCloudEvents struct {
	Id        int    `json:"id" xorm:"not null pk autoincr comment('主键') INT(11)"`
	FileId    int64  `json:"file_id" xorm:"not null default 0 comment('腾讯云文件id') BIGINT(20)"`
	Event     string `json:"event" xorm:"comment('事件内容（json字符串）') MEDIUMTEXT"`
	CreateAt  int    `json:"create_at" xorm:"not null INT(11)"`
	EventType int    `json:"event_type" xorm:"not null default 0 comment('0 上传事件 1 视频转码事件') TINYINT(3)"`
}
