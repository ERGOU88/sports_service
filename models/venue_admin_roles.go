package models

type VenueAdminRoles struct {
	Id         int64  `json:"id" xorm:"pk autoincr BIGINT(11)"`
	Title      string `json:"title" xorm:"VARCHAR(255)"`
	RoleName   string `json:"role_name" xorm:"VARCHAR(255)"`
	Permission string `json:"permission" xorm:"TEXT"`
	Status     int    `json:"status" xorm:"default 0 TINYINT(4)"`
	VenueId    int    `json:"venue_id" xorm:"default 0 INT(10)"`
	CreateAt   int    `json:"create_at" xorm:"default 0 INT(11)"`
	UpdateAt   int    `json:"update_at" xorm:"default 0 INT(11)"`
}
