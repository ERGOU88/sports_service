package models

type VenueEnterLog struct {
	Id            int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	UserId        string `json:"user_id" xorm:"comment('用户user id') VARCHAR(255)"`
	AppointmentId int64  `json:"appointment_id" xorm:"default 0 comment('预约ID') BIGINT(11)"`
	OrderId       string `json:"order_id" xorm:"comment('订单号') VARCHAR(255)"`
	VenueId       int64  `json:"venue_id" xorm:"comment('场馆ID') BIGINT(20)"`
	Type          int    `json:"type" xorm:"TINYINT(4)"`
	Duration      int    `json:"duration" xorm:"comment('体验时长') INT(11)"`
	DtEnter       int    `json:"dt_enter" xorm:"comment('入场时间') INT(11)"`
	DtExit        int    `json:"dt_exit" xorm:"comment('出场时间') INT(11)"`
	DeviceEnter   string `json:"device_enter" xorm:"comment('入场设备号') VARCHAR(255)"`
	DeviceExit    string `json:"device_exit" xorm:"comment('出场设备号') VARCHAR(255)"`
	CreateAt      int    `json:"create_at" xorm:"comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"comment('更新时间') INT(11)"`
	UnitPrice     int    `json:"unit_price" xorm:"INT(10)"`
	UnitDuration  int    `json:"unit_duration" xorm:"INT(255)"`
	Amount        int    `json:"amount" xorm:"INT(255)"`
	HasSettle     int    `json:"has_settle" xorm:"default 0 comment('是否结算0未结算 1自动结算 2人工结算') TINYINT(4)"`
}