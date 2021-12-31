package muser

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
)

// 社交平台账户
type SocialModel struct {
	SocialAccount   *models.SocialAccountLogin
	Engine          *xorm.Session
}

// 实例
func NewSocialPlatform(engine *xorm.Session) *SocialModel {
	return &SocialModel{
		SocialAccount: &models.SocialAccountLogin{},
		Engine: engine,
	}
}

// 设置用户id
func (m *SocialModel) SetUserId(userId string) {
	m.SocialAccount.UserId = userId
}

// 设置用户openId
func (m *SocialModel) SetOpenId(openId string) {
	m.SocialAccount.OpenId = openId
}

// 设置社交平台关联id
func (m *SocialModel) SetUnionId(unionId string) {
	m.SocialAccount.Unionid = unionId
}

// 设置社交平台类型 1 微信 2 QQ 3 微博
func (m *SocialModel) SetSocialType(socialType int) {
	m.SocialAccount.SocialType = socialType
}

// 设置状态 0 正常 1 封禁 默认正常
func (m *SocialModel) SetStatus(status int) {
	m.SocialAccount.Status = status
}

// 设置创建时间
func (m *SocialModel) SetCreateAt(tm int64) {
	m.SocialAccount.CreateAt = int(tm)
}

// 添加用户社交帐号信息
func (m *SocialModel) AddSocialAccountInfo() error {
	if _, err := m.Engine.InsertOne(m.SocialAccount); err != nil {
		log.Log.Errorf("user_trace: add social account err:%s", err)
		return err
	}

	return nil
}

// 获取社交帐号
func (m *SocialModel) GetSocialAccountByType(socialType int, unionid string) *models.SocialAccountLogin {
	ok, err := m.Engine.Where("social_type=? AND unionid=?", socialType, unionid).Get(m.SocialAccount)
	if !ok || err != nil {
		return nil
	}

	return m.SocialAccount
}



