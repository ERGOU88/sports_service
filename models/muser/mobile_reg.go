package muser

import (
	"errors"
	"fmt"
	"sports_service/server/dao"
	"sports_service/server/global/consts"
	"sports_service/server/global/app/log"
	"sports_service/server/global/rdskey"
	"sports_service/server/util"
	"strconv"
	"time"
)

type mobileRegister struct {
	*base
}

// 实例
func NewMobileRegister() *mobileRegister {
	r := &mobileRegister{&base{}}
	return r
}

// 用户手机注册
func (r *mobileRegister) Register(u *UserModel, platform int, mobileNum, clientIp string) error {
	phone, err := strconv.ParseInt(mobileNum, 10, 64)
	if err != nil {
		log.Log.Errorf("reg_trace: strconv.ParseInt(%s)", err.Error())
		return err
	}

	key := rdskey.MakeKey(rdskey.LOGIN_REPEAT, consts.TYPE_MOBILE, mobileNum)
	// 验证重复注册
	ok, err := IsReapeat(key)
	if !ok || err != nil{
		log.Log.Errorf("reg_trace: repeat register mobile num %s", mobileNum)
		return errors.New("repeat register")
	}

	rds := dao.NewRedisDao()
	if _, err := rds.EXPIRE64(key, rdskey.KEY_EXPIRE_MIN); err != nil {
		log.Log.Errorf("reg_trace: expire login repeat key err:%s", err)
	}

	r.newUser(u, phone, platform, clientIp)
	// 腾讯云im导入用户
	//ch := r.tcyAddUser(u.User)
	//result, ok := <- ch
	//if !ok || result.err != nil {
	//	log.Log.Error("reg_trace: read  chan = %v, err = %+v", ok, result.err)
	//	return errors.New("register im user fail")
	//}
	//
	//u.SetTencentImInfo(u.User.UserId, result.sig)

	return nil
}

// 设置用户相关信息
func (r *mobileRegister) newUser(u *UserModel, phone int64, platform int, clientIp string) *UserModel {
	now := time.Now().Unix()
	// 8位uid
	u.SetUid(r.getUserID())
	u.SetNickName(r.newDefaultNickName(phone))
	u.SetPhone(phone)
	u.SetAvatar(r.defaultAvatar())
	u.SetDeviceType(platform)
	u.SetCreateAt(now)
	u.SetUpdateAt(now)
	u.SetLastLoginTime(now)
	u.SetPassword("")
	u.SetUserType(consts.TYPE_MOBILE)
	u.SetLastLoginTime(now)
	u.SetRegisterIp(clientIp)
	return u
}

// 默认昵称
func (r *mobileRegister) newDefaultNickName(mobileNum int64) string {
	str := fmt.Sprint(mobileNum)
	length := len(str)
	// 手机号后4位
	str = str[len(str)-4: length]
	return  fmt.Sprintf("用户%d%s", util.GenerateRandnum(10000, 99999), str)
}

// 验证重复
func IsReapeat(key string) (bool, error) {
	rds := dao.NewRedisDao()
	return rds.SETNX(key)
}
