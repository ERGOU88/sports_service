package sms

import (
	"math/rand"
	"sports_service/server/global/rdskey"
	notify "sports_service/server/tools/goNotify"
	"sports_service/server/global/consts"
	"time"
	"fmt"
	"sports_service/server/global/app/log"
	"sports_service/server/dao"
)

var (
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type SmsModel struct {

}

// 发送短信验证码请求参数
type SendSmsCodeParams struct {
	SendType       string     `json:"send_type"`    // 短信类型 1 账户登陆/注册
	MobileNum      string     `json:"mobile_num"`   // 手机号码
}

// 手机验证码登陆 请求参数
type SmsCodeLoginParams struct {
	MobileNum      string     `binding:"required" json:"mobile_num"`   // 手机号码
	Code           string     `binding:"required" json:"code"`         // 手机验证码
	Platform       int        `json:"platform"`     // 平台 0 android 1 iOS 2 web
}

// 实栗
func NewSmsModel() *SmsModel {
	return new(SmsModel)
}

// 获取验证码
func (m *SmsModel) GetSmsCode() string {
	return fmt.Sprintf("%06d", randSource.Intn(999999))
}

// 获取发送验证码的模版
func (m *SmsModel) GetSendMod(sendType string) string {
	switch sendType {
	// 账户登陆/注册 短信模版
	case consts.ACCOUNT_OPT_TYPE:
		return consts.ACCOUNT_MODE
	default:
		log.Log.Errorf("sms_trace: unsupported sendType, sendType:%s", sendType)
	}

	return ""
}

// 发送短信验证码
func (m *SmsModel) Send(sendMod, mobileNum, code string) error {
	s := &notify.Sms{}
	s.To = mobileNum
	s.ServiceName = consts.SERVICE_NAME
	s.Content = []byte(fmt.Sprintf(sendMod, code))
	if err := s.Send(); err != nil {
		return err
	}

	return nil
}

// 获取24小时内发送短信的限制数量
func (m *SmsModel) GetSendSmsLimitNum(mobileNum string) (int, error) {
	key := rdskey.MakeKey(rdskey.SMS_INTERVAL_NUM, time.Now().Format("2006-01-02"), mobileNum)
	rds := dao.NewRedisDao()
	return rds.GetInt(key)
}

// 增加已发短信的数量（24小时内限制十条）
func (m *SmsModel) IncrSendSmsNum(mobileNum string) error {
	key :=  rdskey.MakeKey(rdskey.SMS_INTERVAL_NUM, time.Now().Format("2006-01-02"), mobileNum)
	rds := dao.NewRedisDao()
	_, err := rds.INCR(key)
	return err
}

// 记录短信验证码次数的key设置过期
func (m *SmsModel) SetSmsIntervalExpire(mobileNum string) (int, error) {
	key :=  rdskey.MakeKey(rdskey.SMS_INTERVAL_NUM, time.Now().Format("2006-01-02"), mobileNum)
	rds := dao.NewRedisDao()
	return rds.EXPIRE(key, rdskey.KEY_EXPIRE_DAY)
}

// 存储验证码 并 设置短信验证码过期时间 5分钟内有效
func (m *SmsModel) SaveSmsCodeByRds(sendType, mobileNum, code string) error {
	key := rdskey.MakeKey(rdskey.SMS_CODE, sendType, mobileNum)
	rds := dao.NewRedisDao()
	return rds.SETEX(key, rdskey.KEY_EXPIRE_MIN * 5, code)
}

// 如果redis能获取到 说明验证码未过期
func (m *SmsModel) GetSmsCodeByRds(sendType, mobileNum string) (string, error) {
	key := rdskey.MakeKey(rdskey.SMS_CODE, sendType, mobileNum)
	rds := dao.NewRedisDao()
	return rds.Get(key)
}

// 删除存储验证码的key
func (m *SmsModel) DelSmsCodeKey(sendType, mobileNum string) error {
	key := rdskey.MakeKey(rdskey.SMS_CODE, sendType, mobileNum)
	rds := dao.NewRedisDao()
	_, err := rds.Del(key)
	return err
}

// 是否已过重发验证码的间隔时间
func (m *SmsModel) HasSmsIntervalPass(sendType, mobileNum string) (bool, error) {
	key := rdskey.MakeKey(rdskey.SMS_INTERVAL_TM, sendType, mobileNum)
	rds := dao.NewRedisDao()
	return rds.EXISTS(key)
}

// 设置重发验证码的间隔时间
func (m *SmsModel) SetSmsIntervalTm(sendType, mobileNum string) error {
	key := rdskey.MakeKey(rdskey.SMS_INTERVAL_TM, sendType, mobileNum)
	rds := dao.NewRedisDao()
	return rds.SETEX(key, rdskey.KEY_EXPIRE_MIN, 1)
}

// 删除限制重发验证码间隔时间的key
func (m *SmsModel) DelSmsIntervalTmKey(sendType, mobileNum string) error {
	key := rdskey.MakeKey(rdskey.SMS_INTERVAL_TM, sendType, mobileNum)
	rds := dao.NewRedisDao()
	_, err := rds.Del(key)
	return err
}

