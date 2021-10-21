package models

type VenueAdminRoles struct {
	Id         int    `json:"id" xorm:"INT(11)"`
	RoleName   string `json:"role_name" xorm:"VARCHAR(255)"`
	Permission string `json:"permission" xorm:"TEXT"`
	Status     int    `json:"status" xorm:"default 0 TINYINT(4)"`
	CreateAt   int    `json:"create_at" xorm:"default 0 INT(11)"`
	UpdateAt   int    `json:"update_at" xorm:"default 0 INT(11)"`
}
