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

// 登陆请求所需的参数
type LoginParams struct {
	Platform  int      `json:"platform" example:"0"`                                // 平台 0 android 1 iOS 2 web
	Token     string   `json:"token" example:"客户端token"`
	OpToken   string   `json:"opToken" example:"客户端返回的运营商token"`
	Operator  string   `json:"operator" example:"客户端返回的运营商，CMCC:中国移动通信, CUCC:中国联通通讯, CTCC:中国电信"`
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

// 添加用户
func (m *UserModel) AddUser() error {
	if _, err := m.Engine.InsertOne(m.User); err != nil {
		log.Log.Errorf("user_trace: add user err:%s", err)
		return err
	}

	return nil
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

