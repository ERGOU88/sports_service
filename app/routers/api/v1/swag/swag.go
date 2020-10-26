package swag

import (
  "sports_service/server/models"
  "sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
)

// swagger api文档（登陆接口返回数据）
type LoginSwag struct {
	Token string        `json:"token"`      // token
	User  *User         `json:"user_info"`  // 用户信息
}

// swag文档展示
type ZoneInfoSwag struct {
	UserInfoResp     *muser.UserInfoResp       `json:"user_info"`
	UserZoneInfoResp *UserZoneInfoResp         `json:"zone_info"`
}

// swag文档展示（综合搜索结果）
type ColligateSearchSwag struct {
	UserList      []*muser.UserSearchResults `json:"user_list"`
	VideoList     []*mvideo.VideoDetailInfo  `json:"video_list"`
}

// swag文档展示（客户端初始化返回数据）
type ClientInitSwag struct {
  Secret        string                      `json:"secret" example:"密钥"`
  AvatarList    []models.DefaultAvatar      `json:"avatar_list"`
  WorldList     []models.WorldMap           `json:"world_list"`
  LabelList     []models.VideoLabelConfig   `json:"label_list"`
  LoginTreaty   string                      `json:"login_treaty" example:"登陆协议h5页"`
  UploadTreaty  string                      `json:"upload_treaty" example:"上传协议h5页"`
  FaqH5         string                      `json:"faq_h5" example:"常见问题h5页"`
  About         string                      `json:"about" example:"关于h5页"`
}

// 个人空间用户信息
type UserZoneInfoResp struct {
	TotalBeLiked     int64  `json:"totalBeLiked" example:"100"`     // 被点赞数
	TotalFans        int64  `json:"totalFans" example:"100"`        // 粉丝数
	TotalAttention   int64  `json:"totalAttention" example:"100"`   // 关注数
	TotalCollect     int64  `json:"totalCollect" example:"100"`     // 收藏的作品数
	TotalPublish     int64  `json:"totalPublish" example:"100"`     // 发布的作品数
	TotalLikes       int64  `json:"totalLikes" example:"100"`       // 点赞的作品数
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
