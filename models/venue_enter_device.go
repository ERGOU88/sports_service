package models

type VenueEnterDevice struct {
	Id         int    `json:"id" xorm:"not null pk autoincr INT(10)"`
	DeviceId   string `json:"device_id" xorm:"VARCHAR(255)"`
	DeviceName string `json:"device_name" xorm:"VARCHAR(255)"`
	VenueId    int64  `json:"venue_id" xorm:"BIGINT(11)"`
	Position   string `json:"position" xorm:"VARCHAR(255)"`
	Status     int    `json:"status" xorm:"TINYINT(1)"`
	CreateAt   int    `json:"create_at" xorm:"INT(11)"`
	UpdateAt   int    `json:"update_at" xorm:"INT(11)"`
}
