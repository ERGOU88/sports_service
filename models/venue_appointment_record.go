package models

type VenueAppointmentRecord struct {
	Id              int64  `json:"id" xorm:"pk autoincr comment('id') BIGINT(20)"`
	UserId          string `json:"user_id" xorm:"not null comment('用户id') VARCHAR(60)"`
	UseUserId       string `json:"use_user_id" xorm:"default '' comment('使用者UserID') VARCHAR(60)"`
	IsWriteOff      int    `json:"is_write_off" xorm:"not null default 0 comment('是否核销 0 未核销 1 已核销') TINYINT(1)"`
	Duration        int    `json:"duration" xorm:"not null default 0 comment('购买相关服务总时长') INT(11)"`
	AppointmentType int    `json:"appointment_type" xorm:"not null default 0 comment('0 场馆预约 1 私教预约 2 大课预约') TINYINT(1)"`
	TimeNode        string `json:"time_node" xorm:"not null default '' comment('时间节点 例如 10:00-12:00') VARCHAR(128)"`
	Date            string `json:"date" xorm:"not null default ' ' comment('预约日期 例如 2021-10-01') VARCHAR(30)"`
	PayOrderId      string `json:"pay_order_id" xorm:"not null default '' comment('订单号') VARCHAR(150)"`
	PurchasedNum    int    `json:"purchased_num" xorm:"not null comment('购买的数量') INT(10)"`
	SeatInfo        string `json:"seat_info" xorm:"not null default '' comment('预约的座位信息 ') VARCHAR(1000)"`
	DeductionTm     int64  `json:"deduction_tm" xorm:"not null default 0 comment('抵扣会员时长') BIGINT(20)"`
	DeductionAmount int64  `json:"deduction_amount" xorm:"not null default 0 comment('抵扣的价格') BIGINT(20)"`
	DeductionNum    int64  `json:"deduction_num" xorm:"not null default 0 comment('抵扣的数量') BIGINT(20)"`
	SingleDuration  int    `json:"single_duration" xorm:"not null default 0 comment('单个时长') INT(11)"`
	CoachId         int64  `json:"coach_id" xorm:"not null default 0 comment('教练id 包含 私教老师/大课老师') BIGINT(20)"`
	CreateAt        int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt        int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	UnitPrice       int    `json:"unit_price" xorm:"not null default 0 comment('单位价格 [分] tips:[用于计费使用] 例：每十五分钟/10元 则该字段为1000') INT(11)"`
	UnitDuration    int    `json:"unit_duration" xorm:"not null default 0 comment('单位时长 [秒] tips:[用于计费使用] 例：每十五分钟/10元 则该字段为900') INT(11)"`
	StartTm         int    `json:"start_tm" xorm:"not null default 0 comment('预约开始时间') INT(11)"`
	EndTm           int    `json:"end_tm" xorm:"not null default 0 comment('预约结束时间') INT(11)"`
	VenueId         int64  `json:"venue_id" xorm:"not null default 0 comment('场馆id') BIGINT(20)"`
	CourseId        int64  `json:"course_id" xorm:"not null default 0 comment('课程id') BIGINT(20)"`
}
