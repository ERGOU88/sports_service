package muser

import (
	"github.com/go-xorm/xorm"
	"regexp"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"sports_service/server/util"
	"strings"
	"fmt"
	"sports_service/server/global/rdskey"
)

type UserModel struct {
	User    *models.User
	Engine  *xorm.Session
}

func NewUserModel(engine *xorm.Session) *UserModel {
	return &UserModel{
		User: new(models.User),
		Engine: engine,
	}
}

// 用户简单信息返回
type UserInfoResp struct {
	UserId        string `json:"userId" example:"2009011314521111"`
	Avatar        string `json:"avatar" example:"头像地址"`
	MobileNum     int32  `json:"mobileNum" example:"13177656222"`
	NickName      string `json:"nickname" example:"昵称 陈二狗"`
	Gender        int32  `json:"gender" example:"0"`
	Signature     string `json:"signature" example:"个性签名"`
	Status        int32  `json:"status" example:"0"`
	IsAnchor      int32  `json:"isAnchor" example:"0"`
	BackgroundImg string `json:"backgroundImg" example:"背景图"`
	Born          string `json:"born" example:"出生日期"`
	Age           int    `json:"age" example:"27"`
	UserType      int    `json:"userType" example:"0"`
	Country       int32  `json:"country" example:"0"`
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

// 登陆请求所需的参数
type LoginParams struct {
	Platform  int      `json:"platform" example:"0"`                                // 平台 0 android 1 iOS 2 web
	Token     string   `json:"token" example:"客户端token"`
	OpToken   string   `json:"opToken" example:"客户端返回的运营商token"`
	Operator  string   `json:"operator" example:"客户端返回的运营商，CMCC:中国移动通信, CUCC:中国联通通讯, CTCC:中国电信"`
}

// 修改用户信息请求参数
type EditUserInfoParams struct {
	AvatarId   int32    `json:"avatarId" example:"6"`            // 系统头像id（暂时仅支持更换系统默认头像）
	NickName   string   `json:"nickName" example:"陈二狗"`        // 昵称
	Born       string   `json:"born" example:"1993-06-20"`       // 出生年月
	Gender     int32    `json:"gender" example:"1"`              // 性别 1 男 2 女
	CountryId  int32    `json:"countryId" example:"1"`           // 国家id
	Signature  string   `json:"signature" example:"emmmmmmmm"`   // 个性签名
}

// 用户反馈请求参数
type FeedbackParam struct {
	Phone    string `json:"phone" example:"手机号"`
	Describe string `binding:"required" json:"describe" example:"问题描述"`    // 反馈内容
	Problem  string `json:"problem" example:"遇到的问题"`                       // 遇到的问题
	Pics     string `json:"pics" example:"图片列表"`                           // 图片（多张逗号分隔）
}

// 个人空间 用户信息请求参数
type UserZoneInfoParam struct {
	UserId  string `json:"userId"`            // 用户userid
}

var validPhone = regexp.MustCompile(`^1\d{10}$`)
// 检验手机号
func (m *UserModel) CheckCellPhoneNumber(mobileNum string) bool {
	return validPhone.MatchString(mobileNum)
}

// 手机号查询用户
func (m *UserModel) FindUserByPhone(mobileNum string) *models.User {
	ok, err := m.Engine.Where("mobile_num=?", mobileNum).Get(m.User)
	if !ok || err != nil {
		log.Log.Errorf("user_trace: find user by phone err:%s", err)
		return nil
	}

	return m.User
}

// userid查询用户
func (m *UserModel) FindUserByUserid(userId string) *models.User {
	ok, err := m.Engine.Where("user_id=?", userId).Get(m.User)
	if !ok || err != nil {
		log.Log.Errorf("user_trace: find user by userid err:%s", err)
		return nil
	}

	return m.User
}

// userid列表查询用户
func (m *UserModel) FindUserByUserids(userIds string, offset, size int) []*models.User {
	var list []*models.User
	sql := fmt.Sprintf("SELECT * FROM user WHERE user_id in(%s) ORDER BY id DESC LIMIT ?, ?", userIds)
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		log.Log.Errorf("user_trace: get user list err:%s", err)
		return nil
	}

	return list
}

// 添加用户
func (m *UserModel) AddUser() error {
	if _, err := m.Engine.InsertOne(m.User); err != nil {
		log.Log.Errorf("user_trace: add user err:%s", err)
		return err
	}

	return nil
}

// 昵称是否重复
func (m *UserModel) IsRepeatOfNickName(nickName string) bool {
	count, _ := m.Engine.Where("nick_name = ?", nickName).Count(&models.User{})
	if count > 0 {
		return true
	}

	return false
}

// 更新用户信息
func (m *UserModel) UpdateUserInfo() error {
	if _, err := m.Engine.Where("id=?", m.User.Id).
		Cols("avatar, nick_name, born, age, gender, country, signature, update_at").
		Update(m.User); err != nil {
			return err
	}

	return nil
}

// 获取世界信息（暂时只有国家）
func (m *UserModel) GetWorldInfo() []*models.WorldMap {
	var list []*models.WorldMap
	if err := m.Engine.Where("status=0").Desc("sortorder").Find(&list); err != nil {
		return nil
	}

	return list
}


// 通过id获取世界信息（暂时只有国家）
func (m *UserModel) GetWorldInfoById(id int32) *models.WorldMap {
	info := new(models.WorldMap)
	ok, err := m.Engine.Where("id=?", id).Get(info)
	if !ok || err != nil {
		return nil
	}

	return info
}

// 获取系统默认头像列表
func (m *UserModel) GetSystemAvatarList() []*models.DefaultAvatar {
	var list []*models.DefaultAvatar
	if err := m.Engine.Desc("sortorder").Find(&list); err != nil {
		return nil
	}

	return list
}

// 记录用户反馈信息
func (m *UserModel) RecordUserFeedback(userId string, param *FeedbackParam, now int) error {
	feedback := new(models.Feedback)
	feedback.Describe = param.Describe
	feedback.Problem = param.Problem
	feedback.CreateAt = now
	feedback.UpdateAt = now
	feedback.UserId = userId
	feedback.Phone = param.Phone
	feedback.Pics = param.Pics

	if _, err := m.Engine.InsertOne(feedback); err != nil {
		return err
	}

	return nil
}

// 通过id获取系统默认头像
func (m *UserModel) GetSystemAvatarById(id int32) *models.DefaultAvatar {
	info := new(models.DefaultAvatar)
	ok, err := m.Engine.Where("id=?", id).Get(info)
	if !ok || err != nil {
		return nil
	}

	return info
}

// 设置uid
func (m *UserModel) SetUid(uid string) {
	m.User.UserId = uid
}

// 设置昵称
func (m *UserModel) SetNickName(name string) {
	m.User.NickName = name
}

// 设置手机号
func (m *UserModel) SetPhone(mobileNum int64) {
	m.User.MobileNum = mobileNum
}

// 设置用户头像
func (m *UserModel) SetAvatar(avatar string) {
	m.User.Avatar = avatar
}

// 设置用户类型
func (m *UserModel) SetUserDeviceType(utype int) {
	m.User.DeviceType = utype
}

// 设备类型
func (m *UserModel) SetDeviceType(deviceType int) {
	m.User.DeviceType = deviceType
}

// 设置性别
func (m *UserModel) SetGender(gender int) {
	m.User.Gender = gender
}

// 设置用户状态
func (m *UserModel) SetStatus(status int) {
	m.User.Status = status
}

// 设置创建时间
func (m *UserModel) SetCreateAt(tm int64) {
	m.User.CreateAt= int(tm)
}

// 设置更新时间
func (m *UserModel) SetUpdateAt(tm int64) {
	m.User.UpdateAt = int(tm)
}

// 设置密码
func (m *UserModel) SetPassword(password string) {
	m.User.Password = password
}

// 设置用户类型 1 手机号 2 微信 3 QQ 4 微博
func (u *UserModel) SetUserType(userType int) {
	u.User.UserType = userType
}

// 设置登陆时间
func (u *UserModel) SetLastLoginTime(tm int64) {
	u.User.LastLoginTime = int(tm)
}

// 获取用户token
func (m *UserModel) GetUserToken(uid string) (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.MakeKey(rdskey.USER_AUTH, uid))
}

// 保存token(1周过期)
func (m *UserModel) SaveUserToken(uid, token string) error {
	rds := dao.NewRedisDao()
	return rds.SETEX(rdskey.MakeKey(rdskey.USER_AUTH, uid), rdskey.KEY_EXPIRE_WEEK, token)
}

// 获取token
func (m *UserModel) GenUserToken(uid, password string) string {
	auth := uid + "_" + strings.ToLower(util.Md5String(uid+"|"+password))
	return auth
}

// 最终按加密规则token校验
func (m *UserModel) TokenValid(account, password, hashcode string) (b bool) {
	md := strings.ToLower(util.MD5(fmt.Sprintf("%s|%s", account, util.Md5String(password))))
	log.Log.Debugf("@@@@@hashcode:%s, md:%s", hashcode, md)
	return md == hashcode
}

