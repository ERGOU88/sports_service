package models

type AppVersionControl struct {
	Id          int64  `json:"id" xorm:"pk autoincr comment('自增ID') BIGINT(20)"`
	VersionName string `json:"version_name" xorm:"not null default '' comment('版本名称') VARCHAR(128)"`
	Version     string `json:"version" xorm:"not null default '' comment('版本号') VARCHAR(60)"`
	VersionCode int    `json:"version_code" xorm:"not null default 0 comment('版本code') index INT(8)"`
	Size        string `json:"size" xorm:"not null default '' comment('包大小') VARCHAR(128)"`
	IsForce     int    `json:"is_force" xorm:"not null default 0 comment('0 不需要强更 1 需要强更') TINYINT(1)"`
	Status      int    `json:"status" xorm:"not null default 0 comment('0 可用 1 不可用') TINYINT(1)"`
	Platform    int    `json:"platform" xorm:"not null default 0 comment('0 android 1 ios') TINYINT(1)"`
	UpgradeUrl  string `json:"upgrade_url" xorm:"not null default '' comment('更新包地址') VARCHAR(256)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	Describe    string `json:"describe" xorm:"not null default '' comment('版本说明') VARCHAR(500)"`
}
