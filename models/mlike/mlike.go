package mlike

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"time"
)

type LikeModel struct {
	Like         *models.ThumbsUp
	Engine       *xorm.Session
}

// 添加点赞请求参数
type GiveLikeParam struct {
	VideoId       int64     `json:"videoId" example:"10001"`     // 点赞的视频id
}

// 取消点赞请求参数
type CancelLikeParam struct {
	VideoId       int64     `json:"videoId" example:"10001"`     // 取消点赞的视频id
}

// 实栗
func NewLikeModel(engine *xorm.Session) *LikeModel {
	return &LikeModel{
		Like: new(models.ThumbsUp),
		Engine: engine,
	}
}

// 获取用户点赞的视频id列表
func (m *LikeModel) GetUserLikeVideos(userId string) []string {
	var list []string
	if err := m.Engine.Where("zan_type=1 AND status=1 AND user_id=?", userId).Cols("type_id").Find(&list); err != nil {
		log.Log.Errorf("like_trace: get like videos err:%s", err)
		return nil
	}

	return list
}


// 添加点赞记录 1 视频点赞 2 帖子点赞 3 评论点赞
func (m *LikeModel) AddGiveLikeByType(userId string, composeId int64, status, zanType int) error {
	m.Like.UserId = userId
	m.Like.TypeId = composeId
	m.Like.Status = status
	m.Like.ZanType = zanType
	m.Like.CreateAt = int(time.Now().Unix())
	if _, err := m.Engine.InsertOne(m.Like); err != nil {
		return err
	}

	return nil
}

// 获取点赞的信息
func (m *LikeModel) GetLikeInfo(userId string, composeId int64, zanType int) *models.ThumbsUp {
	ok, err := m.Engine.Where("user_id=? AND type_id=? AND zan_type=?", userId, composeId, zanType).Get(m.Like)
	if !ok || err != nil {
		return nil
	}

	return m.Like
}

// 更新点赞状态 点赞/取消点赞
func (m *LikeModel) UpdateLikeStatus() error {
	if _, err := m.Engine.Where("id=?", m.Like.Id).
		Cols("status, create_at").
		Update(m.Like); err != nil {
		return err
	}

	return nil
}

