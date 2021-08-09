package models

type VenueInfo struct {
	Id            int64  `json:"id" xorm:"pk autoincr comment('场馆ID') BIGINT(20)"`
	VenueName     string `json:"venue_name" xorm:"not null comment('场馆名称') VARCHAR(60)"`
	Address       string `json:"address" xorm:"not null comment('场馆地址') VARCHAR(100)"`
	Describe      string `json:"describe" xorm:"not null comment('场馆介绍') VARCHAR(1000)"`
	Telephone     string `json:"telephone" xorm:"not null comment('联系电话') VARCHAR(60)"`
	VenueImages   string `json:"venue_images" xorm:"comment('场馆图片 多张逗号分隔') TEXT"`
	BusinessHours string `json:"business_hours" xorm:"not null comment('营业时间') VARCHAR(100)"`
	Services      string `json:"services" xorm:"not null comment('设施及服务') VARCHAR(300)"`
	Longitude     string `json:"longitude" xorm:"not null default '' comment('经度') VARCHAR(30)"`
	Latitude      string `json:"latitude" xorm:"not null default '' comment('纬度') VARCHAR(30)"`
	Status        int    `json:"status" xorm:"not null default 0 comment('0 正常 1 废弃') TINYINT(1)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	ServiceStatus int    `json:"service_status" xorm:"not null default 0 comment('0 正常营业 1 暂停营业') TINYINT(1)"`
}
