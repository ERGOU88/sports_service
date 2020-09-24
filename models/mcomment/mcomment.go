package mcomment

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
)

type CommentModel struct {
	Engine      *xorm.Session
	Comment     *models.VideoComment
	ReceiveAt   *models.ReceivedAt
}

// 视频评论列表
type VideoComments struct {
	Id                  int64               `json:"id"`                      // 评论id
	IsTop               int                 `json:"is_top"`                  // 置顶状态 1 置顶 0 不置顶
	Avatar              string              `json:"avatar"`                  // 用户头像
	UserId              string              `json:"user_id"`                 // 用户id
	UserName            string              `json:"user_name"`               // 用户名称
	CommentLevel        int                 `json:"comment_level"`           // 评论等级[1 一级评论 默认 ，2 二级评论]
	Content             string              `json:"content"`                 // 评论内容
	CreateAt            int                 `json:"create_at"`               // 创建时间
	Status              int                 `json:"status"`                  // 状态 1 有效 0 逻辑删除
	VideoId             int64               `json:"video_id"`                // 视频id
	ReplyList           []*ReplyComment     `json:"reply_list"`              // 回复列表
	LikeNum             int64               `json:"like_num"`                // 点赞数
	IsAttention         int                 `json:"is_attention"`            // 是否关注
	ReplyNum            int64               `json:"reply_num"`               // 总回复数
}

// 回复评论的内容
type ReplyComment struct {
	Id                   int64               `json:"id"`                      // 评论id
	IsTop                int                 `json:"is_top"`                  // 置顶状态 1 置顶 0 不置顶
	Avatar               string              `json:"avatar"`                  // 用户头像
	UserId               string              `json:"user_id"`                 // 用户id
	UserName             string              `json:"user_name"`               // 用户名称
	CommentLevel         int                 `json:"comment_level"`           // 评论等级[1 一级评论 默认 ，2 二级评论]
	Content              string              `json:"content"`                 // 评论内容
	CreateAt             int                 `json:"create_at"`               // 创建时间
	ParentCommentId      int64               `json:"parent_comment_id"`       // 父评论id
	ParentCommentUserId  string              `json:"parent_comment_user_id"`  // 父评论的用户id
	ReplyCommentId       int64               `json:"reply_comment_id"`        // 被回复的评论id
	ReplyCommentUserId   string              `json:"reply_comment_user_id"`   // 被回复的评论用户id
	ReplyCommentUserName string              `json:"reply_comment_user_name"` // 被回复评论的用户昵称
	ReplyCommentAvatar   string              `json:"reply_comment_avatar"`    // 被回复评论的用户头像
	ReplyContent         string              `json:"reply_content"`           // 被回复的内容
	Status               int                 `json:"status"`                  // 状态 1 有效 0 逻辑删除
	VideoId              int64               `json:"video_id"`                // 视频id
	LikeNum              int64               `json:"like_num"`                // 点赞数
	IsAttention          int                 `json:"is_attention"`            // 是否关注
}

// 发布评论请求参数
type PublishCommentParams struct {
	VideoId          int64       `json:"video_id"`      // 视频id
	Content          string      `json:"content"`       // 评论的内容
}

// 回复评论请求参数
type ReplyCommentParams struct {
	VideoId          int64       `json:"video_id"`      // 视频id
	Content          string      `json:"content"`       // 评论的内容
	ReplyId          string      `json:"reply_id"`      // 被回复的评论id
}

// 实栗
func NewCommentModel(engine *xorm.Session) *CommentModel {
	return &CommentModel{
		Engine:  engine,
		Comment: new(models.VideoComment),
		ReceiveAt: new(models.ReceivedAt),
	}
}

// 添加视频评论(包含回复评论)
func (m *CommentModel) AddVideoComment() error {
	if _, err := m.Engine.InsertOne(m.Comment); err != nil {
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
	ok, err := m.Engine.Where("id=? AND status=1", commentId).Get(comment)
	if !ok || err != nil {
		log.Log.Errorf("comment_trace: video comment not found, commentId:%s", commentId)
		return nil
	}

	return comment
}


// 获取视频评论列表(1级评论)
func (m *CommentModel) GetVideoCommentList(videoId string, offset, size int) []*models.VideoComment {
	var list []*models.VideoComment
	if err := m.Engine.Where("video_id=? AND comment_level=1 AND status=1", videoId).
		Desc("is_top").
		Asc("id").
		Limit(size, offset).
		Find(&list); err != nil {
			log.Log.Errorf("comment_trace: get video comment err:%s", err)
		return nil
	}

	return list
}

// 根据评论点赞数排序 获取视频评论列表（1级评论）
func (m *CommentModel) GetVideoCommentListByLike(videoId string, offset, size int) []*VideoComments {
	sql := "SELECT vc.*, count(tu.Id) AS like_num FROM video_comment AS vc " +
		"LEFT JOIN thumbs_up AS tu ON vc.id = tu.type_id AND tu.zan_type=3 AND tu.status=1 WHERE vc.video_id=? " +
		"AND vc.comment_level = 1 AND vc.status=1 " +
		"GROUP BY vc.Id ORDER BY like_num DESC LIMIT ?, ?"

	var list []*VideoComments
	if err := m.Engine.SQL(sql, videoId, offset, size).Find(&list); err != nil {
		log.Log.Errorf("comment_trace: get video comment list by like err:%s", err)
		return nil
	}

	return list
}



// 获取评论下的回复列表
func (m *CommentModel) GetVideoReply(videoId, commentId string, offset, size int) []*ReplyComment {
	var list []*ReplyComment
	if err := m.Engine.Table(&models.VideoComment{}).Where("video_id=? AND comment_level=2 AND parent_comment_id=? AND status=1", videoId, commentId).
		Desc("id").
		Limit(size, offset).
		Find(&list); err != nil {
			log.Log.Errorf("comment_trace: get video reply err:%s", err)
		return nil
	}

	return list
}

// 获取评论总回复数
func (m *CommentModel) GetTotalReplyByComment(commentId string) int64 {
	total, err := m.Engine.Where("parent_comment_id=? AND status=1", commentId).Count(&models.VideoComment{})
	if err != nil {
		log.Log.Errorf("comment_trace get total reply by comment err:%s", err)
		return 0
	}

	return total

}
