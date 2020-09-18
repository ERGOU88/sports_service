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
	VideoId       int64     `binding:"required" json:"videoId" example:"10001"`          // 点赞的视频id
	ToUserId      string    `binding:"required" json:"toUserId" example:"被点赞的用户"`    // 被点赞的用户id
}

// 取消点赞请求参数
type CancelLikeParam struct {
	VideoId       int64     `binding:"required" json:"videoId" example:"10001"`          // 取消点赞的视频id
}

// 实栗
func NewLikeModel(engine *xorm.Session) *LikeModel {
	return &LikeModel{
		Like: new(models.ThumbsUp),
		Engine: engine,
	}
}

type LikeVideosInfo struct {
	TypeId   int64    `json:"type_id"`      // 视频id
	CreateAt int      `json:"create_at"`    // 点赞时间
}
// 获取用户点赞的视频id列表
func (m *LikeModel) GetUserLikeVideos(userId string, offset, size int) []*LikeVideosInfo {
	var list []*LikeVideosInfo
	if err := m.Engine.Table(&models.ThumbsUp{}).Where("zan_type=1 AND status=1 AND user_id=?", userId).
		Cols("type_id", "create_at").
		Desc("id").
		Limit(size, offset).
		Find(&list); err != nil {
		log.Log.Errorf("like_trace: get like videos err:%s", err)
		return nil
	}

	return list
}

// 添加点赞记录 1 视频点赞 2 帖子点赞 3 评论点赞
func (m *LikeModel) AddGiveLikeByType(userId, toUserId string, composeId int64, status, zanType int) error {
	m.Like.UserId = userId
	m.Like.ToUserId = toUserId
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
	m.Like = new(models.ThumbsUp)
	ok, err := m.Engine.Where("user_id=? AND type_id=? AND zan_type=?", userId, composeId, zanType).Get(m.Like)
	if !ok || err != nil {
		return nil
	}

	return m.Like
}

// 更新点赞状态 点赞/取消点赞
func (m *LikeModel) UpdateLikeStatus() error {
	if _, err := m.Engine.ID(m.Like.Id).
		Cols("status, create_at").
		Update(m.Like); err != nil {
		return err
	}

	return nil
}

// 获取用户被点赞总数
func (m *LikeModel) GetUserTotalBeLiked(userId string) int64 {
	total, err := m.Engine.Where("to_user_id=? AND status=1", userId).Count(m.Like)
	if err != nil {
		log.Log.Errorf("like_trace: get user total be liked err:%s, uid:%s", err, userId)
		return 0
	}

	return total
}

// 获取用户点赞总数
func (m *LikeModel) GetUserTotalLikes(userId string) int64 {
	total, err := m.Engine.Where("user_id=? AND status=1", userId).Count(m.Like)
	if err != nil {
		log.Log.Errorf("like_trace: get user total likes err:%s, uid:%s", err, userId)
		return 0
	}

	return total
}

