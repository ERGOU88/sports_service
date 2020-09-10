package muser

import (
	"errors"
	"fmt"
	"sports_service/server/dao"
	"sports_service/server/global/consts"
	"sports_service/server/global/login/log"
	"sports_service/server/global/rdskey"
	"sports_service/server/util"
	"strconv"
	"time"
)

type mobileRegister struct {
}

func NewPhoneRegister() *mobileRegister {
	r := &mobileRegister{}
	return r
}

// 用户手机注册
func (r *mobileRegister) Register(u *UserModel, param *LoginParams) error {
	mobileNum, err := strconv.ParseInt(param.MobileNum, 10, 64)
	if err != nil {
		log.Log.Errorf("reg_trace: strconv.ParseInt(%s)", err.Error())
		return err
	}

	key := rdskey.MakeKey(rdskey.LOGIN_REPEAT, consts.TYPE_PHONE, param.MobileNum)
	// 验证重复注册
	ok, err := IsReapeat(key)
	if !ok || err != nil{
		log.Log.Errorf("reg_trace: repeat register mobile num %s", param.MobileNum)
		return errors.New("repeat register")
	}

	rds := dao.NewRedisDao()
	rds.EXPIRE64(key, rdskey.KEY_EXPIRE_MIN)

	r.newUser(u, mobileNum, param)
	if err := u.AddUser(); err != nil {
		log.Log.Errorf("reg_trace: user register err:%s", err)
		return err
	}

	return nil
}

// 设置用户相关信息
func (r *mobileRegister) newUser(u *UserModel, mobileNum int64, param *LoginParams) *UserModel {
	now := time.Now().Unix()
	// todo 暂时先使用时间 + 4位随机数 生成uid
	u.SetUid(util.NewUserId())
	u.SetNickName(r.newDefaultNickName(mobileNum))
	u.SetPhone(mobileNum)
	u.SetAvatar(r.defaultAvatar())
	u.SetDeviceType(param.Platform)
	u.SetCreateAt(now)
	u.SetUpdateAt(now)
	u.SetLoginTime(now)
	u.SetPassword("")
	return u
}

// 默认头像
func (r *mobileRegister) defaultAvatar() string {
	return ""
}

// 默认昵称
func (r *mobileRegister) newDefaultNickName(mobileNum int64) string {
	str := fmt.Sprint(mobileNum)
	length := len(str)
	// 手机号后4位
	str = str[len(str)-4: length]
	return  fmt.Sprintf("用户%s", str)
}

// 验证重复
func IsReapeat(key string) (bool, error) {
	rds := dao.NewRedisDao()
	return rds.SETNX(key)
}
