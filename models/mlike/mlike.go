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
	ComposeId int64  `binding:"required" json:"compose_id" example:"10000000000"`      // 点赞的视频/评论/帖子 id
}

// 取消点赞请求参数
type CancelLikeParam struct {
	ComposeId int64 `binding:"required" json:"compose_id" example:"10000000000"`       // 取消点赞的视频/评论/帖子 id
}

// 被点赞的视频信息
type BeLikedVideoInfo struct {
	ComposeId     int64                 `json:"compose_id"`      // 作品id
	Title         string                `json:"title"`           // 标题
	Describe      string                `json:"describe"`        // 描述
	Cover         string                `json:"cover"`           // 封面
	VideoAddr     string                `json:"video_addr"`      // 视频地址
	VideoDuration int                   `json:"video_duration"`  // 视频时长
	VideoWidth    int64                 `json:"video_width"`     // 视频宽
	VideoHeight   int64                 `json:"video_height"`    // 视频高
	CreateAt      int                   `json:"create_at"`       // 视频创建时间
	BarrageNum    int                   `json:"barrage_num"`     // 弹幕数
	BrowseNum     int                   `json:"browse_num"`      // 浏览数（播放数）
	//ToUserId      string                `json:"to_user_id"`      // 被点赞的用户id
	//ToUserName    string                `json:"to_user_name"`    // 被点赞的用户昵称
  NickNames    []string `json:"nick_names"`     // 点赞用户昵称（多个）
  TotalLikeNum int      `json:"total_like_num"` // 总点赞数
  Avatar       string   `json:"avatar"`         // 最近点赞的用户头像
	OpTime       int      `json:"op_time"`        // 用户点赞操作时间
	Type         int      `json:"type"`           // 类型 1 视频 2 帖子 3 评论
}

// todo: 后续使用
// 点赞的用户信息
type ZanUserInfo struct {
  UserId        string                `json:"user_id"`         // 点赞的用户id
  Avatar        string                `json:"avatar"`          // 点赞用户头像
  Nickname      string                `json:"nick_name"`       // 点赞的用户昵称
}

// 点赞用户信息
type LikedUserInfo struct {
  UserId      string   `json:"user_id"`
  NickName    string   `json:"nick_name"`
  Avatar      string   `json:"avatar"`
  OpTm        int      `json:"op_tm"`       // 用户点赞时间
}

// 被点赞的信息
type BeLikedInfo struct {
  ComposeId     int64                 `json:"compose_id"`      // 作品id
  Title         string                `json:"title"`           // 标题
  Describe      string                `json:"describe"`        // 描述
  Cover         string                `json:"cover"`           // 封面
  VideoAddr     string                `json:"video_addr"`      // 视频地址
  VideoDuration int                   `json:"video_duration"`  // 视频时长
  VideoWidth    int64                 `json:"video_width"`     // 视频宽
  VideoHeight   int64                 `json:"video_height"`    // 视频高
  CreateAt      int                   `json:"create_at"`       // 视频创建时间
  BarrageNum    int                   `json:"barrage_num"`     // 弹幕数
  BrowseNum     int                   `json:"browse_num"`      // 浏览数（播放数）
  //UserId        string                `json:"user_id"`         // 点赞的用户id
  //Avatars       []string              `json:"avatars"`          // 点赞用户头像
  //Nicknames     []string              `json:"nick_name"`       // 点赞的用户昵称
  TotalLikeNum  int                   `json:"total_like_num"`  // 总点赞数
  UserList      []*LikedUserInfo      `json:"user_list"`       // 点赞的用户列表
  //ToUserId      string                `json:"to_user_id"`      // 被点赞的用户id
  //ToUserAvatar  string                `json:"avatar"`          // 被点赞用户头像
  //ToUserName    string                `json:"to_user_name"`    // 被点赞的用户昵称
  Content       string                `json:"content"`         // 被点赞的评论内容
  OpTime        int                   `json:"op_time"`         // 用户点赞操作时间
  Type          int                   `json:"type"`            // 类型 1 视频 2 帖子 3 评论

  ParentCommentId  int64              `json:"parent_comment_id,omitempty"`  // 父级评论id
}

// 被点赞的评论信息
type BeLikedCommentInfo struct {
	ComposeId     int64                 `json:"compose_id"`      // 作品id
	Title         string                `json:"title"`           // 标题
	Describe      string                `json:"describe"`        // 描述
	Cover         string                `json:"cover"`           // 封面
	VideoAddr     string                `json:"video_addr"`      // 视频地址
	VideoDuration int                   `json:"video_duration"`  // 视频时长
	VideoWidth    int64                 `json:"video_width"`     // 视频宽
	VideoHeight   int64                 `json:"video_height"`    // 视频高
	CreateAt      int                   `json:"create_at"`       // 视频创建时间
	BarrageNum    int                   `json:"barrage_num"`     // 弹幕数
	BrowseNum     int                   `json:"browse_num"`      // 浏览数（播放数）
	//UserId        string                `json:"user_id"`         // 点赞的用户id
	Avatar        string                `json:"avatar"`          // 点赞用户头像
	Nicknames     []string              `json:"nick_name"`       // 点赞的用户昵称
	TotalLikeNum  int                   `json:"total_like_num"`  // 总点赞数
	//ToUserId      string                `json:"to_user_id"`      // 被点赞的用户id
	//ToUserAvatar  string                `json:"avatar"`          // 被点赞用户头像
	//ToUserName    string                `json:"to_user_name"`    // 被点赞的用户昵称
	Content       string                `json:"content"`         // 被点赞的评论内容
	OpTime        int                   `json:"op_time"`         // 用户点赞操作时间
	Type          int                   `json:"type"`            // 类型 1 视频 2 帖子 3 评论
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

// 获取用户被点赞的记录 包含 视频、评论等
func (m *LikeModel) GetBeLikedList(toUserId string, offset, size int) []*models.ThumbsUp {
	var list []*models.ThumbsUp
	if err := m.Engine.Where("to_user_id=? AND status=1", toUserId).Desc("create_at", "id").Limit(size, offset).Find(&list); err != nil {
		log.Log.Errorf("like_trace: get be liked list err:%s", err)
		return nil
	}

	return list
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
	total, err := m.Engine.Where("to_user_id=? AND status=1", userId).Count(&models.ThumbsUp{})
	if err != nil {
		log.Log.Errorf("like_trace: get user total be liked err:%s, uid:%s", err, userId)
		return 0
	}

	return total
}

// 获取用户点赞总数
func (m *LikeModel) GetUserTotalLikes(userId string) int64 {
	total, err := m.Engine.Where("user_id=? AND status=1", userId).Count(&models.ThumbsUp{})
	if err != nil {
		log.Log.Errorf("like_trace: get user total likes err:%s, uid:%s", err, userId)
		return 0
	}

	return total
}

// 根据点赞类型、作品id 获取点赞数量
func (m *LikeModel) GetLikeNumByType(composeId int64, zanType int) int64 {
	count, err := m.Engine.Where("type_id=? AND zan_type=? AND status=1", composeId, zanType).Count(&models.ThumbsUp{})
	if err != nil {
		log.Log.Errorf("like_trace: get like num by type err:%s", err)
		return 0
	}

	return count
}

// 获取未读的被点赞的数量
func (m *LikeModel) GetUnreadBeLikedCount(userId, readTm string) int64 {
	count, err := m.Engine.Where("to_user_id=? AND status=1 AND create_at > ?", userId, readTm).Count(&models.ThumbsUp{})
	if err != nil {
		log.Log.Errorf("like_trace: get unread be liked count err:%s", err)
		return 0
	}

	return count
}
