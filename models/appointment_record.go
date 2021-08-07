package models

type AppointmentRecord struct {
	Id              int64  `json:"id" xorm:"pk autoincr comment('id') BIGINT(20)"`
	UserId          string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	RelatedId       string `json:"related_id" xorm:"not null comment('关联id 私教/场馆/课程') index VARCHAR(60)"`
	AppointmentType int    `json:"appointment_type" xorm:"not null default 0 comment('0 场馆预约 1 私教预约') TINYINT(1)"`
	BeginNode       string `json:"begin_node" xorm:"not null default '' comment('开始时间节点 例如 10:00') VARCHAR(128)"`
	EndNode         string `json:"end_node" xorm:"not null default '' comment('结束时间节点 例如 12:00') VARCHAR(128)"`
	Date            string `json:"date" xorm:"not null default ' ' comment('预约日期') VARCHAR(30)"`
	PayOrderId      string `json:"pay_order_id" xorm:"not null comment('订单号') VARCHAR(150)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	PurchasedNum    int    `json:"purchased_num" xorm:"not null comment('购买的数量') INT(10)"`
	Status          int    `json:"status" xorm:"not null default 0 comment('0 待支付 1 订单超时/未支付 2 已支付 3 已完成 4 已取消  5 退款中 6 已退款 7 软删除') TINYINT(1)"`
}
