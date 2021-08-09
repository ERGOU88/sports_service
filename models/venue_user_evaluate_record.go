package models

type VenueUserEvaluateRecord struct {
	Id        int64  `json:"id" xorm:"pk autoincr comment('自增主键') BIGINT(20)"`
	UserId    string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	CoachId   string `json:"coach_id" xorm:"not null comment('教练id') index VARCHAR(60)"`
	Star      int    `json:"star" xorm:"not null comment('评价几颗星 1~5星') TINYINT(2)"`
	OrderId   string `json:"order_id" xorm:"not null comment('订单id') VARCHAR(256)"`
	OrderType int    `json:"order_type" xorm:"not null default 1 comment('订单类型 默认 1 私教课程') TINYINT(1)"`
	Content   string `json:"content" xorm:"default '' comment('评价描述') VARCHAR(1000)"`
	LabelInfo string `json:"label_info" xorm:"not null default '' comment('用户选取的评价标签信息') VARCHAR(1000)"`
	Status    int    `json:"status" xorm:"not null default 0 comment('0 有效 1 废弃') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"not null comment('更新时间') INT(11)"`
}
