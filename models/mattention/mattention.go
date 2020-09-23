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
	UserId       string     `binding:"required" json:"user_id" example:"需关注的用户id"`     // 需关注的用户id
}

// 取消关注请求参数
type CancelAttentionParam struct {
	UserId       string     `binding:"required" json:"user_id" example:"被取消关注的用户id"`  // 被取消关注的用户id
}

func NewAttentionModel(engine *xorm.Session) *AttentionModel {
	return &AttentionModel{
		UserAttention: new(models.UserAttention),
		Engine: engine,
	}
}

// 获取关注的信息
func (m *AttentionModel) GetAttentionInfo(attentionUid, userId string) *models.UserAttention {
	m.UserAttention = new(models.UserAttention)
	ok, err := m.Engine.Where("attention_uid=? AND user_id=?", attentionUid, userId).Get(m.UserAttention)
	if !ok || err != nil {
		return nil
	}

	return m.UserAttention
}

// 添加关注记录 (attentionUid 关注的用户id userId 被关注的用户id)
func (m *AttentionModel) AddAttention(attentionUid, userId string, status int) error {
	m.UserAttention.UserId = userId
	m.UserAttention.AttentionUid = attentionUid
	m.UserAttention.CreateAt = int(time.Now().Unix())
	m.UserAttention.Status = status
	if _, err := m.Engine.InsertOne(m.UserAttention); err != nil {
		return err
	}

	return nil
}

// 更新关注状态 关注/取消关注
func (m *AttentionModel) UpdateAttentionStatus() error {
	if _, err := m.Engine.ID(m.UserAttention.Id).
		Cols("status, create_at").
		Update(m.UserAttention); err != nil {
		return err
	}

	return nil
}

// 获取用户关注的列表
func (m *AttentionModel) GetAttentionList(attentionUid string) []string {
	var list []string
	if err := m.Engine.Table(&models.UserAttention{}).Where("status=1 AND attention_uid=?", attentionUid).Cols("user_id").Find(&list); err != nil {
		log.Log.Errorf("attention_trace: get attention list err:%s", err)
		return nil
	}

	return list
}

// 获取用户的粉丝列表
func (m *AttentionModel) GetFansList(userId string) []string {
	var list []string
	if err := m.Engine.Table(&models.UserAttention{}).Where("status=1 AND user_id=?", userId).Cols("attention_uid").Find(&list); err != nil {
		log.Log.Errorf("attention_trace: get fans list err:%s", err)
		return nil
	}

	return list
}

// 获取用户的关注总数
func (m *AttentionModel) GetTotalAttention(userId string) int64 {
	total, err := m.Engine.Where("status=1 AND attention_uid=?", userId).Count(m.UserAttention)
	if err != nil {
		log.Log.Errorf("attention_trace: get attention total err:%s, uid:%s", err, userId)
		return 0
	}

	return total
}

// 获取用户的粉丝总数
func (m *AttentionModel) GetTotalFans(userId string) int64 {
	total, err := m.Engine.Where("status=1 AND user_id=?", userId).Count(m.UserAttention)
	if err != nil {
		log.Log.Errorf("attention_trace: get fans total err:%s, uid:%s", err, userId)
		return 0
	}

	return total
}


// 用户搜索联系人（关注/粉丝）
type SearchContactRes struct {
	UserId        string `json:"user_id" example:"2009011314521111"`
	Avatar        string `json:"avatar" example:"头像地址"`
	NickName      string `json:"nick_name" example:"昵称 陈二狗"`
	Gender        int32  `json:"gender" example:"0"`
	Signature     string `json:"signature" example:"个性签名"`
	Status        int32  `json:"status" example:"0"`
	IsAnchor      int32  `json:"is_anchor" example:"0"`
	BackgroundImg string `json:"background_img" example:"背景图"`
	Born          string `json:"born" example:"出生日期"`
	Age           int    `json:"age" example:"27"`
	IsAttention   int32  `json:"is_attention"`
}

const (
 SEARCH_ATTENTION = "SELECT u.*, ua.status as is_attention FROM user_attention AS ua INNER JOIN user as u ON u.`user_id` = ua.`user_id` " +
 	"WHERE u.nick_name LIKE '%?%' OR u.`user_id` LIKE '%?%' AND u.status=0 AND ua.attention_uid = ? AND ua.status=1 ORDER BY ua.Id DESC LIMIT ?, ?"
)
// 搜索关注的用户
func (m *AttentionModel) SearchAttentionUser(userId, name string, offset, size int) []*SearchContactRes {
	var list []*SearchContactRes
	ok, err := m.Engine.Table(&models.UserAttention{}).SQL(SEARCH_ATTENTION, name, name, userId, offset, size).Get(&list)
	if !ok || err != nil {
		return nil
	}

	return list
}

const (
	SEARCH_FANS = "SELECT u.*, ua.status as is_attention FROM user_attention AS ua INNER JOIN user as u ON u.`user_id` = ua.`attention_uid` " +
		"WHERE u.nick_name LIKE '%?%' OR u.`user_id` LIKE '%?%' AND u.status=0 AND ua.user_id = ? AND ua.status=1 ORDER BY ua.Id DESC LIMIT ?, ?"
)
// 搜索粉丝
func (m *AttentionModel) SearchFans(userId, name string, offset, size int) []*SearchContactRes {
	var list []*SearchContactRes
	ok, err := m.Engine.Table(&models.UserAttention{}).SQL(SEARCH_FANS, name, name, userId, offset, size).Get(&list)
	if !ok || err != nil {
		return nil
	}

	return list
}

