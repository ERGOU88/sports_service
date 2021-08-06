package comment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mlike"
	"sports_service/server/models/mposting"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	redismq "sports_service/server/redismq/event"
	"sports_service/server/tools/tencentCloud"
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
	post        *mposting.PostingModel
}

func New(c *gin.Context) CommentModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return CommentModule{
		context: c,
		comment: mcomment.NewCommentModel(socket),
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		like: mlike.NewLikeModel(socket),
		attention: mattention.NewAttentionModel(socket),
		post: mposting.NewPostingModel(socket),
		engine: socket,
	}
}

// 新版评论
func (svc *CommentModule) V2PublishComment(userId string, params *mcomment.V2PubCommentParams) (int, *mcomment.CommentList) {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("video_trace: session begin err:%s", err)
		return errdef.ERROR, nil
	}

	//contentLen := util.GetStrLen([]rune(params.Content))
	// 最少1字符 最多1000字符
	contentLen := len(params.Content)
	if contentLen < consts.COMMENT_MIN_LEN || contentLen > consts.COMMENT_MAX_LEN {
		log.Log.Errorf("comment_trace: invalid content length, len:%d", contentLen)
		svc.engine.Rollback()
		return errdef.COMMENT_INVALID_LEN, nil
	}

	client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测评论内容
	b := util.IsSpace([]rune(params.Content))
	isPass, err := client.TextModeration(params.Content)
	if !isPass || !b {
		log.Log.Errorf("comment_trace: validate comment err: %s，pass: %v", err, isPass)
		svc.engine.Rollback()
		return errdef.COMMENT_INVALID_CONTENT, nil
	}

	// 查询用户是否存在
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("comment_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS, nil
	}

	now := int(time.Now().Unix())
	var (
		cover, toUserId string
		msgType int32
		commentId, composeId int64
	)

	resp := &mcomment.CommentList{}
	switch params.CommentType {
	// 视频评论
	case consts.COMMENT_TYPE_VIDEO:
		svc.comment.VideoComment.UserId = userId
		svc.comment.VideoComment.Content = params.Content
		svc.comment.VideoComment.CommentLevel = consts.COMMENT_PUBLISH
		//svc.comment.VideoComment.Avatar = user.Avatar
		//svc.comment.VideoComment.UserName = user.NickName
		svc.comment.VideoComment.CreateAt = now
		svc.comment.VideoComment.ParentCommentId = 0
		svc.comment.VideoComment.VideoId = params.ComposeId
		svc.comment.VideoComment.Status = 1
		// 添加评论
		if err := svc.comment.AddVideoComment(); err != nil {
			log.Log.Errorf("comment_trace: add comment err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_PUBLISH_FAIL, nil
		}

		// 查找视频是否存在
		video := svc.video.FindVideoById(fmt.Sprint(params.ComposeId))
		if video == nil  {
			log.Log.Errorf("comment_trace: video not found, videoId:%d", params.ComposeId)
			svc.engine.Rollback()
			return errdef.VIDEO_NOT_EXISTS, nil
		}

		// 视频状态 != 1 (1为视频审核成功)
		if fmt.Sprint(video.Status) != consts.VIDEO_AUDIT_SUCCESS {
			log.Log.Errorf("comment_trace: video status not audit success, videoId:%d", params.ComposeId)
			svc.engine.Rollback()
			return errdef.VIDEO_NOT_EXISTS, nil
		}

		// 更新视频总计（视频评论总数）
		if err := svc.video.UpdateVideoCommentNum(params.ComposeId, int(now), 1); err != nil {
			log.Log.Errorf("comment_trace: update video comment num err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_PUBLISH_FAIL, nil
		}

		svc.comment.ReceiveAt.TopicType = consts.TYPE_VIDEO_COMMENT
		cover = video.Cover
		msgType = consts.VIDEO_COMMENT_MSG
		toUserId = video.UserId
		commentId = svc.comment.VideoComment.Id
		composeId = video.VideoId

		resp.Id = svc.comment.VideoComment.Id
		resp.Avatar = user.Avatar
		resp.UserName = user.NickName
		resp.VideoId = params.ComposeId
		resp.CommentLevel = svc.comment.VideoComment.CommentLevel
		resp.UserId = user.UserId
		resp.Content = svc.comment.VideoComment.Content
		resp.CreateAt = svc.comment.VideoComment.CreateAt
		resp.Status = svc.comment.VideoComment.Status

	// 帖子评论
	case consts.COMMENT_TYPE_POST:
		svc.comment.PostComment.UserId = userId
		svc.comment.PostComment.Content = params.Content
		svc.comment.PostComment.CommentLevel = consts.COMMENT_PUBLISH
		svc.comment.PostComment.CreateAt = now
		svc.comment.PostComment.ParentCommentId = 0
		svc.comment.PostComment.PostId = params.ComposeId
		svc.comment.PostComment.Status = 1
		// 添加评论
		if err := svc.comment.AddPostComment(); err != nil {
			log.Log.Errorf("comment_trace: add comment err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_PUBLISH_FAIL, nil
		}

		// 查找帖子是否存在
		post, err := svc.post.GetPostById(fmt.Sprint(params.ComposeId))
		if post == nil || err != nil {
			log.Log.Errorf("comment_trace: post not found, postId:%d", params.ComposeId)
			svc.engine.Rollback()
			return errdef.POST_NOT_EXISTS, nil
		}


		// 状态 != 1 (1为帖子审核成功)
		if fmt.Sprint(post.Status) != consts.POST_AUDIT_SUCCESS {
			log.Log.Errorf("comment_trace: post status not audit success, postId:%d", params.ComposeId)
			svc.engine.Rollback()
			return errdef.VIDEO_NOT_EXISTS, nil
		}

		// 更新帖子总计（帖子评论总数）
		if err := svc.post.UpdatePostCommentNum(params.ComposeId, int(now), 1); err != nil {
			log.Log.Errorf("comment_trace: update post comment num err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_PUBLISH_FAIL, nil
		}

		svc.comment.ReceiveAt.TopicType = consts.TYPE_POST_COMMENT
		msgType = consts.POST_COMMENT_MSG
		toUserId = post.UserId
		commentId = svc.comment.PostComment.Id
		composeId = post.Id

		resp.Id = svc.comment.PostComment.Id
		resp.Avatar = user.Avatar
		resp.UserName = user.NickName
		resp.PostId = params.ComposeId
		resp.CommentLevel = svc.comment.PostComment.CommentLevel
		resp.UserId = user.UserId
		resp.Content = params.Content
		resp.CreateAt = svc.comment.PostComment.CreateAt
		resp.Status = svc.comment.PostComment.Status
	default:
		log.Log.Errorf("comment_trace: invalid commentType:%d", params.CommentType)
		return errdef.INVALID_PARAMS, nil
	}

	svc.comment.ReceiveAt.UserId = userId
	// 被@的用户 1级评论 则@的是视频up主 / 帖子发布者
	svc.comment.ReceiveAt.ToUserId = toUserId
	svc.comment.ReceiveAt.ComposeId = commentId
	svc.comment.ReceiveAt.UpdateAt = now
	svc.comment.ReceiveAt.CommentLevel = consts.COMMENT_PUBLISH
	svc.comment.ReceiveAt.Status = 1
	// 评论也是@
	//if err := svc.comment.AddReceiveAt(); err != nil {
	//	log.Log.Errorf("comment_trace: add receive at err:%s", err)
	//	svc.engine.Rollback()
	//	return errdef.COMMENT_PUBLISH_FAIL, 0
	//}

	// 添加@
	if len(params.AtInfo) > 0 {
		atList := make([]*models.ReceivedAt, 0)
		atList = append(atList, svc.comment.ReceiveAt)
		for _, val := range params.AtInfo {
			user := svc.user.FindUserByUserid(val)
			if user == nil {
				log.Log.Errorf("post_trace: at user not found, userId:%s", val)
				continue
			}

			at := &models.ReceivedAt{
				ToUserId:     val,
				UserId:       userId,
				ComposeId:    commentId,
				TopicType:    params.CommentType,
				UpdateAt:     now,
				CommentLevel: consts.COMMENT_PUBLISH,
				CreateAt:     now,
				Status:       1,
			}

			atList = append(atList, at)
		}

		affected, err := svc.post.AddReceiveAtList(atList)
		if err != nil || int(affected) != len(atList) {
			log.Log.Errorf("post_trace: add receive at list fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_PUBLISH_FAIL, nil
		}
	}

	svc.engine.Commit()

	// 视频/帖子 评论推送
	//event.PushEventMsg(config.Global.AmqpDsn, toUserId, user.NickName, cover, params.Content, msgType)
	redismq.PushEventMsg(redismq.NewEvent(toUserId, fmt.Sprint(composeId), user.NickName, cover, params.Content, msgType))

	return errdef.SUCCESS, resp
}

// 发布评论
func (svc *CommentModule) PublishComment(userId string, params *mcomment.PublishCommentParams) (int, int64) {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("video_trace: session begin err:%s", err)
		return errdef.ERROR, 0
	}

	//contentLen := util.GetStrLen([]rune(params.Content))
	// 最少1字符 最多1000字符
	contentLen := len(params.Content)
	if contentLen < consts.COMMENT_MIN_LEN || contentLen > consts.COMMENT_MAX_LEN {
		log.Log.Errorf("comment_trace: invalid content length, len:%d", contentLen)
		svc.engine.Rollback()
		return errdef.COMMENT_INVALID_LEN, 0
	}

	client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测评论内容
	b := util.IsSpace([]rune(params.Content))
	isPass, err := client.TextModeration(params.Content)
	if !isPass || !b {
		log.Log.Errorf("comment_trace: validate comment err: %s，pass: %v", err, isPass)
		svc.engine.Rollback()
		return errdef.COMMENT_INVALID_CONTENT, 0
	}

	// 查询用户是否存在
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("comment_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS, 0
	}

	// 查找视频是否存在
	video := svc.video.FindVideoById(fmt.Sprint(params.VideoId))
	if video == nil {
		log.Log.Errorf("comment_trace: video not found, videoId:%d", params.VideoId)
		svc.engine.Rollback()
		return errdef.VIDEO_NOT_EXISTS, 0
	}

	// 视频状态 != 1 (视频审核成功)
	if fmt.Sprint(video.Status) != consts.VIDEO_AUDIT_SUCCESS {
		log.Log.Errorf("comment_trace: video status not audit success, videoId:%d", params.VideoId)
		svc.engine.Rollback()
		return errdef.VIDEO_NOT_EXISTS, 0
	}

	now := time.Now().Unix()
	svc.comment.VideoComment.UserId = userId
	svc.comment.VideoComment.Content = params.Content
	svc.comment.VideoComment.CommentLevel = consts.COMMENT_PUBLISH
	//svc.comment.VideoComment.Avatar = user.Avatar
	//svc.comment.VideoComment.UserName = user.NickName
	svc.comment.VideoComment.CreateAt = int(now)
	svc.comment.VideoComment.ParentCommentId = 0
	svc.comment.VideoComment.VideoId = params.VideoId
	svc.comment.VideoComment.Status = 1
	// 添加评论
	if err := svc.comment.AddVideoComment(); err != nil {
		log.Log.Errorf("comment_trace: add comment err:%s", err)
		svc.engine.Rollback()
		return errdef.COMMENT_PUBLISH_FAIL, 0
	}

	svc.comment.ReceiveAt.UserId = userId
	// 被@的用户 1级评论 则@的是视频up主
	svc.comment.ReceiveAt.ToUserId = video.UserId
	svc.comment.ReceiveAt.ComposeId = svc.comment.VideoComment.Id
	svc.comment.ReceiveAt.TopicType = consts.TYPE_VIDEO_COMMENT
	svc.comment.ReceiveAt.UpdateAt = int(now)
	svc.comment.ReceiveAt.CommentLevel = consts.COMMENT_PUBLISH
	// 评论也是@
	if err := svc.comment.AddReceiveAt(); err != nil {
		log.Log.Errorf("comment_trace: add receive at err:%s", err)
		svc.engine.Rollback()
		return errdef.COMMENT_PUBLISH_FAIL, 0
	}

	// 更新视频总计（视频评论总数）
	if err := svc.video.UpdateVideoCommentNum(video.VideoId, int(now), 1); err != nil {
		log.Log.Errorf("comment_trace: update video comment num err:%s", err)
		svc.engine.Rollback()
		return errdef.COMMENT_PUBLISH_FAIL, 0
	}

	svc.engine.Commit()

	// 视频评论推送
	//event.PushEventMsg(config.Global.AmqpDsn, video.UserId, user.NickName, video.Cover, params.Content, consts.VIDEO_COMMENT_MSG)
	redismq.PushEventMsg(redismq.NewEvent(video.UserId, fmt.Sprint(video.VideoId), user.NickName, video.Cover, params.Content, consts.VIDEO_COMMENT_MSG))

	return errdef.SUCCESS, svc.comment.VideoComment.Id
}

// 回复评论
func (svc *CommentModule) PublishReply(userId string, params *mcomment.ReplyCommentParams) (int, *mcomment.ReplyComment) {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("video_trace: session begin err:%s", err)
		return errdef.ERROR, nil
	}

	// 最少10字符 最多1000字符
	contentLen := util.GetStrLen([]rune(params.Content))
	if contentLen < consts.COMMENT_MIN_LEN || contentLen > consts.COMMENT_MAX_LEN {
		log.Log.Errorf("comment_trace: invalid content length, len:%d", contentLen)
		svc.engine.Rollback()
		return errdef.COMMENT_INVALID_LEN, nil
	}

	client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测评论内容
	isPass, err := client.TextModeration(params.Content)
	if !isPass {
		log.Log.Errorf("comment_trace: validate reply content err: %s，pass: %v", err, isPass)
		svc.engine.Rollback()
		return errdef.COMMENT_INVALID_REPLY, nil
	}

	// 查询用户是否存在
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("comment_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS, nil
	}

	// todo: 用户是否能回复自己？
	//if strings.Compare(userId, replyInfo.UserId) != -1 {
	//
	//}

	if params.CommentType <= 0 {
		// 默认为视频回复
		params.CommentType = consts.COMMENT_TYPE_VIDEO
	}

	var (
		cover, toUserId, content string
		msgType int32
		commentId, composeId int64
		atType int
	)
	now := int(time.Now().Unix())
	resp := &mcomment.ReplyComment{}
	switch params.CommentType {
	// 视频评论
	case consts.COMMENT_TYPE_VIDEO:
		// 查询被回复的评论是否存在
		replyInfo := svc.comment.GetVideoCommentById(fmt.Sprint(params.ReplyId))
		if replyInfo == nil {
			log.Log.Error("comment_trace: reply comment not found, commentId:%s", params.ReplyId)
			svc.engine.Rollback()
			return errdef.COMMENT_NOT_FOUND, nil
		}

		// 查找视频是否存在
		video := svc.video.FindVideoById(fmt.Sprint(replyInfo.VideoId))
		if video == nil  {
			log.Log.Errorf("comment_trace: video not found, videoId:%d", replyInfo.VideoId)
			svc.engine.Rollback()
			return errdef.VIDEO_NOT_EXISTS, nil
		}

		// 视频状态 != 1 (1为视频审核成功)
		if fmt.Sprint(video.Status) != consts.VIDEO_AUDIT_SUCCESS {
			log.Log.Errorf("comment_trace: video status not audit success, postId:%d", replyInfo.VideoId)
			svc.engine.Rollback()
			return errdef.VIDEO_NOT_EXISTS, nil
		}

		svc.comment.VideoComment.UserId = userId
		svc.comment.VideoComment.Content = params.Content
		//svc.comment.VideoComment.Avatar = user.Avatar
		//svc.comment.VideoComment.UserName = user.NickName
		svc.comment.VideoComment.CreateAt = now
		svc.comment.VideoComment.VideoId = replyInfo.VideoId
		svc.comment.VideoComment.CommentLevel = consts.COMMENT_REPLY
		svc.comment.VideoComment.Status = 1
		//svc.comment.VideoComment.CommentType = replyInfo.CommentType

		svc.comment.VideoComment.ReplyCommentUserId = replyInfo.UserId
		svc.comment.VideoComment.ReplyCommentId = replyInfo.Id
		svc.comment.VideoComment.ParentCommentUserId = replyInfo.ParentCommentUserId
		svc.comment.VideoComment.ParentCommentId = replyInfo.ParentCommentId

		isAt := 1
		// 1级评论 parentid为0 如果被回复的评论 parentid 为0，说明当前回复的是1级评论 否则 回复的为2级评论
		if replyInfo.ParentCommentId == 0 {
			svc.comment.VideoComment.ParentCommentId = replyInfo.Id
			svc.comment.VideoComment.ParentCommentUserId = replyInfo.UserId
			isAt = 0
		}

		if err := svc.comment.AddVideoComment(); err != nil {
			log.Log.Errorf("comment_trace: add video reply err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_REPLY_FAIL, nil
		}

		// 更新视频总计（视频评论总数）
		if err := svc.video.UpdateVideoCommentNum(replyInfo.VideoId, now, 1); err != nil {
			log.Log.Errorf("comment_trace: update video comment num err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_PUBLISH_FAIL, nil
		}

		svc.comment.ReceiveAt.TopicType = consts.TYPE_VIDEO_COMMENT
		cover = video.Cover
		msgType = consts.VIDEO_REPLY_MSG
		toUserId = replyInfo.UserId
		content = replyInfo.Content
		commentId = svc.comment.VideoComment.Id
		composeId = replyInfo.VideoId

		resp.Id = svc.comment.VideoComment.Id
		resp.Avatar = user.Avatar
		resp.UserName = user.NickName
		resp.VideoId = params.ComposeId
		resp.CommentLevel = svc.comment.VideoComment.CommentLevel
		resp.UserId = user.UserId
		resp.Content = svc.comment.VideoComment.Content
		resp.CreateAt = svc.comment.VideoComment.CreateAt
		resp.Status = svc.comment.VideoComment.Status
		resp.IsAt = isAt
		resp.ReplyContent = replyInfo.Content
		resp.ReplyCommentUserId = replyInfo.UserId
		resp.ReplyCommentId = replyInfo.Id

		if uinfo := svc.user.FindUserByUserid(replyInfo.UserId); uinfo != nil {
			resp.ReplyCommentAvatar = uinfo.Avatar
			resp.ReplyCommentUserName = uinfo.NickName
		}

		atType = consts.TYPE_VIDEO_COMMENT

	// 帖子评论
	case consts.COMMENT_TYPE_POST:
		// 查询被回复的评论是否存在
		replyInfo := svc.comment.GetPostCommentById(fmt.Sprint(params.ReplyId))
		if replyInfo == nil {
			log.Log.Error("comment_trace: reply comment not found, commentId:%s", params.ReplyId)
			svc.engine.Rollback()
			return errdef.COMMENT_NOT_FOUND, nil
		}

		// todo: 查询帖子是否存在 封面待确认？？
		post, err := svc.post.GetPostById(fmt.Sprint(replyInfo.PostId))
		if post == nil || err != nil {
			log.Log.Errorf("comment_trace: post not found, postId:%d", replyInfo.PostId)
			svc.engine.Rollback()
			return errdef.POST_NOT_EXISTS, nil
		}

		// 帖子状态 != 1 (1审核成功)
		if fmt.Sprint(post.Status) != consts.POST_AUDIT_SUCCESS {
			log.Log.Errorf("comment_trace: post status not audit success, postId:%d", replyInfo.PostId)
			svc.engine.Rollback()
			return errdef.VIDEO_NOT_EXISTS, nil
		}

		svc.comment.PostComment.UserId = userId
		svc.comment.PostComment.Content = params.Content
		svc.comment.PostComment.CreateAt = now
		svc.comment.PostComment.PostId = replyInfo.PostId
		svc.comment.PostComment.CommentLevel = consts.COMMENT_REPLY
		svc.comment.PostComment.Status = 1

		svc.comment.PostComment.ReplyCommentUserId = replyInfo.UserId
		svc.comment.PostComment.ReplyCommentId = replyInfo.Id
		svc.comment.PostComment.ParentCommentUserId = replyInfo.ParentCommentUserId
		svc.comment.PostComment.ParentCommentId = replyInfo.ParentCommentId

		isAt := 1
		// 1级评论 parentid为0 如果被回复的评论 parentid 为0，说明当前回复的是1级评论 否则 回复的为2级评论
		if replyInfo.ParentCommentId == 0 {
			svc.comment.PostComment.ParentCommentId = replyInfo.Id
			svc.comment.PostComment.ParentCommentUserId = replyInfo.UserId
			isAt = 0
		}

		if err := svc.comment.AddPostComment(); err != nil {
			log.Log.Errorf("comment_trace: add post reply err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_REPLY_FAIL, nil
		}

		// 更新视频总计（帖子评论总数）
		if err := svc.post.UpdatePostCommentNum(replyInfo.PostId, now, 1); err != nil {
			log.Log.Errorf("comment_trace: update post comment num err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_PUBLISH_FAIL, nil
		}

		svc.comment.ReceiveAt.TopicType = consts.TYPE_POST_COMMENT
		msgType = consts.POST_REPLY_MSG
		toUserId = replyInfo.UserId
		content = replyInfo.Content
		commentId = svc.comment.PostComment.Id
		composeId = replyInfo.PostId

		resp.Id = svc.comment.PostComment.Id
		resp.Avatar = user.Avatar
		resp.UserName = user.NickName
		resp.VideoId = params.ComposeId
		resp.CommentLevel = svc.comment.PostComment.CommentLevel
		resp.UserId = user.UserId
		resp.Content = svc.comment.PostComment.Content
		resp.CreateAt = svc.comment.PostComment.CreateAt
		resp.Status = svc.comment.PostComment.Status
		resp.IsAt = isAt
		resp.ReplyContent = replyInfo.Content
		resp.ReplyCommentUserId = replyInfo.UserId
		resp.ReplyCommentId = replyInfo.Id

		if uinfo := svc.user.FindUserByUserid(replyInfo.UserId); uinfo != nil {
			resp.ReplyCommentAvatar = uinfo.Avatar
			resp.ReplyCommentUserName = uinfo.NickName
		}

		atType = consts.TYPE_POST_COMMENT

	default:
		log.Log.Errorf("comment_trace: invalid commentType:%d", params.CommentType)
		return errdef.INVALID_PARAMS, nil
	}

	svc.comment.ReceiveAt.UserId = userId
	// 被@的用户
	svc.comment.ReceiveAt.ToUserId = toUserId
	svc.comment.ReceiveAt.ComposeId = svc.comment.VideoComment.Id
	svc.comment.ReceiveAt.UpdateAt = now
	svc.comment.ReceiveAt.CommentLevel = consts.COMMENT_REPLY
	// 回复 记录到 @
	if err := svc.comment.AddReceiveAt(); err != nil {
		log.Log.Errorf("comment_trace: add receive at err:%s", err)
		svc.engine.Rollback()
		return errdef.COMMENT_REPLY_FAIL, nil
	}

	// 添加@
	if len(params.AtInfo) > 0 {
		var atList []*models.ReceivedAt
		for _, val := range params.AtInfo {
			user := svc.user.FindUserByUserid(val)
			if user == nil {
				log.Log.Errorf("post_trace: at user not found, userId:%s", val)
				continue
			}

			at := &models.ReceivedAt{
				ToUserId:     val,
				UserId:       userId,
				ComposeId:    commentId,
				TopicType:    atType,
				CreateAt:     now,
				CommentLevel: consts.COMMENT_PUBLISH,
				UpdateAt:     now,
			}

			atList = append(atList, at)
		}

		affected, err := svc.post.AddReceiveAtList(atList)
		if err != nil || int(affected) != len(atList) {
			log.Log.Errorf("post_trace: add receive at list fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.COMMENT_REPLY_FAIL, nil
		}
	}


	svc.engine.Commit()
	// 视频/帖子 回复推送
	//event.PushEventMsg(config.Global.AmqpDsn, toUserId, user.NickName, cover, content, msgType)
	redismq.PushEventMsg(redismq.NewEvent(toUserId, fmt.Sprint(composeId), user.NickName, cover, content, msgType))
	return errdef.SUCCESS, resp
}

// 检查视频/帖子
func (svc *CommentModule) CheckComposeInfo(composeId string, commentType int) (int, int) {
	zanType := -1
	switch commentType {
	case consts.COMMENT_TYPE_VIDEO:
		video := svc.video.FindVideoById(composeId)
		// 视频不存在
		if video == nil {
			log.Log.Errorf("comment_trace: video not found or not pass, composeId:%s", composeId)
			return errdef.VIDEO_NOT_EXISTS, zanType
		}

		// 视频状态 != 1 (审核成功)
		if fmt.Sprint(video.Status) != consts.VIDEO_AUDIT_SUCCESS {
			log.Log.Errorf("comment_trace: video not audit success, videoId:%s", video.VideoId)
			return errdef.VIDEO_NOT_EXISTS, zanType
		}

		zanType = consts.TYPE_VIDEO_COMMENT

	case consts.COMMENT_TYPE_POST:
		post, err := svc.post.GetPostById(composeId)
		if post == nil || err != nil {
			log.Log.Errorf("comment_trace: post not found, composeId:%s", composeId)
			return errdef.POST_NOT_EXISTS, zanType
		}

		// 帖子状态 != 1 (审核成功)
		if fmt.Sprint(post.Status) != consts.POST_AUDIT_SUCCESS {
			log.Log.Errorf("comment_trace: post not audit success, postId:%s", post.Id)
			return errdef.POST_NOT_EXISTS, zanType
		}

		zanType = consts.TYPE_POST_COMMENT
	default:
		return errdef.INVALID_PARAMS, zanType
	}

	return errdef.SUCCESS, zanType
}

// 获取评论
func (svc *CommentModule) GetComments(userId, composeId, sortType string, commentType, page, size int) (int, []*mcomment.CommentList) {
	code, zanType := svc.CheckComposeInfo(composeId, commentType)
	if code != errdef.SUCCESS {
		return code, []*mcomment.CommentList{}
	}

	// 热门排序（按点赞数）
	if sortType == consts.SORT_HOT {
		log.Log.Debugf("comment_trace: get comments by hot")
		return errdef.SUCCESS, svc.GetCommentsByLiked(userId, composeId, zanType, page, size)
	}

	var list []*mcomment.CommentList
	offset := (page - 1) * size

	if commentType == consts.COMMENT_TYPE_POST {
		list = svc.GetPostCommentList(userId, composeId, offset, size, zanType)
	} else {
		list = svc.GetVideoCommentList(userId, composeId, offset, size, zanType)
	}


	return errdef.SUCCESS, list
}

// 获取帖子评论列表
func (svc *CommentModule) GetPostCommentList(userId, composeId string, offset, size, zanType int) []*mcomment.CommentList {
	// 获取评论
	comments := svc.comment.GetPostCommentList(composeId, offset, size)
	if len(comments) == 0 {
		log.Log.Errorf("comment_trace: no comments, composeId:%s", composeId)
		return []*mcomment.CommentList{}
	}

	list := make([]*mcomment.CommentList, len(comments))
	// contents 存储 评论id——>评论的内容
	contents := make(map[int64]string, 0)
	// userInfo 存储 用户id——>头像、昵称
	//userInfo := make(map[string]*tmpUser, 0)
	for index, item := range comments {
		comment := new(mcomment.CommentList)
		comment.Id = item.Id
		comment.Status = item.Status
		comment.PostId = item.PostId
		comment.UserId = item.UserId
		comment.CreateAt = item.CreateAt
		comment.IsTop = item.IsTop
		comment.Content = item.Content
		// 已被逻辑删除
		if comment.Status == 0 {
			comment.Content = "原内容已删除"
		}

		comment.CommentLevel = item.CommentLevel
		// 总回复数
		comment.ReplyNum = svc.comment.GetTotalReplyByPostComment(fmt.Sprint(comment.Id))
		// 评论点赞数
		comment.LikeNum = svc.like.GetLikeNumByType(item.Id, consts.TYPE_POST_COMMENT)

		// 如果总回复数 > 3 条
		if comment.ReplyNum > 3 {
			// 1 客户端展示查看更多
			comment.HasMore = 1
		}

		//comment.Avatar = item.Avatar
		//comment.UserName = item.UserName
		//user := new(tmpUser)
		//user.NickName = item.UserName
		//user.Avatar = item.Avatar
		//userInfo[item.UserId] = user
		// todo: 用户信息需使用最新数据
		user := svc.user.FindUserByUserid(comment.UserId)
		if user != nil {
			comment.Avatar = user.Avatar
			comment.UserName = user.NickName
		}

		contents[item.Id] = item.Content

		// 获取每个评论下的回复列表 (默认取三条)
		comment.ReplyList = svc.comment.GetPostReplyList(composeId, fmt.Sprint(item.Id), 0, 3)
		for _, reply := range comment.ReplyList {
			//user := new(tmpUser)
			//user.NickName = reply.UserName
			//user.Avatar = reply.Avatar
			//userInfo[reply.UserId] = user
			user = svc.user.FindUserByUserid(reply.UserId)
			if user != nil {
				reply.Avatar = user.Avatar
				reply.UserName = user.NickName
			}

			contents[reply.Id] = reply.Content
			// 评论点赞数
			reply.LikeNum = svc.like.GetLikeNumByType(reply.Id, zanType)

			// todo: 被回复的用户名、用户头像使用最新数据
			user = svc.user.FindUserByUserid(reply.ReplyCommentUserId)
			if user != nil {
				reply.ReplyCommentAvatar = user.Avatar
				reply.ReplyCommentUserName = user.NickName
			}

			// 如果回复的是1级评论 不展示@内容 否则展示   0 不是@消息 1是
			if reply.ParentCommentId != reply.ReplyCommentId || reply.ReplyCommentId != item.Id {
				reply.IsAt = 1
			}

			// 默认回复的是1级评论
			reply.ReplyContent = comment.Content
			// 被回复的内容
			content, ok := contents[reply.ReplyCommentId]
			if ok {
				reply.ReplyContent = content
			}

			if reply.Status == 0 {
				reply.ReplyContent = "原内容已删除"
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

			// 获取点赞的信息
			if likeInfo := svc.like.GetLikeInfo(userId, comment.Id, zanType); likeInfo != nil {
				comment.IsLike = likeInfo.Status
			}
		}

		if len(comment.ReplyList) == 0 {
			comment.ReplyList = []*mcomment.ReplyComment{}
		}

		list[index] = comment

	}

	return list
}

// 获取视频评论列表
func (svc *CommentModule) GetVideoCommentList(userId, composeId string, offset, size, zanType int) []*mcomment.CommentList {
	// 获取评论
	comments := svc.comment.GetVideoCommentList(composeId, offset, size)
	if len(comments) == 0 {
		log.Log.Errorf("comment_trace: no comments, composeId:%s", composeId)
		return []*mcomment.CommentList{}
	}

	list := make([]*mcomment.CommentList, len(comments))
	//type tmpUser struct {
	//	NickName   string
	//	Avatar     string
	//}
	// contents 存储 评论id——>评论的内容
	contents := make(map[int64]string, 0)
	// userInfo 存储 用户id——>头像、昵称
	//userInfo := make(map[string]*tmpUser, 0)
	for index, item := range comments {
		comment := new(mcomment.CommentList)
		comment.Id = item.Id
		comment.Status = item.Status
		comment.VideoId = item.VideoId
		comment.UserId = item.UserId
		comment.CreateAt = item.CreateAt
		comment.IsTop = item.IsTop
		comment.Content = item.Content
		// 已被逻辑删除
		if comment.Status == 0 {
			comment.Content = "原内容已删除"
		}

		comment.CommentLevel = item.CommentLevel
		// 总回复数
		comment.ReplyNum = svc.comment.GetTotalReplyByVideoComment(fmt.Sprint(comment.Id))
		// 评论点赞数
		comment.LikeNum = svc.like.GetLikeNumByType(item.Id, consts.TYPE_VIDEO_COMMENT)

		// 如果总回复数 > 3 条
		if comment.ReplyNum > 3 {
			// 1 客户端展示查看更多
			comment.HasMore = 1
		}

		//comment.Avatar = item.Avatar
		//comment.UserName = item.UserName
		//user := new(tmpUser)
		//user.NickName = item.UserName
		//user.Avatar = item.Avatar
		//userInfo[item.UserId] = user
		// todo: 用户信息需使用最新数据
		user := svc.user.FindUserByUserid(comment.UserId)
		if user != nil {
			comment.Avatar = user.Avatar
			comment.UserName = user.NickName
		}

		contents[item.Id] = item.Content

		// 获取每个评论下的回复列表 (默认取三条)
		comment.ReplyList = svc.comment.GetVideoReplyList(composeId, fmt.Sprint(item.Id), 0, 3)
		for _, reply := range comment.ReplyList {
			//user := new(tmpUser)
			//user.NickName = reply.UserName
			//user.Avatar = reply.Avatar
			//userInfo[reply.UserId] = user
			user = svc.user.FindUserByUserid(reply.UserId)
			if user != nil {
				reply.Avatar = user.Avatar
				reply.UserName = user.NickName
			}

			contents[reply.Id] = reply.Content
			// 评论点赞数
			reply.LikeNum = svc.like.GetLikeNumByType(reply.Id, zanType)

			// 被回复的用户名、用户头像
			//uinfo, ok := userInfo[reply.ReplyCommentUserId]
			//if ok {
			//	reply.ReplyCommentAvatar = uinfo.Avatar
			//	reply.ReplyCommentUserName = uinfo.NickName
			//}
			// todo: 被回复的用户名、用户头像使用最新数据
			user = svc.user.FindUserByUserid(reply.ReplyCommentUserId)
			if user != nil {
				reply.ReplyCommentAvatar = user.Avatar
				reply.ReplyCommentUserName = user.NickName
			}

			// 如果回复的是1级评论 不展示@内容 否则展示   0 不是@消息 1是
			if reply.ParentCommentId != reply.ReplyCommentId || reply.ReplyCommentId != item.Id {
				reply.IsAt = 1
			}

			// 默认回复的是1级评论
			reply.ReplyContent = comment.Content
			// 被回复的内容
			content, ok := contents[reply.ReplyCommentId]
			if ok {
				reply.ReplyContent = content
			}

			if reply.Status == 0 {
				reply.ReplyContent = "原内容已删除"
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

			// 获取点赞的信息
			if likeInfo := svc.like.GetLikeInfo(userId, comment.Id, zanType); likeInfo != nil {
				comment.IsLike = likeInfo.Status
			}
		}

		if len(comment.ReplyList) == 0 {
			comment.ReplyList = []*mcomment.ReplyComment{}
		}

		list[index] = comment

	}

	return list
}

// 根据评论点赞数排序 获取评论列表
func (svc *CommentModule) GetCommentsByLiked(userId, composeId string, zanType, page, size int) []*mcomment.CommentList {
	offset := (page - 1) * size
	var comments []*mcomment.CommentList
	switch zanType {
	case consts.TYPE_VIDEO_COMMENT:
		// 获取视频评论(按点赞数排序)
		comments = svc.comment.GetVideoCommentListByLike(composeId, zanType, offset, size)
	case consts.TYPE_POST_COMMENT:
		// 获取帖子评论(按点赞数排序)
		comments = svc.comment.GetPostCommentListByLike(composeId, zanType, offset, size)
	}

	if len(comments) == 0 {
		log.Log.Errorf("comment_trace: get comment fail, composeId:%s", composeId)
		return []*mcomment.CommentList{}
	}

	//type tmpUser struct {
	//	NickName   string
	//	Avatar     string
	//}
	// contents 存储 评论id——>评论的内容
	contents := make(map[int64]string, 0)
	// userInfo 存储 用户id——>头像、昵称
	//userInfo := make(map[string]*tmpUser, 0)
	for _, item := range comments {
		if item.Status == 0 {
			item.Content = "原内容已删除"
		}

		if item.VideoId > 0 {
			// 视频评论总回复数
			item.ReplyNum = svc.comment.GetTotalReplyByVideoComment(fmt.Sprint(item.Id))
			// 获取每个评论下的回复列表 (默认取三条)
			item.ReplyList = svc.comment.GetVideoReplyList(composeId, fmt.Sprint(item.Id), 0, 3)
		}

		if item.PostId > 0 {
			// 帖子评论总回复数
			item.ReplyNum = svc.comment.GetTotalReplyByPostComment(fmt.Sprint(item.Id))
			// 获取每个评论下的回复列表 (默认取三条)
			item.ReplyList = svc.comment.GetPostReplyList(composeId, fmt.Sprint(item.Id), 0, 3)
		}

		// 如果总回复数 > 3 条
		if item.ReplyNum > 3 {
			// 1 客户端展示查看更多
			item.HasMore = 1
		}


		//user := new(tmpUser)
		//user.NickName = item.UserName
		//user.Avatar = item.Avatar
		//userInfo[item.UserId] = user
		// todo: 评论用户的头像、昵称需使用最新的数据
		user := svc.user.FindUserByUserid(item.UserId)
		if user != nil {
			item.Avatar = user.Avatar
			item.UserName = user.NickName
		}

		contents[item.Id] = item.Content

		for _, reply := range item.ReplyList {
			//user := new(tmpUser)
			//user.NickName = reply.UserName
			//user.Avatar = reply.Avatar
			//userInfo[reply.UserId] = user
			// 已被逻辑删除
			if reply.Status == 0 {
				reply.Content = "原内容已删除"
			}

			contents[reply.Id] = reply.Content
			// 评论点赞数
			reply.LikeNum = svc.like.GetLikeNumByType(reply.Id, zanType)
			// 被回复的用户名、用户头像
			//uinfo, ok := userInfo[reply.ReplyCommentUserId]
			//if ok {
			//	reply.ReplyCommentAvatar = uinfo.Avatar
			//	reply.ReplyCommentUserName = uinfo.NickName
			//}
			// todo: 被回复的用户名、用户头像使用最新数据
			user := svc.user.FindUserByUserid(reply.UserId)
			if user != nil {
				reply.UserName = user.NickName
				reply.Avatar = user.Avatar
			}

			// todo: 被回复的用户名、用户头像使用最新数据
			user = svc.user.FindUserByUserid(reply.UserId)
			if user != nil {
				reply.ReplyCommentAvatar = user.Avatar
				reply.ReplyCommentUserName = user.NickName
			}

			// 如果回复的是1级评论 不展示@内容 否则展示   0 不是@消息 1是
			if reply.ParentCommentId != reply.ReplyCommentId || reply.ReplyCommentId != item.Id {
				reply.IsAt = 1
			}

			// 默认回复的是1级评论
			reply.ReplyContent = item.Content
			// 被回复的内容
			content, ok := contents[reply.ReplyCommentId]
			if ok {
				reply.ReplyContent = content
			}

			// 已逻辑删除
			if reply.Status == 0 {
				reply.ReplyContent = "原内容已删除"
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
				item.IsAttention = attention.Status
			}

			// 获取点赞的信息
			if likeInfo := svc.like.GetLikeInfo(userId, item.Id, zanType); likeInfo != nil {
				item.IsLike = likeInfo.Status
			}
		}

		if len(item.ReplyList) == 0 {
			item.ReplyList = []*mcomment.ReplyComment{}
		}

	}

	return comments
}

// 获取评论回复列表
func (svc *CommentModule) GetCommentReplyList(userId, composeId, commentId string, commentType, page, size int) (int, []*mcomment.ReplyComment) {
	code, zanType := svc.CheckComposeInfo(composeId, commentType)
	if code != errdef.SUCCESS {
		return code, []*mcomment.ReplyComment{}
	}

	var (
		replyList []*mcomment.ReplyComment
		content string
	)

	offset := (page - 1) * size
	switch commentType {
	case consts.COMMENT_TYPE_VIDEO:
		// 查询视频评论是否存在
		comment := svc.comment.GetVideoCommentById(commentId)
		if comment == nil {
			log.Log.Error("comment_trace: comment not found, commentId:%s", commentId)
			return errdef.COMMENT_NOT_FOUND, []*mcomment.ReplyComment{}
		}
		content = comment.Content

		replyList = svc.comment.GetVideoReplyList(composeId, commentId, offset, size)
	case consts.COMMENT_TYPE_POST:
		// 查询帖子评论是否存在
		comment := svc.comment.GetPostCommentById(commentId)
		if comment == nil {
			log.Log.Error("comment_trace: comment not found, commentId:%s", commentId)
			return errdef.COMMENT_NOT_FOUND, []*mcomment.ReplyComment{}
		}
		content = comment.Content

		replyList = svc.comment.GetPostReplyList(composeId, commentId, offset, size)
	}

	if len(replyList) == 0 {
		log.Log.Errorf("comment_trace: not found comment reply, commentId:%s", commentId)
		return errdef.SUCCESS, []*mcomment.ReplyComment{}
	}

	//type tmpUser struct {
	//	NickName   string
	//	Avatar     string
	//}
	// contents 存储 评论id——>评论的内容
	contents := make(map[int64]string, 0)
	// userInfo 存储 用户id——>头像、昵称
	//userInfo := make(map[string]*tmpUser, 0)
	for _, reply := range replyList {
		//user := new(tmpUser)
		//user.NickName = reply.UserName
		//user.Avatar = reply.Avatar
		//userInfo[reply.UserId] = user
		user := svc.user.FindUserByUserid(reply.UserId)
		if user != nil {
			reply.UserName = user.NickName
			reply.Avatar = user.Avatar
		}

		if reply.Status == 0 {
			reply.Content = "原内容已删除"
		}

		contents[reply.Id] = reply.Content
		// 评论点赞数
		reply.LikeNum = svc.like.GetLikeNumByType(reply.Id, zanType)
		// 被回复的用户名、用户头像
		//uinfo, ok := userInfo[reply.ReplyCommentUserId]
		//if ok {
		//	reply.ReplyCommentAvatar = uinfo.Avatar
		//	reply.ReplyCommentUserName = uinfo.NickName
		//}
		// todo: 用户信息需使用最新数据
		user = svc.user.FindUserByUserid(reply.ReplyCommentUserId)
		if user != nil {
			reply.ReplyCommentAvatar = user.Avatar
			reply.ReplyCommentUserName = user.NickName
		}

		// 如果回复的是1级评论 不展示@内容 否则展示   0 不是@消息 1是
		if reply.ParentCommentId != reply.ReplyCommentId || fmt.Sprint(reply.ReplyCommentId) != commentId {
			reply.IsAt = 1
		}

		// 默认回复的是1级评论
		reply.ReplyContent = content
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

			// 获取点赞的信息
			if likeInfo := svc.like.GetLikeInfo(userId, reply.Id, zanType); likeInfo != nil {
				reply.IsLike = likeInfo.Status
			}
		}
	}

	return errdef.SUCCESS, replyList
}

// 如果从消息页 点击某条@数据跳转到视频详情时 则 需要组装@消息的详情
func (svc *CommentModule) GetFirstComment(userId, commentId string) *mcomment.CommentList {
	var first *mcomment.CommentList
	// 如果从消息页 点击某条@数据跳转到视频详情时 则 需要组装点击的@消息的详情
	if commentId != "" {
		comment := svc.comment.GetVideoCommentById(commentId)
		if comment != nil {
			// todo:
			first = &mcomment.CommentList{
				Id: comment.Id,
				LikeNum:  svc.like.GetLikeNumByType(comment.Id, consts.TYPE_VIDEO_COMMENT),
				IsTop: comment.IsTop,
				CommentLevel: comment.CommentLevel,
				Content: comment.Content,
				CreateAt: comment.CreateAt,
				Status: comment.Status,
				VideoId: comment.VideoId,
				ReplyNum:  svc.comment.GetTotalReplyByVideoComment(fmt.Sprint(comment.Id)),
			}

			user := svc.user.FindUserByUserid(comment.UserId)
			if user != nil {
				first.Avatar = user.Avatar
				first.UserName = user.NickName
				first.UserId = user.UserId
				// 获取点赞的信息
				if likeInfo := svc.like.GetLikeInfo(userId, comment.Id, consts.TYPE_VIDEO_COMMENT); likeInfo != nil {
					first.IsLike = likeInfo.Status
				}
			}

			// contents 存储 评论id——>评论的内容
			content := make(map[int64]string, 0)
			content[comment.Id] = comment.Content
			first.ReplyList = svc.comment.GetVideoReplyList(fmt.Sprint(comment.VideoId), fmt.Sprint(comment.Id), 0, 3)
			for _, reply := range first.ReplyList {
				// 已被逻辑删除
				if reply.Status == 0 {
					reply.Content = "原内容已删除"
				}

				content[reply.Id] = reply.Content
				// 评论点赞数
				reply.LikeNum = svc.like.GetLikeNumByType(reply.Id, consts.TYPE_VIDEO_COMMENT)
				// todo: 被回复的用户名、用户头像使用最新数据
				user = svc.user.FindUserByUserid(reply.ReplyCommentUserId)
				if user != nil {
					reply.ReplyCommentAvatar = user.Avatar
					reply.ReplyCommentUserName = user.NickName
					reply.ReplyCommentUserId = user.UserId
					if user.UserId != "" {
						// 是否关注
						if attention := svc.attention.GetAttentionInfo(userId, reply.UserId); attention != nil {
							reply.IsAttention = attention.Status
						}

						// 获取点赞的信息
						if likeInfo := svc.like.GetLikeInfo(userId, reply.Id, consts.TYPE_VIDEO_COMMENT); likeInfo != nil {
							reply.IsLike = likeInfo.Status
						}
					}
				}

				// 被回复的内容
				content, ok := content[reply.ReplyCommentId]
				if ok {
					reply.ReplyContent = content
				}
			}
		}
	}

	return first
}

// 添加评论举报
func (svc *CommentModule) AddCommentReport(params *mcomment.CommentReportParam) int {
	// 查询评论是否存在
	comment := svc.comment.GetVideoCommentById(fmt.Sprint(params.CommentId))
	if comment == nil {
		log.Log.Error("comment_trace: comment not found, commentId:%s", fmt.Sprint(params.CommentId))
		return errdef.COMMENT_NOT_FOUND
	}

	svc.comment.Report.UserId = params.UserId
	svc.comment.Report.CommentId = params.CommentId
	svc.comment.Report.Reason = params.Reason
	svc.comment.Report.CreateAt = int(time.Now().Unix())
	svc.comment.Report.CommentType = params.CommentType
	if _, err := svc.comment.AddCommentReport(); err != nil {
		log.Log.Errorf("comment_trace: add comment report err:%s", err)
		return errdef.COMMENT_REPORT_FAIL
	}

	return errdef.SUCCESS
}
