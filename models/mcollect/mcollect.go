package mcollect

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
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

// 获取收藏的视频id列表
func (m *CollectModel) GetCollectVideos(userId string) []string {
	var list []string
	if err := m.Engine.Where("status=1 AND user_id=?", userId).Cols("video_id").Find(&list); err != nil {
		log.Log.Errorf("collect_trace: get collect videos err:%s", err)
		return nil
	}

	return list
}
