package mattention

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"time"
)

type AttentionModel struct {
	UserAttention    *models.UserAttention
	Engine           *xorm.Session
}

// 添加关注请求参数
type AddAttentionParam struct {
	UserId       string     `json:"userId" example:"需关注的用户id"`     // 需关注的用户id
}

// 取消关注请求参数
type CancelAttentionParam struct {
	UserId       string     `json:"userId" example:"被取消关注的用户id"`  // 被取消关注的用户id
}

func NewAttentionModel(engine *xorm.Session) *AttentionModel {
	return &AttentionModel{
		UserAttention: new(models.UserAttention),
		Engine: engine,
	}
}

// 获取关注的信息
func (m *AttentionModel) GetAttentionInfo(attentionUid, userId string) *models.UserAttention {
	ok, err := m.Engine.Where("attention_uid=? AND user_id=?", attentionUid, userId).Get(m.UserAttention)
	if !ok || err != nil {
		return nil
	}

	return m.UserAttention
}

// 添加关注记录 (attentionUid 关注的用户id userId 被关注的用户id)
func (m *AttentionModel) AddAttention(attentionUid, userId string) error {
	m.UserAttention.UserId = userId
	m.UserAttention.AttentionUid = attentionUid
	m.UserAttention.CreateAt = int(time.Now().Unix())
	if _, err := m.Engine.InsertOne(m.UserAttention); err != nil {
		return err
	}

	return nil
}

// 更新关注状态 重新关注/取消关注
func (m *AttentionModel) UpdateAttentionStatus() error {
	if _, err := m.Engine.Where("id=?", m.UserAttention.Id).
		Cols("status, create_at").
		Update(m.UserAttention); err != nil {
		return err
	}

	return nil
}

// 获取用户关注的列表
func (m *AttentionModel) GetAttentionList(attentionUid string) []string {
	var list []string
	if err := m.Engine.Where("status=1 AND attention_uid=?", attentionUid).Cols("user_id").Find(&list); err != nil {
		log.Log.Errorf("attention_trace: get attention list err:%s", err)
		return nil
	}

	return list
}

// 获取用户的粉丝列表
func (m *AttentionModel) GetFansList(userId string) []string {
	var list []string
	if err := m.Engine.Where("status=1 AND user_id=?", userId).Cols("attention_uid").Find(&list); err != nil {
		log.Log.Errorf("attention_trace: get fans list err:%s", err)
		return nil
	}

	return list
}

