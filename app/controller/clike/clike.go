package clike

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"strings"
	"time"
)

type LikeModule struct {
	context    *gin.Context
	engine     *xorm.Session
	user       *muser.UserModel
	like       *mlike.LikeModel
	video      *mvideo.VideoModel
	comment    *mcomment.CommentModel
	attention  *mattention.AttentionModel
}

func New(c *gin.Context) LikeModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return LikeModule{
		context: c,
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		comment: mcomment.NewCommentModel(socket),
		like: mlike.NewLikeModel(socket),
		attention: mattention.NewAttentionModel(socket),
		engine: socket,
	}
}

// 点赞视频
func (svc *LikeModule) GiveLikeForVideo(userId, toUserId string, videoId int64) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("like_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("like_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 查找视频是否存在
	if video := svc.video.FindVideoById(fmt.Sprint(videoId)); video == nil || fmt.Sprint(video.Status) != consts.VIDEO_AUDIT_SUCCESS  {
		log.Log.Errorf("like_trace: like video not found, videoId:%d", videoId)
		svc.engine.Rollback()
		return errdef.LIKE_VIDEO_NOT_EXISTS
	}

	// 获取点赞的视频信息
	info := svc.like.GetLikeInfo(userId, videoId, consts.TYPE_VIDEOS)
	// 是否已点赞
	// 已点赞
	if info != nil && info.Status == consts.ALREADY_GIVE_LIKE {
		log.Log.Errorf("like_trace: already give like, userId:%s, videoId:%d", userId, videoId)
		svc.engine.Rollback()
		return errdef.LIKE_ALREADY_EXISTS
	}

	now :=  int(time.Now().Unix())
	// 更新视频点赞总计 +1
	if err := svc.video.UpdateVideoLikeNum(videoId, now, consts.CONFIRM_OPERATE); err != nil {
		log.Log.Errorf("like_trace: update video like num err:%s", err)
		svc.engine.Rollback()
		return errdef.LIKE_VIDEO_FAIL
	}

	// 未点赞
	// 记录存在 且 状态为 未点赞 更新状态为 已点赞
	if info != nil && info.Status == consts.NOT_GIVE_LIKE {
		info.Status = consts.ALREADY_GIVE_LIKE
		info.CreateAt = now
		if err := svc.like.UpdateLikeStatus(); err != nil {
			log.Log.Errorf("like_trace: update like status err:%s", err)
			svc.engine.Rollback()
			return errdef.LIKE_VIDEO_FAIL
		}

	} else {
		// 添加点赞记录
		if err := svc.like.AddGiveLikeByType(userId, toUserId, videoId, consts.ALREADY_GIVE_LIKE, consts.TYPE_VIDEOS); err != nil {
			log.Log.Errorf("like_trace: add like video record err:%s", err)
			svc.engine.Rollback()
			return errdef.LIKE_VIDEO_FAIL
		}
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 取消点赞（视频）
func (svc *LikeModule) CancelLikeForVideo(userId string, videoId int64) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("like_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("like_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 查找视频是否存在
	if video := svc.video.FindVideoById(fmt.Sprint(videoId)); video == nil {
		log.Log.Errorf("like_trace: cancel like video not found, videoId:%d", videoId)
		svc.engine.Rollback()
		return errdef.LIKE_VIDEO_NOT_EXISTS
	}

	// 获取点赞的信息 判断是否已点赞 记录不存在 则 未点过赞
	info := svc.like.GetLikeInfo(userId, videoId, consts.TYPE_VIDEOS)
	if info == nil {
		log.Log.Errorf("like_trace: record not found, not give like, userId:%s, videoId:%d", userId, videoId)
		svc.engine.Rollback()
		return errdef.LIKE_RECORD_NOT_EXISTS
	}

	// 状态 ！= 已点赞 提示重复操作
	if info.Status != consts.ALREADY_GIVE_LIKE {
		log.Log.Errorf("like_trace: already cancel like, userId:%s, videoId:%d", userId, videoId)
		svc.engine.Rollback()
		return errdef.LIKE_REPEAT_CANCEL
	}

	now :=  int(time.Now().Unix())
	// 更新视频点赞总计 -1
	if err := svc.video.UpdateVideoLikeNum(videoId, now, consts.CANCEL_OPERATE); err != nil {
		log.Log.Errorf("like_trace: update video like num err:%s", err)
		svc.engine.Rollback()
		return errdef.LIKE_CANCEL_FAIL
	}

	info.Status = consts.NOT_GIVE_LIKE
	info.CreateAt = now
	// 更新状态 未点赞
	if err := svc.like.UpdateLikeStatus(); err != nil {
		log.Log.Errorf("like_trace: update like status err:%s", err)
		svc.engine.Rollback()
		return errdef.LIKE_CANCEL_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 获取用户点赞的视频列表
func (svc *LikeModule) GetUserLikeVideos(userId string, page, size int) []*mvideo.VideosInfoResp {
	offset := (page - 1) * size
	infos := svc.like.GetUserLikeVideos(userId, offset, size)
	if len(infos) == 0 {
		return []*mvideo.VideosInfoResp{}
	}

	// mp key videoId   value 用户视频点赞的时间
	mp := make(map[int64]int)
	// 当前页所有视频id
	videoIds := make([]string, len(infos))
	for index, like := range infos {
		mp[like.TypeId] = like.CreateAt
		videoIds[index] = fmt.Sprint(like.TypeId)
	}

	vids := strings.Join(videoIds, ",")
	// 获取点赞的视频列表信息
	videoList := svc.video.FindVideoListByIds(vids)
	if len(videoList) == 0 {
		log.Log.Errorf("like_trace: not found like video list info, len:%d, videoIds:%s", len(videoList), vids)
		return []*mvideo.VideosInfoResp{}
	}

	// 重新组装数据
	list := make([]*mvideo.VideosInfoResp, len(videoList))
	for index, video := range videoList {
		resp := new(mvideo.VideosInfoResp)
		resp.VideoId = video.VideoId
		resp.Title = video.Title
		resp.Describe = video.Describe
		resp.Cover = video.Cover
		resp.VideoAddr = video.VideoAddr
		resp.IsRecommend = video.IsRecommend
		resp.IsTop = video.IsTop
		resp.VideoDuration = video.VideoDuration
		resp.VideoWidth = video.VideoWidth
		resp.VideoHeight = video.VideoHeight
		resp.CreateAt = video.CreateAt
		resp.UserId = video.UserId
		// 获取用户信息
		if user := svc.user.FindUserByUserid(video.UserId); user != nil {
			resp.Avatar = user.Avatar
			resp.Nickname = user.NickName
		}

		// 是否关注
		attentionInfo := svc.attention.GetAttentionInfo(userId, video.UserId)
		if attentionInfo != nil {
			resp.IsAttention = attentionInfo.Status
		}

		collectAt, ok := mp[video.VideoId]
		if ok {
			// 用户给视频点赞的时间
			resp.OpTime = collectAt
		}

		list[index] = resp
	}

	return list
}

// 点赞评论
func (svc *LikeModule) GiveLikeForComment(userId, toUserId string, commentId int64) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("like_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("like_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 查找视频是否存在
	if video := svc.comment.GetVideoCommentById(fmt.Sprint(commentId)); video == nil {
		log.Log.Errorf("like_trace: like comment not found, commentId:%d", commentId)
		svc.engine.Rollback()
		return errdef.LIKE_COMMENT_NOT_EXISTS
	}

	// 获取点赞的评论信息
	info := svc.like.GetLikeInfo(userId, commentId, consts.TYPE_COMMENT)
	// 是否已点赞
	// 已点赞
	if info != nil && info.Status == consts.ALREADY_GIVE_LIKE {
		log.Log.Errorf("like_trace: already give like, userId:%s, commentId:%d", userId, commentId)
		svc.engine.Rollback()
		return errdef.LIKE_ALREADY_EXISTS
	}

	now :=  int(time.Now().Unix())
	// 未点赞
	// 记录存在 且 状态为 未点赞 更新状态为 已点赞
	if info != nil && info.Status == consts.NOT_GIVE_LIKE {
		info.Status = consts.ALREADY_GIVE_LIKE
		info.CreateAt = now
		if err := svc.like.UpdateLikeStatus(); err != nil {
			log.Log.Errorf("like_trace: update like comment status err:%s", err)
			svc.engine.Rollback()
			return errdef.LIKE_COMMENT_FAIL
		}

	} else {
		// 添加点赞记录
		if err := svc.like.AddGiveLikeByType(userId, toUserId, commentId, consts.ALREADY_GIVE_LIKE, consts.TYPE_COMMENT); err != nil {
			log.Log.Errorf("like_trace: add like comment record err:%s", err)
			svc.engine.Rollback()
			return errdef.LIKE_COMMENT_FAIL
		}
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 取消点赞（评论）
func (svc *LikeModule) CancelLikeForComment(userId string, commentId int64) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("like_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("like_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 查找评论是否存在
	if comment := svc.comment.GetVideoCommentById(fmt.Sprint(commentId)); comment == nil {
		log.Log.Errorf("like_trace: cancel like comment not found, commentId:%d", commentId)
		svc.engine.Rollback()
		return errdef.LIKE_COMMENT_NOT_EXISTS
	}

	// 获取点赞的信息 判断是否已点赞 记录不存在 则 未点过赞
	info := svc.like.GetLikeInfo(userId, commentId, consts.TYPE_COMMENT)
	if info == nil {
		log.Log.Errorf("like_trace: record not found, not give like, userId:%s, commentId:%d", userId, commentId)
		svc.engine.Rollback()
		return errdef.LIKE_RECORD_NOT_EXISTS
	}

	// 状态 ！= 已点赞 提示重复操作
	if info.Status != consts.ALREADY_GIVE_LIKE {
		log.Log.Errorf("like_trace: already cancel like, userId:%s, commentId:%d", userId, commentId)
		svc.engine.Rollback()
		return errdef.LIKE_REPEAT_CANCEL
	}

	now :=  int(time.Now().Unix())
	info.Status = consts.NOT_GIVE_LIKE
	info.CreateAt = now
	// 更新状态 未点赞
	if err := svc.like.UpdateLikeStatus(); err != nil {
		log.Log.Errorf("like_trace: update like status err:%s", err)
		svc.engine.Rollback()
		return errdef.LIKE_CANCEL_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}
