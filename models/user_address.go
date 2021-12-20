package models

type UserAddress struct {
	Id          int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	UserId      string `json:"user_id" xorm:"not null default '0' comment('会员id') VARCHAR(60)"`
	Name        string `json:"name" xorm:"not null default '' comment('用户姓名') VARCHAR(50)"`
	Mobile      string `json:"mobile" xorm:"not null default '' comment('手机') VARCHAR(15)"`
	Telephone   string `json:"telephone" xorm:"not null default '' comment('联系电话') VARCHAR(30)"`
	ProvinceId  int    `json:"province_id" xorm:"not null default 0 comment('省id') INT(11)"`
	CityId      int    `json:"city_id" xorm:"not null default 0 comment('市id') INT(11)"`
	DistrictId  int    `json:"district_id" xorm:"not null default 0 comment('区县id') INT(11)"`
	CommunityId int    `json:"community_id" xorm:"not null default 0 comment('社区id') INT(11)"`
	Address     string `json:"address" xorm:"not null default '' comment('地址信息') VARCHAR(255)"`
	FullAddress string `json:"full_address" xorm:"not null default '' comment('详细地址信息') VARCHAR(255)"`
	PostalCode  string `json:"postal_code" xorm:"not null default '' comment('邮编') VARCHAR(50)"`
	Longitude   string `json:"longitude" xorm:"not null default '' comment('经度') VARCHAR(255)"`
	Latitude    string `json:"latitude" xorm:"not null default '' comment('纬度') VARCHAR(255)"`
	IsDefault   int    `json:"is_default" xorm:"not null default 0 comment('是否是默认地址 0 默认') TINYINT(4)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
