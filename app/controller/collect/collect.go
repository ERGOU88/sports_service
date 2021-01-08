package collect

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
  "sports_service/server/models/mlike"
  "sports_service/server/models/muser"
  "sports_service/server/models/mvideo"
  "sports_service/server/rabbitmq/event"
  "sports_service/server/util"
  "strings"
  "time"
  "sports_service/server/app/config"
)

type CollectModule struct {
	context     *gin.Context
	engine      *xorm.Session
	collect     *mcollect.CollectModel
	user        *muser.UserModel
	video       *mvideo.VideoModel
	attention   *mattention.AttentionModel
	like        *mlike.LikeModel
}

func New(c *gin.Context) CollectModule {
  socket := dao.Engine.NewSession()
	defer socket.Close()
	return CollectModule{
		context: c,
		collect: mcollect.NewCollectModel(socket),
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		attention: mattention.NewAttentionModel(socket),
		like: mlike.NewLikeModel(socket),
		engine: socket,
	}
}

// 添加收藏
func (svc *CollectModule) AddCollect(userId string, videoId int64) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("collect_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("collect_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 查找视频是否存在
	video := svc.video.FindVideoById(fmt.Sprint(videoId))
	if video == nil {
		log.Log.Errorf("collect_trace: video not found, videoId:%d", videoId)
		svc.engine.Rollback()
		return errdef.COLLECT_VIDEO_NOT_EXISTS
	}

	// 获取收藏的信息
	info := svc.collect.GetCollectInfo(userId, videoId, consts.TYPE_VIDEO)
	// 是否已收藏
	// 已收藏过
	if info != nil && info.Status == consts.ALREADY_COLLECT {
		log.Log.Errorf("collect_trace: already collect, userId:%s, videoId:%d", userId, videoId)
		svc.engine.Rollback()
		return errdef.COLLECT_ALREADY_EXISTS
	}

	now :=  int(time.Now().Unix())
	// 更新视频收藏总计 +1
	if err := svc.video.UpdateVideoCollectNum(videoId, now, consts.CONFIRM_OPERATE); err != nil {
		log.Log.Errorf("collect_trace: update video collect num err:%s", err)
		svc.engine.Rollback()
		return errdef.COLLECT_VIDEO_FAIL
	}

	// 未收藏
	// 记录存在 且 状态为未收藏 更新状态为收藏
	if info != nil && info.Status == consts.NO_COLLECT {
		info.Status = consts.ALREADY_COLLECT
		info.UpdateAt = int(time.Now().Unix())
		if err := svc.collect.UpdateCollectStatus(); err != nil {
			log.Log.Errorf("collect_trace: update collect status err:%s", err)
			svc.engine.Rollback()
			return errdef.COLLECT_VIDEO_FAIL
		}
	} else {
		// 添加收藏记录
		if err := svc.collect.AddCollectVideo(userId, video.UserId, videoId, consts.ALREADY_COLLECT, consts.TYPE_VIDEO); err != nil {
			log.Log.Errorf("collect_trace: add collect record err:%s", err)
			svc.engine.Rollback()
			return errdef.COLLECT_VIDEO_FAIL
		}
	}

	svc.engine.Commit()

  // 发送收藏视频推送
  event.PushEventMsg(config.Global.AmqpDsn, video.UserId, user.NickName, video.Cover, "", consts.COLLECT_VIDEO_MSG)

	return errdef.SUCCESS
}

// 取消收藏
func (svc *CollectModule) CancelCollect(userId string, videoId int64) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("collect_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("collect_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 查找视频是否存在
	if video := svc.video.FindVideoById(fmt.Sprint(videoId)); video == nil {
		log.Log.Errorf("collect_trace: video not found, videoId:%d", videoId)
		svc.engine.Rollback()
		return errdef.VIDEO_NOT_EXISTS
	}

	// 获取收藏的信息 判断是否已收藏 记录不存在 则 未收藏
	info := svc.collect.GetCollectInfo(userId, videoId, consts.TYPE_VIDEO)
	if info == nil {
		log.Log.Errorf("collect_trace: record not found, no collect, userId:%s, videoId:%d", userId, videoId)
		svc.engine.Rollback()
		return errdef.COLLECT_RECORD_NOT_EXISTS
	}

	// 状态 	!= 已收藏 提示重复操作
	if info.Status != consts.ALREADY_COLLECT {
		log.Log.Errorf("collect_trace: already cancel collect, userId:%s, videoId:%d", userId, videoId)
		svc.engine.Rollback()
		return errdef.COLLECT_REPEAT_CANCEL
	}

	now :=  int(time.Now().Unix())
	// 更新视频收藏总计 -1
	if err := svc.video.UpdateVideoCollectNum(videoId, now, consts.CANCEL_OPERATE); err != nil {
		log.Log.Errorf("collect_trace: update video collect num err:%s", err)
		svc.engine.Rollback()
		return errdef.COLLECT_CANCEL_FAIL
	}

	// 更新状态 未收藏
	info.Status = consts.NO_COLLECT
	info.UpdateAt = int(time.Now().Unix())
	if err := svc.collect.UpdateCollectStatus(); err != nil {
		log.Log.Errorf("collect_trace: update collect status err:%s", err)
		svc.engine.Rollback()
		return errdef.COLLECT_CANCEL_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 获取用户收藏的视频列表
func (svc *CollectModule) GetUserCollectVideos(userId string, page, size int) []*mvideo.VideosInfoResp {
  if userId == "" {
    log.Log.Errorf("collect_trace: userId not exists! userId:%s", userId)
    return []*mvideo.VideosInfoResp{}
  }

	offset := (page - 1) * size
	infos := svc.collect.GetCollectList(userId, offset, size)
	if len(infos) == 0 {
		return []*mvideo.VideosInfoResp{}
	}

	// mp key videoId  value 用户收藏视频的时间
	mp := make(map[int64]int)
	// 当前页所有视频id
	videoIds := make([]string, len(infos))
	for index, like := range infos {
		mp[like.ComposeId] = like.UpdateAt
		videoIds[index] = fmt.Sprint(like.ComposeId)
	}

	vids := strings.Join(videoIds, ",")
	// 获取收藏的视频列表信息
	videoList := svc.video.FindVideoListByIds(vids)
	if len(videoList) == 0 {
		log.Log.Errorf("collect_trace: not found collect video list info, len:%d, videoIds:%s", len(videoList), vids)
		return []*mvideo.VideosInfoResp{}
	}

	// 重新组装数据
	list := make([]*mvideo.VideosInfoResp, len(videoList))
	for index, video := range videoList {
		resp := new(mvideo.VideosInfoResp)
		resp.VideoId = video.VideoId
		resp.Title = util.TrimHtml(video.Title)
		resp.Describe = util.TrimHtml(video.Describe)
		resp.Cover = video.Cover
    resp.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
		resp.IsRecommend = video.IsRecommend
		resp.IsTop = video.IsTop
		resp.VideoDuration = video.VideoDuration
		resp.VideoWidth = video.VideoWidth
		resp.VideoHeight = video.VideoHeight
		resp.CreateAt = video.CreateAt
		resp.UserId = video.UserId
		if user := svc.user.FindUserByUserid(video.UserId); user != nil {
			resp.Avatar = user.Avatar
			resp.Nickname = user.NickName
      // 是否关注
      attentionInfo := svc.attention.GetAttentionInfo(userId, video.UserId)
      if attentionInfo != nil {
        resp.IsAttention = attentionInfo.Status
      }

      // 获取该视频的用户已播时长
      if record := svc.video.GetUserPlayDurationRecord(userId, fmt.Sprint(video.VideoId)); record != nil {
        resp.TimeElapsed = record.PlayDuration
      }
		}

		collectAt, ok := mp[video.VideoId]
		if ok {
			// 用户收藏视频的时间
			resp.OpTime = collectAt
		}

		list[index] = resp
	}

	return list
}

// 删除收藏记录
func (svc *CollectModule) DeleteCollectByIds(userId string, param *mcollect.DeleteCollectParam) int {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("collect_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	ids := strings.Join(param.ComposeIds, ",")
	if err := svc.collect.DeleteCollectByIds(userId, ids); err != nil {
		log.Log.Errorf("collect_trace: delete collect by ids err:%s", err)
		return errdef.COLLECT_DELETE_FAIL
	}

	return errdef.SUCCESS
}

