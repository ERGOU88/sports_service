package comment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"sports_service/server/util"
	"time"
)

type CommentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	comment     *mcomment.CommentModel
	collect     *mcollect.CollectModel
	user        *muser.UserModel
	video       *mvideo.VideoModel
	like        *mlike.LikeModel
	attention   *mattention.AttentionModel
}

func New(c *gin.Context) CommentModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return CommentModule{
		context: c,
		comment: mcomment.NewCommentModel(socket),
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		like: mlike.NewLikeModel(socket),
		attention: mattention.NewAttentionModel(socket),
		engine: socket,
	}
}

// 发布评论 todo: 脏词过滤
func (svc *CommentModule) PublishComment(userId string, params *mcomment.PublishCommentParams) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("video_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("comment_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 最少10字符 最多1000字符
	contentLen := util.GetStrLen([]rune(params.Content))
	if contentLen < consts.COMMENT_MIN_LEN || contentLen > consts.COMMENT_MAX_LEN {
		log.Log.Errorf("comment_trace: invalid content length, len:%d", contentLen)
		svc.engine.Rollback()
		return errdef.COMMENT_INVALID_CONTENT
	}

	// 查找视频是否存在
	video := svc.video.FindVideoById(fmt.Sprint(params.VideoId))
	if video == nil {
		log.Log.Errorf("comment_trace: video not found, videoId:%d", params.VideoId)
		svc.engine.Rollback()
		return errdef.VIDEO_NOT_EXISTS
	}

	now := time.Now().Unix()
	svc.comment.Comment.UserId = userId
	svc.comment.Comment.Content = params.Content
	svc.comment.Comment.CommentLevel = consts.COMMENT_PUBLISH
	svc.comment.Comment.Avatar = user.Avatar
	svc.comment.Comment.UserName = user.NickName
	svc.comment.Comment.CreateAt = int(now)
	svc.comment.Comment.ParentCommentId = 0
	svc.comment.Comment.VideoId = params.VideoId
	svc.comment.Comment.Status = 1
	// 添加视频评论
	if err := svc.comment.AddVideoComment(); err != nil {
		log.Log.Errorf("comment_trace: add video comment err:%s", err)
		svc.engine.Rollback()
		return errdef.COMMENT_PUBLISH_FAIL
	}

	svc.comment.ReceiveAt.UserId = userId
	// 被@的用户 1级评论 则@的是视频up主
	svc.comment.ReceiveAt.ToUserId = video.UserId
	svc.comment.ReceiveAt.CommentId = svc.comment.Comment.Id
	svc.comment.ReceiveAt.TopicType = consts.TYPE_COMMENT
	svc.comment.ReceiveAt.CreateAt = int(now)
	// 评论也是@
	if err := svc.comment.AddReceiveAt(); err != nil {
		log.Log.Errorf("comment_trace: add receive at err:%s", err)
		svc.engine.Rollback()
		return errdef.COMMENT_PUBLISH_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 回复评论 todo: 脏词过滤
func (svc *CommentModule) PublishReply(userId string, params *mcomment.ReplyCommentParams) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("video_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("comment_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 最少10字符 最多1000字符
	contentLen := util.GetStrLen([]rune(params.Content))
	if contentLen < consts.COMMENT_MIN_LEN || contentLen > consts.COMMENT_MAX_LEN {
		log.Log.Errorf("comment_trace: invalid content length, len:%d", contentLen)
		svc.engine.Rollback()
		return errdef.COMMENT_INVALID_CONTENT
	}

	// 查找视频是否存在
	if video := svc.video.FindVideoById(fmt.Sprint(params.VideoId)); video == nil {
		log.Log.Errorf("comment_trace: video not found, videoId:%d", params.VideoId)
		svc.engine.Rollback()
		return errdef.VIDEO_NOT_EXISTS
	}

	// 查询被回复的评论是否存在
	replyInfo := svc.comment.GetVideoCommentById(params.ReplyId)
	if replyInfo == nil {
		log.Log.Error("comment_trace: reply comment not found, commentId:%s", params.ReplyId)
		svc.engine.Rollback()
		return errdef.COMMENT_NOT_FOUND
	}

	// todo: 用户是否能回复自己？
	//if strings.Compare(userId, replyInfo.UserId) != -1 {
	//
	//}
	now := time.Now().Unix()
	svc.comment.Comment.UserId = userId
	svc.comment.Comment.Content = params.Content
	svc.comment.Comment.Avatar = user.Avatar
	svc.comment.Comment.UserName = user.NickName
	svc.comment.Comment.CreateAt = int(now)
	svc.comment.Comment.VideoId = params.VideoId
	svc.comment.Comment.CommentLevel = consts.COMMENT_REPLY
	svc.comment.Comment.Status = 1

	svc.comment.Comment.ReplyCommentUserId = replyInfo.UserId
	svc.comment.Comment.ReplyCommentId = replyInfo.Id
	svc.comment.Comment.ParentCommentUserId = replyInfo.ParentCommentUserId
	svc.comment.Comment.ParentCommentId = replyInfo.ParentCommentId
	// 1级评论 parentid为0 如果被回复的评论 parentid 为0，说明当前回复的是1级评论 否则 回复的为2级评论
	if replyInfo.ParentCommentId == 0 {
		svc.comment.Comment.ParentCommentId = replyInfo.Id
		svc.comment.Comment.ParentCommentUserId = replyInfo.UserId
	}

	if err := svc.comment.AddVideoComment(); err != nil {
		log.Log.Errorf("comment_trace: add video reply err:%s", err)
		svc.engine.Rollback()
		return errdef.COMMENT_REPLY_FAIL
	}

	svc.comment.ReceiveAt.UserId = userId
	// 被@的用户
	svc.comment.ReceiveAt.ToUserId = replyInfo.UserId
	svc.comment.ReceiveAt.CommentId = svc.comment.Comment.Id
	svc.comment.ReceiveAt.TopicType = consts.TYPE_COMMENT
	svc.comment.ReceiveAt.CreateAt = int(now)
	// 回复 记录到 @
	if err := svc.comment.AddReceiveAt(); err != nil {
		log.Log.Errorf("comment_trace: add receive at err:%s", err)
		svc.engine.Rollback()
		return errdef.COMMENT_REPLY_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 获取视频评论
func (svc *CommentModule) GetVideoComment(userId, videoId, sortType string, page, size int) []*mcomment.VideoComments {
	video := svc.video.FindVideoById(videoId)
	// 视频不存在 或 视频未过审
	if video == nil || fmt.Sprint(video.Status) != consts.VIDEO_AUDIT_SUCCESS {
		log.Log.Errorf("comment_trace: video not found or not pass, videoId:%s", videoId)
		return nil
	}

	offset := (page - 1) * size
	// 获取评论
	comments := svc.comment.GetVideoCommentList(videoId, offset, size)
	if len(comments) == 0 {
		log.Log.Errorf("comment_trace: no comments, videoId:%s", videoId)
		return nil
	}

	list := make([]*mcomment.VideoComments, len(comments))
	type tmpUser struct {
		NickName   string
		Avatar     string
	}
	// contents 存储 评论id——>评论的内容
	contents := make(map[int64]string, 0)
	// userInfo 存储 用户id——>头像、昵称
	userInfo := make(map[string]*tmpUser, 0)
	for index, item := range comments {
		comment := new(mcomment.VideoComments)
		comment.Id = item.Id
		comment.Status = item.Status
		comment.VideoId = item.VideoId
		comment.UserId = item.UserId
		comment.Avatar = item.Avatar
		comment.UserName = item.UserName
		comment.CreateAt = item.CreateAt
		comment.IsTop = item.IsTop
		comment.Content = item.Content
		comment.CommentLevel = item.CommentLevel
		// 总回复数
		comment.ReplyNum = svc.comment.GetTotalReplyByComment(fmt.Sprint(comment.Id))
		// 评论点赞数
		comment.LikeNum = svc.like.GetLikeNumByType(item.Id, consts.TYPE_COMMENT)

		user := new(tmpUser)
		user.NickName = item.UserName
		user.Avatar = item.Avatar
		userInfo[item.UserId] = user

		contents[item.Id] = item.Content

		// 获取每个评论下的回复列表 (默认取三条)
		comment.ReplyList = svc.comment.GetVideoReply(videoId, fmt.Sprint(item.Id), 0, 3)
		for _, reply := range comment.ReplyList {
			user := new(tmpUser)
			user.NickName = reply.UserName
			user.Avatar = reply.Avatar
			userInfo[reply.UserId] = user

			contents[reply.Id] = reply.Content
			// 评论点赞数
			reply.LikeNum = svc.like.GetLikeNumByType(reply.Id, consts.TYPE_COMMENT)
			// 被回复的用户名、用户头像
			uinfo, ok := userInfo[reply.ReplyCommentUserId]
			if ok {
				reply.ReplyCommentAvatar = uinfo.Avatar
				reply.ReplyCommentUserName = uinfo.NickName
			}

			// 被回复的内容
			content, ok := contents[reply.ReplyCommentId]
			if ok {
				reply.ReplyContent = content
			}

			if userId != "" {
				// 是否关注
				if attention := svc.attention.GetAttentionInfo(userId, reply.UserId); attention != nil {
					reply.IsAttention = attention.Status
				}
			}
		}

		if userId != "" {
			// 是否关注
			if attention := svc.attention.GetAttentionInfo(userId, item.UserId); attention != nil {
				comment.IsAttention = attention.Status
			}
		}

		list[index] = comment

	}

	// 热度排序（点赞数）
	if sortType == consts.SORT_HOT {
		log.Log.Debug("sort by hot")
		util.PartialSort(SortComment(list), len(list))
	}

	return list
}

// 获取评论回复列表
func (svc *CommentModule) GetCommentReplyList(userId, videoId, commentId string, page, size int) (int, []*mcomment.ReplyComment) {
	video := svc.video.FindVideoById(videoId)
	// 视频不存在 或 视频未过审
	if video == nil || fmt.Sprint(video.Status) != consts.VIDEO_AUDIT_SUCCESS {
		log.Log.Errorf("comment_trace: video not found or not pass, videoId:%s", videoId)
		return errdef.VIDEO_NOT_EXISTS, nil
	}

	// 查询评论是否存在
	comment := svc.comment.GetVideoCommentById(commentId)
	if comment == nil {
		log.Log.Error("comment_trace: comment not found, commentId:%s", commentId)
		return errdef.COMMENT_NOT_FOUND, nil
	}

	offset := (page - 1) * size
	replyList := svc.comment.GetVideoReply(videoId, commentId, offset, size)
	if len(replyList) == 0 {
		log.Log.Errorf("comment_trace: not found comment reply, commentId:%s", commentId)
		return errdef.COMMENT_REPLY_NOT_FOUND, nil
	}

	type tmpUser struct {
		NickName   string
		Avatar     string
	}
	// contents 存储 评论id——>评论的内容
	contents := make(map[int64]string, 0)
	// userInfo 存储 用户id——>头像、昵称
	userInfo := make(map[string]*tmpUser, 0)
	for _, reply := range replyList {
		user := new(tmpUser)
		user.NickName = reply.UserName
		user.Avatar = reply.Avatar
		userInfo[reply.UserId] = user

		contents[reply.Id] = reply.Content
		// 评论点赞数
		reply.LikeNum = svc.like.GetLikeNumByType(reply.Id, consts.TYPE_COMMENT)
		// 被回复的用户名、用户头像
		uinfo, ok := userInfo[reply.ReplyCommentUserId]
		if ok {
			reply.ReplyCommentAvatar = uinfo.Avatar
			reply.ReplyCommentUserName = uinfo.NickName
		}

		// 被回复的内容
		content, ok := contents[reply.ReplyCommentId]
		if ok {
			reply.ReplyContent = content
		}

		if userId != "" {
			// 是否关注
			if attention := svc.attention.GetAttentionInfo(userId, reply.UserId); attention != nil {
				reply.IsAttention = attention.Status
			}
		}
	}

	return errdef.SUCCESS, replyList
}
