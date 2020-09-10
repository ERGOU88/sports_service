package swag

// 登陆请求所需的参数
type LoginParamsSwag struct {
	MobileNum string `binding:"required" json:"mobileNum" example:"手机号码"`    // 手机号
	Platform  int    `json:"platform" example:"1"`                              // 平台
}

// swagger api文档（登陆接口返回数据）
type LoginSwag struct {
	Token string        `json:"token"` // token
	User  *User         `json:"user"`  // 用户信息
}

type User struct {
	Id            int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	NickName      string `json:"nick_name" xorm:"not null default '' comment('昵称') VARCHAR(45)"`
	MobileNum     int64  `json:"mobile_num" xorm:"not null comment('手机号码') BIGINT(20)"`
	Password      string `json:"password" xorm:"not null comment('用户密码') VARCHAR(128)"`
	UserId        string `json:"user_id" xorm:"not null comment('用户id') unique VARCHAR(60)"`
	Gender        int    `json:"gender" xorm:"not null default 0 comment('性别 0男性 1女性') TINYINT(1)"`
	Born          string `json:"born" xorm:"not null default '' comment('出生日期') VARCHAR(128)"`
	Age           int    `json:"age" xorm:"not null default 0 comment('年龄') INT(3)"`
	Avatar        string `json:"avatar" xorm:"not null default '' comment('头像') VARCHAR(100)"`
	Status        int    `json:"status" xorm:"default 0 comment('0 正常 1 封禁') TINYINT(1)"`
	LastLoginTime int    `json:"last_login_time" xorm:"comment('最后登录时间') INT(11)"`
	Signature     string `json:"signature" xorm:"not null default '' comment('签名') VARCHAR(200)"`
	DeviceType    int    `json:"device_type" xorm:"comment('设备类型 0 android 1 iOS 2 web') TINYINT(2)"`
	City          string `json:"city" xorm:"not null default '' comment('城市') VARCHAR(64)"`
	IsAnchor      int    `json:"is_anchor" xorm:"not null default 0 comment('0不是主播 1为主播') TINYINT(1)"`
	ChannelId     int    `json:"channel_id" xorm:"not null default 0 comment('渠道id') INT(11)"`
	BackgroundImg string `json:"background_img" xorm:"not null default '' comment('背景图') VARCHAR(255)"`
	Title         string `json:"title" xorm:"not null default '' comment('称号/特殊身份') VARCHAR(255)"`
	CreateAt      int    `json:"create_at" xorm:"not null default 0 comment('创建时间') INT(11)"`
	UpdateAt      int    `json:"update_at" xorm:"not null default 0 comment('更新时间') INT(11)"`
}
