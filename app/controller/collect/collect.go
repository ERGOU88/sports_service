package collect

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"fmt"
	"strings"
	"time"
)

type CollectModule struct {
	context     *gin.Context
	engine      *xorm.Session
	collect     *mcollect.CollectModel
	user        *muser.UserModel
	video       *mvideo.VideoModel
}

func New(c *gin.Context) CollectModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return CollectModule{
		context: c,
		collect: mcollect.NewCollectModel(socket),
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		engine: socket,
	}
}

// 添加收藏
func (svc *CollectModule) AddCollect(userId string, videoId int64) int {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("collect_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	// 查找视频是否存在
	if video := svc.video.FindVideoById(fmt.Sprint(videoId)); video == nil {
		log.Log.Errorf("collect_trace: video not found, videoId:%d", videoId)
		return errdef.COLLECT_VIDEO_NOT_EXISTS
	}

	// 获取收藏的信息
	info := svc.collect.GetCollectInfo(userId, videoId)
	// 是否已收藏
	// 已收藏过
	if info != nil && info.Status == consts.ALREADY_COLLECT {
		log.Log.Errorf("collect_trace: already collect, userId:%s, videoId:%d", userId, videoId)
		return errdef.COLLECT_ALREADY_EXISTS
	}

	// 未收藏
	// 记录存在 且 状态为未收藏 更新状态为收藏
	if info != nil && info.Status == consts.NO_COLLECT {
		info.Status = consts.ALREADY_COLLECT
		info.UpdateAt = int(time.Now().Unix())
		if err := svc.collect.UpdateCollectStatus(); err != nil {
			log.Log.Errorf("collect_trace: update collect status err:%s", err)
			return errdef.COLLECT_VIDEO_FAIL
		}
	}

	// 添加收藏记录
	if err := svc.collect.AddCollectVideo(userId, videoId, consts.ALREADY_COLLECT); err != nil {
		log.Log.Errorf("collect_trace: add collect record err:%s", err)
		return errdef.COLLECT_VIDEO_FAIL
	}

	return errdef.SUCCESS
}

// 取消收藏
func (svc *CollectModule) CancelCollect(userId string, videoId int64) int {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("collect_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	// 查找视频是否存在
	if video := svc.video.FindVideoById(fmt.Sprint(videoId)); video == nil {
		log.Log.Errorf("collect_trace: video not found, videoId:%d", videoId)
		return errdef.VIDEO_NOT_EXISTS
	}

	// 获取收藏的信息 判断是否已收藏 记录不存在 则 未收藏
	info := svc.collect.GetCollectInfo(userId, videoId)
	if info == nil {
		log.Log.Errorf("collect_trace: record not found, no collect, userId:%s, videoId:%d", userId, videoId)
		return errdef.COLLECT_RECORD_NOT_EXISTS
	}

	// 状态 	!= 已收藏 提示重复操作
	if info.Status != consts.ALREADY_COLLECT {
		log.Log.Errorf("collect_trace: already cancel collect, userId:%s, videoId:%d", userId, videoId)
		return errdef.COLLECT_REPEAT_CANCEL
	}

	// 更新状态 未收藏
	info.Status = consts.NO_COLLECT
	info.UpdateAt = int(time.Now().Unix())
	if err := svc.collect.UpdateCollectStatus(); err != nil {
		log.Log.Errorf("collect_trace: update collect status err:%s", err)
		return errdef.COLLECT_CANCEL_FAIL
	}

	return errdef.SUCCESS
}

// 获取用户收藏的视频列表
func (svc *CollectModule) GetUserCollectVideos(userId string, page, size int) []*models.Videos {
	videoIds := svc.collect.GetCollectVideos(userId)
	if len(videoIds) == 0 {
		return nil
	}

	offset := (page - 1) * size
	vids := strings.Join(videoIds, ",")
	// 获取收藏的视频列表
	videoList := svc.video.FindVideoListByIds(vids, offset, size)
	if len(videoList) == 0 {
		log.Log.Errorf("collect_trace: not found video list info, len:%d, videoIds:%s", len(videoList), vids)
		return nil
	}

	return videoList
}
