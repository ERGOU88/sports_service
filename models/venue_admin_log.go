package models

type VenueAdminLog struct {
	Id        int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	AdminId   int    `json:"admin_id" xorm:"INT(11)"`
	Operation string `json:"operation" xorm:"VARCHAR(255)"`
	Log       string `json:"log" xorm:"TEXT"`
	Ip        string `json:"ip" xorm:"VARCHAR(255)"`
	CreateAt  int    `json:"create_at" xorm:"INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"INT(11)"`
}
