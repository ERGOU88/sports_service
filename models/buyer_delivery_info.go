package models

type BuyerDeliveryInfo struct {
	Id          int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	OrderId     string `json:"order_id" xorm:"not null default '' comment('关联订单编号') VARCHAR(50)"`
	UserId      string `json:"user_id" xorm:"not null default '' comment('买家uid') VARCHAR(60)"`
	Name        string `json:"name" xorm:"not null default '' comment('买家姓名') VARCHAR(50)"`
	Mobile      string `json:"mobile" xorm:"not null default '' comment('买家手机') VARCHAR(15)"`
	Telephone   string `json:"telephone" xorm:"not null default '' comment('买家固定电话') VARCHAR(30)"`
	ProvinceId  int    `json:"province_id" xorm:"not null default 0 comment('买家省id') INT(11)"`
	CityId      int    `json:"city_id" xorm:"not null default 0 comment('买家市id') INT(11)"`
	DistrictId  int    `json:"district_id" xorm:"not null default 0 comment('买家区县id') INT(11)"`
	CommunityId int    `json:"community_id" xorm:"not null default 0 comment('买家社区id') INT(11)"`
	Address     string `json:"address" xorm:"not null default '' comment('买家地址') VARCHAR(255)"`
	FullAddress string `json:"full_address" xorm:"not null default '' comment('买家详细地址') VARCHAR(255)"`
	Longitude   string `json:"longitude" xorm:"not null default '' comment('买家地址经度') VARCHAR(50)"`
	Latitude    string `json:"latitude" xorm:"not null default '' comment('买家地址纬度') VARCHAR(50)"`
	BuyerIp     string `json:"buyer_ip" xorm:"not null default '' comment('买家ip') VARCHAR(20)"`
	BuyerRemark string `json:"buyer_remark" xorm:"not null default '' comment('买家留言信息') VARCHAR(50)"`
	CreateAt    int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt    int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
