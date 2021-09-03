package models

type VenueAppointmentRecordBk struct {
	Id              int64  `json:"id" xorm:"pk autoincr comment('id') BIGINT(20)"`
	UserId          string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	UseUserId       string `json:"use_user_id" xorm:"default '' comment('使用UserID') VARCHAR(255)"`
	RelatedId       int64  `json:"related_id" xorm:"not null comment('关联id 私教课程/场馆/大课') index BIGINT(20)"`
	AppointmentType int    `json:"appointment_type" xorm:"not null default 0 comment('0 场馆预约 1 私教预约') TINYINT(1)"`
	TimeNode        string `json:"time_node" xorm:"not null default '' comment('时间节点 例如 10:00-12:00') VARCHAR(128)"`
	Date            string `json:"date" xorm:"not null default ' ' comment('预约日期') VARCHAR(30)"`
	PayOrderId      string `json:"pay_order_id" xorm:"not null comment('订单号') VARCHAR(150)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	PurchasedNum    int    `json:"purchased_num" xorm:"not null comment('购买的数量') INT(10)"`
	Status          int    `json:"status" xorm:"not null default 0 comment('-1 软删除 0 待支付 1 订单超时/未支付/已取消 2 已支付 3 已完成  4 退款中 5 已退款 6 已过期 ') TINYINT(1)"`
	SeatInfo        string `json:"seat_info" xorm:"not null default '0' comment('预约的座位信息 ') VARCHAR(1000)"`
	DeductionTm     int64  `json:"deduction_tm" xorm:"not null default 0 comment('抵扣会员时长') BIGINT(20)"`
	CoachId         int64  `json:"coach_id" xorm:"not null default 0 comment('教练id 包含 私教老师/大课老师') BIGINT(20)"`
}
