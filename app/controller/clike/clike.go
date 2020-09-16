package clike

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"fmt"
	"strings"
	"time"
)

type LikeModule struct {
	context    *gin.Context
	engine     *xorm.Session
	user       *muser.UserModel
	like       *mlike.LikeModel
	video      *mvideo.VideoModel
}

func New(c *gin.Context) LikeModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return LikeModule{
		context: c,
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		like: mlike.NewLikeModel(socket),
		engine: socket,
	}
}

// 点赞视频
func (svc *LikeModule) GiveLikeForVideo(userId string, videoId int64) int {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("like_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	// 查找视频是否存在
	if video := svc.video.FindVideoById(fmt.Sprint(videoId)); video == nil {
		log.Log.Errorf("like_trace: like video not found, videoId:%d", videoId)
		return errdef.LIKE_VIDEO_NOT_EXISTS
	}

	// 获取点赞的视频信息
	info := svc.like.GetLikeInfo(userId, videoId, consts.TYPE_VIDEO_LIKE)
	// 是否已点赞
	// 已点赞
	if info != nil && info.Status == consts.ALREADY_GIVE_LIKE {
		log.Log.Errorf("like_trace: already give like, userId:%s, videoId:%d", userId, videoId)
		return errdef.LIKE_ALREADY_EXISTS
	}

	// 未点赞
	// 记录存在 且 状态为 未点赞 更新状态为 已点赞
	if info != nil && info.Status == consts.NOT_GIVE_LIKE {
		info.Status = consts.ALREADY_GIVE_LIKE
		info.CreateAt = int(time.Now().Unix())
		if err := svc.like.UpdateLikeStatus(); err != nil {
			log.Log.Errorf("like_trace: update like status err:%s", err)
			return errdef.LIKE_VIDEO_FAIL
		}
	}

	// 添加点赞记录
	if err := svc.like.AddGiveLikeByType(userId, videoId, consts.ALREADY_GIVE_LIKE, consts.TYPE_VIDEO_LIKE); err != nil {
		log.Log.Errorf("like_trace: add like video record err:%s", err)
		return errdef.LIKE_VIDEO_FAIL
	}

	return errdef.SUCCESS
}

// 取消点赞（视频）
func (svc *LikeModule) CancelLikeForVideo(userId string, videoId int64) int {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("like_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	// 查找视频是否存在
	if video := svc.video.FindVideoById(fmt.Sprint(videoId)); video == nil {
		log.Log.Errorf("like_trace: cancel like video not found, videoId:%d", videoId)
		return errdef.LIKE_VIDEO_NOT_EXISTS
	}

	// 获取点赞的信息 判断是否已点赞 记录不存在 则 未点过赞
	info := svc.like.GetLikeInfo(userId, videoId, consts.TYPE_VIDEO_LIKE)
	if info == nil {
		log.Log.Errorf("like_trace: record not found, not give like, userId:%s, videoId:%d", userId, videoId)
		return errdef.LIKE_RECORD_NOT_EXISTS
	}

	// 状态 ！= 已点赞 提示重复操作
	if info.Status != consts.ALREADY_GIVE_LIKE {
		log.Log.Errorf("like_trace: already cancel like, userId:%s, videoId:%d", userId, videoId)
		return errdef.LIKE_REPEAT_CANCEL
	}

	info.Status = consts.NOT_GIVE_LIKE
	info.CreateAt = int(time.Now().Unix())
	// 更新状态 未点赞
	if err := svc.like.UpdateLikeStatus(); err != nil {
		log.Log.Errorf("like_trace: update like status err:%s", err)
		return errdef.LIKE_CANCEL_FAIL
	}

	return errdef.SUCCESS
}

// 获取用户点赞的视频列表
func (svc *LikeModule) GetUserLikeVideos(userId string, page, size int) []*models.Videos {
	videoIds := svc.like.GetUserLikeVideos(userId)
	if len(videoIds) == 0 {
		return nil
	}

	offset := (page - 1) * size
	vids := strings.Join(videoIds, ",")
	// 获取点赞的视频列表信息
	videoList := svc.video.FindVideoListByIds(vids, offset, size)
	if len(videoList) == 0 {
		log.Log.Errorf("like_trace: not found like video list info, len:%d, videoIds:%s", len(videoList), vids)
		return nil
	}

	return videoList
}
