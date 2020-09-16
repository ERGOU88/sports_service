package mcollect

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
	"time"
)

type CollectModel struct {
	CollectVideo *models.CollectVideoRecord
	Engine       *xorm.Session
}

// 添加收藏请求参数
type AddCollectParam struct {
	VideoId       int64     `json:"videoId" example:"10001"`     // 收藏的视频id
}

// 取消收藏请求参数
type CancelCollectParam struct {
	VideoId       int64     `json:"videoId" example:"10001"`     // 取消收藏的视频id
}

// 实栗
func NewCollectModel(engine *xorm.Session) *CollectModel {
	return &CollectModel{
		CollectVideo: new(models.CollectVideoRecord),
		Engine:       engine,
	}
}

// 添加视频收藏
func (m *CollectModel) AddCollectVideo(userId string, videoId int64, status int) error {
	m.CollectVideo.UserId = userId
	m.CollectVideo.VideoId = videoId
	m.CollectVideo.UpdateAt = int(time.Now().Unix())
	m.CollectVideo.CreateAt = int(time.Now().Unix())
	m.CollectVideo.Status = status
	if _, err := m.Engine.InsertOne(m.CollectVideo); err != nil {
		return err
	}

	return nil
}

// 取消视频收藏
func (m *CollectModel) CancelCollectVideo() {
	return
}

//func (m *AttentionModel) GetAttentionInfo(attentionUid, userId string) *models.UserAttention {
//	ok, err := m.Engine.Where("attention_uid=? AND user_id=?", attentionUid, userId).Get(m.UserAttention)
//	if !ok || err != nil {
//		return nil
//	}
//
//	return m.UserAttention
//}

// 获取收藏的信息
func (m *CollectModel) GetCollectInfo(userId string, videoId int64) *models.CollectVideoRecord {
	ok, err := m.Engine.Where("user_id=? AND video_id=?", userId, videoId).Get(m.CollectVideo)
	if !ok || err != nil {
		return nil
	}

	return m.CollectVideo
}

// 更新收藏状态 收藏/取消收藏
func (m *CollectModel) UpdateCollectStatus() error {
	if _, err := m.Engine.Where("id=?", m.CollectVideo.Id).
		Cols("status, update_at").
		Update(m.CollectVideo); err != nil {
		return err
	}

	return nil
}

