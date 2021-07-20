package mcomment

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/backend/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"fmt"
)

type CommentModel struct {
	Engine       *xorm.Session
	VideoComment *models.VideoComment
	ReceiveAt    *models.ReceivedAt
	Report       *models.CommentReport
	PostComment  *models.PostingComment
}

// 评论列表 [视频/帖子]
type CommentList struct {
	Id                  int64               `json:"id" example:"1000000000"`                      // 评论id
	IsTop               int                 `json:"is_top"  example:"1"`                          // 置顶状态 1 置顶 0 不置顶
	Avatar              string              `json:"avatar"  example:"用户头像"`                    // 用户头像
	UserId              string              `json:"user_id" example:"用户id"`                     // 用户id
	UserName            string              `json:"user_name" example:"用户昵称"`                  // 用户名称
	CommentLevel        int                 `json:"comment_level" example:"1"`                    // 评论等级[1 一级评论 默认 ，2 二级评论]
	Content             string              `json:"content"  example:"评论的内容"`                 // 评论内容
	CreateAt            int                 `json:"create_at" example:"1600000000"`              // 创建时间
	Status              int                 `json:"status" example:"0"`                          // 状态 1 有效 0 逻辑删除
	VideoId             int64               `json:"video_id,omitempty" example:"1000000000"`     // 视频id
	PostId              int64               `json:"post_id,omitempty"`                           // 帖子ID
	ReplyList           []*ReplyComment     `json:"reply_list"`                                  // 回复列表
	LikeNum             int64               `json:"like_num" example:"100"`                      // 点赞数
	IsAttention         int                 `json:"is_attention" example:"0"`                    // 是否关注
	ReplyNum            int64               `json:"reply_num" example:"100"`                     // 总回复数
	IsLike              int                 `json:"is_like" example:"0"`                         // 是否点赞
	HasMore             int                 `json:"has_more" example:"0"`                        // 是否显示更多 0 不展示  1 展示
}

// 回复评论的内容
type ReplyComment struct {
	Id                   int64               `json:"id" example:"1600000000"`                          // 评论id
	IsTop                int                 `json:"is_top" example:"1"`                               // 置顶状态 1 置顶 0 不置顶
	Avatar               string              `json:"avatar" xorm:"-" example:"头像"`                             // 用户头像
	UserId               string              `json:"user_id" example:"用户id"`                          // 用户id
	UserName             string              `json:"user_name" xorm:"-" example:"用户昵称"`                       // 用户名称
	CommentLevel         int                 `json:"comment_level" example:"1"`                         // 评论等级[1 一级评论 默认 ，2 二级评论]
	Content              string              `json:"content" example:"评论内容"`                          // 评论内容
	CreateAt             int                 `json:"create_at" example:"1600000000"`                    // 创建时间
	ParentCommentId      int64               `json:"parent_comment_id" example:"1000000000"`            // 父评论id
	ParentCommentUserId  string              `json:"parent_comment_user_id" example:"父评论的用户id"`     // 父评论的用户id
	ReplyCommentId       int64               `json:"reply_comment_id" example:"1000000000"`             // 被回复的评论id
	ReplyCommentUserId   string              `json:"reply_comment_user_id" example:"被回复的评论用户id"`   // 被回复的评论用户id
	ReplyCommentUserName string              `json:"reply_comment_user_name" example:"被回复评论的用户昵称"`// 被回复评论的用户昵称
	ReplyCommentAvatar   string              `json:"reply_comment_avatar" example:"被回复评论的用户头像"`   // 被回复评论的用户头像
	ReplyContent         string              `json:"reply_content"  example:"被回复的内容"`                // 被回复的内容
	Status               int                 `json:"status" example:"1"`                                 // 状态 1 有效 0 逻辑删除
	VideoId              int64               `json:"video_id,omitempty" example:"1000000000"`            // 视频id
	PostId               int64               `json:"post_id,omitempty" example:"1000000000"`             // 帖子id
	LikeNum              int64               `json:"like_num" example:"100"`                             // 点赞数
	IsAttention          int                 `json:"is_attention" example:"0"`                           // 是否关注
	IsLike               int                 `json:"is_like" example:"0"`                                // 是否点赞
	IsAt                 int                 `json:"is_at" example:"0"`                                  // 是否为@消息 0不是@消息 1是
}

// 视频评论数据（后台展示）
type VideoCommentInfo struct {
	Id            int64                 `json:"id"`              // 评论id
	VideoId       int64                 `json:"video_id"`        // 视频id
	Title         string                `json:"title"`           // 标题
	Describe      string                `json:"describe"`        // 描述
	Cover         string                `json:"cover"`           // 封面
	VideoAddr     string                `json:"video_addr"`      // 视频地址
	VideoDuration int                   `json:"video_duration"`  // 视频时长
	VideoWidth    int64                 `json:"video_width"`     // 视频宽
	VideoHeight   int64                 `json:"video_height"`    // 视频高
	Status        int32                 `json:"status"`          // 评论状态 (1 有效 0 逻辑删除)
	UserId        string                `json:"user_id"`         // 用户id
	Content       string                `json:"content"`         // 评论/回复的内容
	CreateAt      int                   `json:"create_at"`       // 用户评论的时间
	CommentLevel  int                   `json:"comment_level"`   // 1 为评论 2 为回复
	LikeNum       int64                 `json:"like_num"`        // 点赞数
	ReplyNum      int64                 `json:"reply_num"`       // 当前评论的回复数
}


// 发布评论请求参数
type PublishCommentParams struct {
	VideoId          int64       `binding:"required" json:"video_id"`      // 视频id
	Content          string      `binding:"required" json:"content"`       // 评论的内容
}

// 新版发布评论请求参数
type V2PubCommentParams struct {
	ComposeId       int64        `binding:"required" json:"compose_id"`    // 作品id 视频/帖子id
	Content         string       `binding:"required" json:"content"`       // 评论的内容
	CommentType     int          `binding:"required" json:"comment_type"`  // 评论的类型 1视频 2帖子
	AtInfo          []string     `json:"at_info"`                          // @信息 [用户uid]
}

// 回复评论请求参数
type ReplyCommentParams struct {
	VideoId          int64       `binding:"required" json:"video_id"`      // 视频id todo: 改为视频/帖子id
	Content          string      `binding:"required" json:"content"`       // 评论的内容
	ReplyId          string      `binding:"required" json:"reply_id"`      // 被回复的评论id
	CommentType      int         `json:"comment_type"`                     // 评论的类型 1视频 2帖子
	AtInfo           []string    `json:"at_info"`                          // @信息 [用户uid]
}

// 后台删除评论
type DelCommentParam struct {
	CommentId      string     `binding:"required" json:"comment_id"`       // 评论id
}

// 评论举报
type CommentReportParam struct {
	CommentId    int64      `json:"comment_id" binding:"required"`
	UserId       string     `json:"user_id"`
	Reason       string     `json:"reason"`
	CommentType  int        `json:"comment_type"`                    // 1 视频评论 2 帖子评论
}

// 实栗
func NewCommentModel(engine *xorm.Session) *CommentModel {
	return &CommentModel{
		Engine:       engine,
		VideoComment: new(models.VideoComment),
		ReceiveAt:    new(models.ReceivedAt),
		Report:       new(models.CommentReport),
	}
}

// 添加视频评论(包含回复评论)
func (m *CommentModel) AddVideoComment() error {
	if _, err := m.Engine.InsertOne(m.VideoComment); err != nil {
		return err
	}

	return nil
}

// 添加帖子评论
func (m *CommentModel) AddPostComment() error {
	if _, err := m.Engine.InsertOne(m.PostComment); err != nil {
		return err
	}

	return nil
}

// 添加用户收到的@
func (m *CommentModel) AddReceiveAt() error {
	if _, err := m.Engine.InsertOne(m.ReceiveAt); err != nil {
		return err
	}

	return nil
}

// 查询用户收到的@们
func (m *CommentModel) GetReceiveAtList(userId string, offset, size int) []*models.ReceivedAt {
	var list []*models.ReceivedAt
	if err := m.Engine.Where("to_user_id=?", userId).Desc("id").Limit(size, offset).Find(&list); err != nil {
		log.Log.Errorf("comment_trace: get receive at by userid err:%s, userId:%s", err, userId)
		return nil
	}

	return list
}

// 获取用户总评论数
func (m *CommentModel) GetUserTotalComments(userId string) int64 {
	count, err := m.Engine.Where("user_id=?", userId).Count(&models.VideoComment{})
	if err != nil {
		log.Log.Errorf("comment_trace: get user total comments err:%s", err)
		return 0
	}

	return count
}


// 获取未读的被@的数量
func (m *CommentModel) GetUnreadAtCount(userId, readTm string) int64 {
	count, err := m.Engine.Where("to_user_id=? AND create_at > ?", userId, readTm).Count(&models.ReceivedAt{})
	if err != nil {
		log.Log.Errorf("comment_trace: get unread at count err:%s", err)
		return 0
	}

	return count
}

// 通过评论id获取评论信息
func (m *CommentModel) GetVideoCommentById(commentId string) *models.VideoComment {
	comment := new(models.VideoComment)
	ok, err := m.Engine.Where("id=?", commentId).Get(comment)
	if !ok || err != nil {
		log.Log.Errorf("comment_trace: video comment not found, commentId:%s", commentId)
		return nil
	}

	// 已逻辑删除
	if comment.Status == 0 {
		comment.Content = "原内容已删除"
	}

	return comment
}

// 通过评论id 查询该评论下的所有回复id
func (m *CommentModel) GetVideoReplyIdsById(commentId string) []string {
	var replyIds []string
	if err := m.Engine.Table(&models.VideoComment{}).Cols("id").Where("reply_comment_id=?", commentId).Find(&replyIds); err != nil {
		log.Log.Errorf("comment_trace: get video reply ids err:%s", err)
	}

	return replyIds
}

// 通过评论id获取帖子评论信息
func (m *CommentModel) GetPostCommentById(commentId string) *models.PostingComment {
	comment := new(models.PostingComment)
	ok, err := m.Engine.Where("id=?", commentId).Get(comment)
	if !ok || err != nil {
		log.Log.Errorf("comment_trace: post comment not found, commentId:%s", commentId)
		return nil
	}

	// 已逻辑删除
	if comment.Status == 0 {
		comment.Content = "原内容已删除"
	}

	return comment
}

// 删除视频评论
func (m *CommentModel) DelVideoComments(commentIds string) error {
	sql := fmt.Sprintf("DELETE FROM `video_comment` WHERE `id` IN (%s)", commentIds)
	if _, err := m.Engine.Exec(sql); err != nil {
		log.Log.Errorf("comment_trace: delete comments by ids err:%s", err)
		return err
	}

	return nil
}

// 获取视频评论列表(1级评论)
func (m *CommentModel) GetVideoCommentList(composeId string, offset, size int) []*models.VideoComment {
	var list []*models.VideoComment
	if err := m.Engine.Where("video_id=? AND comment_level=1", composeId).
		Desc("is_top").
		Asc("id").
		Limit(size, offset).
		Find(&list); err != nil {
		log.Log.Errorf("comment_trace: get  comment list err:%s", err)
		return nil
	}

	return list
}

// 获取帖子评论列表(1级评论)
func (m *CommentModel) GetPostCommentList(composeId string, offset, size int) []*models.PostingComment {
	var list []*models.PostingComment
	if err := m.Engine.Where("post_id=? AND comment_level=1", composeId).
		Desc("is_top").
		Asc("id").
		Limit(size, offset).
		Find(&list); err != nil {
		log.Log.Errorf("comment_trace: get  comment list err:%s", err)
		return nil
	}

	return list
}

// 根据评论点赞数排序 获取视频评论列表（1级评论）
func (m *CommentModel) GetVideoCommentListByLike(composeId string, zanType, offset, size int) []*CommentList {
	sql := "SELECT vc.*, count(tu.Id) AS like_num FROM video_comment AS vc " +
		"LEFT JOIN thumbs_up AS tu ON vc.id = tu.type_id AND tu.zan_type=? AND tu.status=1 WHERE vc.video_id=? " +
		"AND vc.comment_level = 1 " +
		"GROUP BY vc.Id ORDER BY like_num DESC, vc.id DESC LIMIT ?, ?"

	var list []*CommentList
	if err := m.Engine.SQL(sql, zanType, composeId, offset, size).Find(&list); err != nil {
		log.Log.Errorf("comment_trace: get video comment list by like err:%s", err)
		return nil
	}

	return list
}

// 根据评论点赞数排序 获取帖子评论列表（1级评论）
func (m *CommentModel) GetPostCommentListByLike(composeId string, zanType, offset, size int) []*CommentList {
	sql := "SELECT pc.*, count(tu.Id) AS like_num FROM post_comment AS pc " +
		"LEFT JOIN thumbs_up AS tu ON pc.id = tu.type_id AND tu.zan_type=? AND tu.status=1 WHERE pc.post_id=? " +
		"AND pc.comment_level = 1 " +
		"GROUP BY pc.Id ORDER BY like_num DESC, pc.id DESC LIMIT ?, ?"

	var list []*CommentList
	if err := m.Engine.SQL(sql, zanType, composeId, offset, size).Find(&list); err != nil {
		log.Log.Errorf("comment_trace: get post comment list by like err:%s", err)
		return nil
	}

	return list
}

// 获取视频评论下的回复列表
func (m *CommentModel) GetVideoReplyList(videoId, commentId string, offset, size int) []*ReplyComment {
	var list []*ReplyComment
	if err := m.Engine.Table(&models.VideoComment{}).Where("video_id=? AND comment_level=2 AND parent_comment_id=?", videoId, commentId).
		Asc("id").
		Limit(size, offset).
		Find(&list); err != nil {
		log.Log.Errorf("comment_trace: get reply list err:%s", err)
		return nil
	}

	return list
}

// 获取帖子评论下的回复列表
func (m *CommentModel) GetPostReplyList(videoId, commentId string, offset, size int) []*ReplyComment {
	var list []*ReplyComment
	if err := m.Engine.Table(&models.PostingComment{}).Where("post_id=? AND comment_level=2 AND parent_comment_id=?", videoId, commentId).
		Asc("id").
		Limit(size, offset).
		Find(&list); err != nil {
		log.Log.Errorf("comment_trace: get reply list err:%s", err)
		return nil
	}

	return list
}

// 获取视频评论总回复数
func (m *CommentModel) GetTotalReplyByVideoComment(commentId string) int64 {
	total, err := m.Engine.Where("parent_comment_id=?", commentId).Count(&models.VideoComment{})
	if err != nil {
		log.Log.Errorf("comment_trace get total reply by comment err:%s", err)
		return 0
	}

	return total

}

// 获取帖子评论总回复数
func (m *CommentModel) GetTotalReplyByPostComment(commentId string) int64 {
	total, err := m.Engine.Where("parent_comment_id=?", commentId).Count(&models.PostingComment{})
	if err != nil {
		log.Log.Errorf("comment_trace get total reply by post comment err:%s", err)
		return 0
	}

	return total

}


// 后台获取视频评论列表（查询所有 或 用户id或视频id进行查询 可通过 1 时间、 2 点赞数、 3 回复数排序 默认时间倒序）
func (m *CommentModel) GetVideoCommentsBySort(userId, videoId, sortType, condition string, offset, size int) []*VideoCommentInfo {
	sql := "SELECT vc.*, count(distinct(tu.Id)) AS like_num, count(vc2.id) AS reply_num FROM video_comment AS vc " +
		"LEFT JOIN thumbs_up AS tu ON vc.id = tu.type_id AND tu.zan_type=3 AND tu.status=1 " +
		"LEFT JOIN video_comment AS vc2 ON vc.id=vc2.reply_comment_id AND vc2.comment_level=2 " +
		"WHERE vc.status=1 "

	if userId != "" {
		sql += fmt.Sprintf("AND vc.user_id=%s ", userId)
	}

	if videoId != "" {
		sql += fmt.Sprintf("AND vc.video_id=%s ", videoId)
	}

	sql += "GROUP BY vc.id "

	switch condition {
	case consts.SORT_BY_TIME:
		sql += "ORDER BY vc.create_at "
	case consts.SORT_BY_LIKE:
		sql += "ORDER BY like_num "
	case consts.SORT_BY_REPLY:
		sql += "ORDER BY reply_num "
	default:
		sql += "ORDER BY vc.create_at "
	}

	// 1 正序 默认倒序
	if sortType == "1" {
		sql += "ASC "
	} else {
		sql += "DESC "
	}

	sql += "LIMIT ?, ?"
	var list []*VideoCommentInfo
	if err := m.Engine.Table(&models.VideoComment{}).SQL(sql, offset, size).Find(&list); err != nil {
		log.Log.Errorf("comment_trace: get comment list by sort, err:%s", err)
		return []*VideoCommentInfo{}
	}

	return list
}

// 获取评论总数
func (m *CommentModel) GetVideoCommentTotal() int64 {
	count, err := m.Engine.Where("status=1").Count(&models.VideoComment{})
	if err != nil {
		log.Log.Errorf("comment_trace: get total comments err:%s", err)
		return 0
	}

	return count
}

// 通过用户id 获取视频评论总数
func (m *CommentModel) GetCommentTotalByUserId(userId string) int64 {
	count, err := m.Engine.Where("status=1 AND user_id=?", userId).Count(&models.VideoComment{})
	if err != nil {
		return 0
	}

	return count
}

// 通过视频id 获取评论总数
func (m *CommentModel) GetCommentTotalByVideoId(composeId string) int64 {
	count, err := m.Engine.Where("status=1 AND video_id=?", composeId).Count(&models.VideoComment{})
	if err != nil {
		return 0
	}

	return count
}

// 添加评论举报
func (m *CommentModel) AddCommentReport() (int64, error) {
	return m.Engine.InsertOne(m.Report)
}

// 更新评论/回复 信息
func (m *CommentModel) UpdateCommentInfo(condition, cols string) (int64, error) {
	return m.Engine.Where(condition).Cols(cols).Update(m.VideoComment)
}
