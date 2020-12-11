package models

type User struct {
	Age           int    `json:"age" xorm:"not null default 0 comment('年龄') INT(3)"`
	Avatar        string `json:"avatar" xorm:"not null default '' comment('头像') VARCHAR(300)"`
	BackgroundImg string `json:"background_img" xorm:"not null default '' comment('背景图') VARCHAR(255)"`
	Born          string `json:"born" xorm:"not null default '' comment('出生日期') VARCHAR(128)"`
	ChannelId     int    `json:"channel_id" xorm:"not null default 0 comment('渠道id') INT(11)"`
	City          string `json:"city" xorm:"not null default '' comment('城市') VARCHAR(64)"`
	Country       int    `json:"country" xorm:"not null default 0 comment('国家') INT(3)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	DeviceToken   string `json:"device_token" xorm:"not null default '' comment('设备token') VARCHAR(100)"`
	DeviceType    int    `json:"device_type" xorm:"comment('设备类型 0 android 1 iOS 2 web') TINYINT(2)"`
	Gender        int    `json:"gender" xorm:"not null default 0 comment('性别 0人妖 1男性 2女性') TINYINT(1)"`
	Id            int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	IsAnchor      int    `json:"is_anchor" xorm:"not null default 0 comment('0不是主播 1为主播') TINYINT(1)"`
	LastLoginTime int    `json:"last_login_time" xorm:"comment('最后登录时间') INT(11)"`
	MobileNum     int64  `json:"mobile_num" xorm:"not null comment('手机号码') BIGINT(20)"`
	NickName      string `json:"nick_name" xorm:"not null default '' comment('昵称') VARCHAR(45)"`
	Password      string `json:"password" xorm:"not null comment('用户密码') VARCHAR(128)"`
	RegIp         string `json:"reg_ip" xorm:"default ' ' comment('注册ip') VARCHAR(30)"`
	Signature     string `json:"signature" xorm:"not null default '' comment('签名') VARCHAR(200)"`
	Status        int    `json:"status" xorm:"default 0 comment('0 正常 1 封禁') TINYINT(1)"`
	Title         string `json:"title" xorm:"not null default '' comment('称号/特殊身份') VARCHAR(255)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') unique VARCHAR(60)"`
	UserType      int    `json:"user_type" xorm:"not null default 0 comment('用户类型 0 手机号 1 微信 2 QQ 3 微博') TINYINT(2)"`
}
