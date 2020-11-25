package cvideo

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mbanner"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mlabel"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"sports_service/server/tools/tencentCloud/vod"
	"sports_service/server/util"
	"strconv"
	"strings"
	"time"
	cloud "sports_service/server/tools/tencentCloud"
)

type VideoModule struct {
	context      *gin.Context
	engine       *xorm.Session
	video        *mvideo.VideoModel
	user         *muser.UserModel
	attention    *mattention.AttentionModel
	banner       *mbanner.BannerModel
	like         *mlike.LikeModel
	collect      *mcollect.CollectModel
	label        *mlabel.LabelModel
}

func New(c *gin.Context) VideoModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return VideoModule{
		context: c,
		video: mvideo.NewVideoModel(socket),
		user: muser.NewUserModel(socket),
		attention: mattention.NewAttentionModel(socket),
		banner: mbanner.NewBannerMolde(socket),
		like: mlike.NewLikeModel(socket),
		collect: mcollect.NewCollectModel(socket),
		label: mlabel.NewLabelModel(socket),
		engine: socket,
	}
}

// 记录用户发布的视频信息(先放入缓存 接收到腾讯云回调事件 再写库）
func (svc *VideoModule) RecordPubVideoInfo(userId string, params *mvideo.VideoPublishParams) int {
	// 通过任务id获取用户id 是否为同一个用户
	uid, err := svc.video.GetUploadUserIdByTaskId(params.TaskId)
	if err != nil || strings.Compare(uid, userId) != 0 {
		log.Log.Errorf("video_trace: user not match, cur userId:%s, uid:%s", userId, uid)
		return errdef.VIDEO_PUBLISH_FAIL
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测视频描述
	isPass, err := client.TextModeration(params.Describe)
	if !isPass {
		log.Log.Errorf("video_trace: validate describe err: %s，pass: %v", err, isPass)
		return errdef.VIDEO_INVALID_DESCRIBE
	}

	// 检测视频标题
	isPass, err = client.TextModeration(params.Title)
	if !isPass {
		log.Log.Errorf("video_trace: validate title err: %s，pass: %v", err, isPass)
		return errdef.VIDEO_INVALID_TITLE
	}

	// 检测自定义标签
	if params.CustomLabels != "" {
    isPass, err = client.TextModeration(params.CustomLabels)
    if !isPass {
      log.Log.Errorf("video_trace: validate title err: %s，pass: %v", err, isPass)
      return errdef.VIDEO_INVALID_CUSTOM_LABEL
    }
  }

	info, _ := util.JsonFast.Marshal(params)
	// 先记录到缓存
	if err := svc.video.RecordPublishInfo(userId, string(info), params.TaskId); err != nil {
		log.Log.Errorf("video_trace: record publish info by redis err:%s", err)
		return errdef.VIDEO_PUBLISH_FAIL
	}

	return errdef.SUCCESS
}

// 用户发布视频
// 事务处理
// 数据记录到视频审核表 同时 标签记录到 视频标签表（多条记录 同一个videoId对应N个labelId 生成N条记录）
func (svc *VideoModule) UserPublishVideo(userId string, params *mvideo.VideoPublishParams) error {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("video_trace: session begin err:%s", err)
		return err
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errors.New("user not found")
	}

	now := time.Now().Unix()
	svc.video.Videos.UserId = userId
	svc.video.Videos.Cover = params.Cover
	svc.video.Videos.Title = util.TrimHtml(params.Title)
	svc.video.Videos.Describe = util.TrimHtml(params.Describe)
	svc.video.Videos.VideoAddr = params.VideoAddr
	svc.video.Videos.VideoDuration = params.VideoDuration
	svc.video.Videos.CreateAt = int(now)
	svc.video.Videos.UpdateAt = int(now)
	svc.video.Videos.UserType = consts.PUBLISH_VIDEO_BY_USER
	svc.video.Videos.VideoWidth = params.VideoWidth
	svc.video.Videos.VideoHeight = params.VideoHeight
	fileId, _ := strconv.Atoi(params.FileId)
	svc.video.Videos.FileId = int64(fileId)

	// 视频发布
	affected, err := svc.video.VideoPublish()
	if err != nil || affected != 1 {
		log.Log.Errorf("video_trace: publish video err:%s, affected:%d", err, affected)
		svc.engine.Rollback()
		return err
	}

	labelIds := strings.Split(params.VideoLabels, ",")
	// 组装多条记录 写入视频标签表
	labelInfos := make([]*models.VideoLabels, 0)
	for _, labelId := range labelIds {
		if svc.label.GetLabelInfoByMem(labelId) == nil {
			log.Log.Errorf("video_trace: label not found, labelId:%s", labelId)
			continue
		}

		info := new(models.VideoLabels)
		info.VideoId = svc.video.Videos.VideoId
		info.LabelId = labelId
		info.LabelName = svc.label.GetLabelNameByMem(labelId)
		info.CreateAt = int(now)
		labelInfos = append(labelInfos, info)
	}

  if len(labelInfos) > 0 {
    // 添加视频标签（多条）
    affected, err = svc.video.AddVideoLabels(labelInfos)
    if err != nil || int(affected) != len(labelInfos) {
      svc.engine.Rollback()
      log.Log.Errorf("video_trace: add video labels err:%s", err)
      return errors.New("add video labels error")
    }
  }

	svc.video.Statistic.VideoId = svc.video.Videos.VideoId
	svc.video.Statistic.CreateAt = int(now)
	svc.video.Statistic.UpdateAt = int(now)
	// 初始化视频统计数据
	if err := svc.video.AddVideoStatistic(); err != nil {
		log.Log.Errorf("video_trace: add video statistic err:%s", err)
		svc.engine.Rollback()
		return err
	}

	svc.engine.Commit()
	return nil
}

// 用户浏览过的视频记录 todo:视频标签 暂时只有视频 后续会有其他
func (svc *VideoModule) UserBrowseVideosRecord(userId string, page, size int) []*mvideo.VideosInfoResp {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", userId)
		return []*mvideo.VideosInfoResp{}
	}

	offset := (page - 1) * size
	records := svc.video.GetBrowseRecord(userId, consts.TYPE_BROWSE_VIDEOS, offset, size)
	if len(records) == 0 {
		return []*mvideo.VideosInfoResp{}
	}

	// mp key composeId   value 用户浏览的时间
	mp := make(map[int64]int)
	// 当前页所有视频id
	videoIds := make([]string, len(records))
	for index, info := range records {
		mp[info.ComposeId] = info.UpdateAt
		videoIds[index] = fmt.Sprint(info.ComposeId)
	}

	vids := strings.Join(videoIds, ",")
	// 获取浏览的视频列表信息
	videoList := svc.video.FindVideoListByIds(vids)
	if len(videoList) == 0 {
		log.Log.Errorf("video_trace: not found browse video list info, len:%d, videoIds:%s", len(videoList), vids)
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

		// 获取用户信息
		if user := svc.user.FindUserByUserid(video.UserId); user != nil {
			resp.Avatar = user.Avatar
			resp.Nickname = user.NickName
      // 是否关注
      attentionInfo := svc.attention.GetAttentionInfo(userId, video.UserId)
      if attentionInfo != nil {
        resp.IsAttention = attentionInfo.Status
      }

		}

		collectAt, ok := mp[video.VideoId]
		if ok {
			// 用户浏览视频的时间
			resp.OpTime = collectAt
		}

		list[index] = resp
	}

	return list
}

// 删除历史浏览记录
func (svc *VideoModule) DeleteHistoryByIds(userId string, param *mvideo.DeleteHistoryParam) int {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	ids := strings.Join(param.ComposeIds, ",")
	if err := svc.video.DeleteHistoryByIds(userId, ids); err != nil {
		log.Log.Errorf("video_trace: delete history by ids err:%s", err)
		return errdef.VIDEO_DELETE_HISTORY_FAIL
	}

	return errdef.SUCCESS
}

// 获取用户发布的列表（暂时只有视频）
func (svc *VideoModule) GetUserPublishList(userId, status, condition string, page, size int) []*mvideo.VideosInfo {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", userId)
		return []*mvideo.VideosInfo{}
	}

	offset := (page - 1) * size
	field := svc.GetConditionFieldByPublish(condition)

  // 获取用户发布的视频列表[通过审核状态和条件查询]
  list := svc.video.GetUserPublishVideos(offset, size, userId, status, field)
  if len(list) == 0 {
    log.Log.Errorf("video_trace: not publish video, userId:%s", userId)
    return []*mvideo.VideosInfo{}
  }

	for _, val := range list {
	  // todo: 已播时长（毫秒）
	  val.TimeElapsed = 10000
	  val.StatusCn = svc.GetVideoStatusCn(fmt.Sprint(val.Status))
	  val.VideoAddr = svc.video.AntiStealingLink(val.VideoAddr)
	  val.Describe = util.TrimHtml(val.Describe)
	  val.Title = util.TrimHtml(val.Title)
  }

  return list
}

// 获取视频状态（中文展示）
func (svc *VideoModule) GetVideoStatusCn(status string) string {
  switch status {
  case consts.VIDEO_UNDER_REVIEW:
    return "审核中"
  case consts.VIDEO_AUDIT_SUCCESS:
    return "已发布"
  case consts.VIDEO_AUDIT_FAILURE:
    return "未通过"
  }

  return "未知"
}

// -1 发布时间 0 播放数 1 弹幕数 2 评论数 3 点赞数 4 分享数
func (svc *VideoModule) GetConditionCn(condition string) string {
  switch condition {
  // 发布时间
  case consts.VIDEO_CONDITION_TIME:
    return "发布时间"
  // 播放数
  case consts.VIDEO_CONDITION_PLAY:
    return "播放数"
  // 弹幕数
  case consts.VIDEO_CONDITION_BARRAGE:
    return "弹幕数"
  // 评论数
  case consts.VIDEO_CONDITION_COMMENT:
    return "评论数"
  // 点赞数
  case consts.VIDEO_CONDITION_LIKE:
    return "点赞数"
  // 分享数
  case consts.VIDEO_CONDITION_SHARE:
    return "分享数"
  default:
    log.Log.Errorf("video_trace: unsupported condition, condition: %s", condition)
  }

  return "发布时间"
}

// 条件查询发布的内容
// -1 发布时间 0 播放数 1 弹幕数 2 评论数 3 点赞数 4 分享数
func (svc *VideoModule) GetConditionFieldByPublish(condition string) string {
	switch condition {
	// 发布时间
	case consts.VIDEO_CONDITION_TIME:
		return consts.CONDITION_FIELD_TIME
	// 播放数
	case consts.VIDEO_CONDITION_PLAY:
		return consts.CONDITION_FIELD_PLAY
	// 弹幕数
	case consts.VIDEO_CONDITION_BARRAGE:
		return consts.CONDITION_FIELD_BARRAGE
	// 评论数
	case consts.VIDEO_CONDITION_COMMENT:
		return consts.CONDITION_FIELD_COMMENT
	// 点赞数
	case consts.VIDEO_CONDITION_LIKE:
		return consts.CONDITION_FIELD_LIKE
	// 分享数
	case consts.VIDEO_CONDITION_SHARE:
		return consts.CONDITION_FIELD_SHARE
	default:
		log.Log.Errorf("video_trace: unsupported condition, condition: %s", condition)
	}

	return consts.CONDITION_FIELD_TIME
}

// 删除发布的视频
// 事务处理
// 产品逻辑：正式发布的视频 只能逻辑删除 且 其他流水数据 不可删除。 未审核或审核失败的 需删除 发布的视频、视频标签、视频总计数据
func (svc *VideoModule) DeletePublishVideo(userId, videoId string) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("video_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 查询视频信息
	video := svc.video.FindVideoById(videoId)
	if video == nil {
		log.Log.Errorf("video_trace: video not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.VIDEO_NOT_EXISTS
	}

	// 如果是已审核通过的视频 只能逻辑删除
	if fmt.Sprint(video.Status) == consts.VIDEO_AUDIT_SUCCESS {
		// 状态3 为逻辑删除
		video.Status = 3
		if err := svc.video.UpdateVideoStatus(userId, videoId); err != nil {
			log.Log.Errorf("video_trace: update video status err:%s", err)
			svc.engine.Rollback()
			return errdef.VIDEO_DELETE_PUBLISH_FAIL
		}

		svc.engine.Commit()
		return errdef.SUCCESS
	}

	// 视频为未审核/审核失败 物理删除发布的视频、视频标签、视频总计
	if err := svc.video.DelPublishById(userId, videoId); err != nil {
		log.Log.Errorf("video_trace: delete publish by id err:%s", err)
		svc.engine.Rollback()
		return errdef.VIDEO_DELETE_PUBLISH_FAIL
	}

	// 删除视频标签
	if err := svc.video.DelVideoLabels(videoId); err != nil {
		log.Log.Errorf("video_trace: delete video labels err:%s", err)
		svc.engine.Rollback()
		return errdef.VIDEO_DELETE_LABEL_FAIL
	}

	// 删除视频统计数据
	if err := svc.video.DelVideoStatistic(videoId); err != nil {
		log.Log.Errorf("video_trace: delete video statistic err:%s", err)
		svc.engine.Rollback()
		return errdef.VIDEO_DELETE_STATISTIC_FAIL
	}

	svc.engine.Commit()
	return errdef.SUCCESS
}

// 获取推荐的视频列表
func (svc *VideoModule) GetRecommendVideos(userId string, page, size int) []*mvideo.RecommendVideoInfo {
	offset := (page - 1) * size
	list := svc.video.GetRecommendVideoList(offset, size)
	if len(list) == 0 {
		return []*mvideo.RecommendVideoInfo{}
	}

	// 重新组装数据
	for _, video := range list {
		// 查询用户信息
		userInfo := svc.user.FindUserByUserid(video.UserId)
		if userInfo == nil {
			log.Log.Errorf("video_trace: user not found, uid:%s", video.UserId)
			continue
		}

    video.Describe = util.TrimHtml(video.Describe)
    video.Title = util.TrimHtml(video.Title)

		video.Avatar = userInfo.Avatar
		video.Nickname = userInfo.NickName
		video.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)

		// 用户未登录
		if userId == "" {
			log.Log.Error("video_trace: no login")
			continue
		}

		// 获取点赞的信息
		if likeInfo := svc.like.GetLikeInfo(userId, video.VideoId, consts.TYPE_VIDEO); likeInfo != nil {
			video.IsLike = likeInfo.Status
		}

		// 获取收藏的信息
		if collectInfo := svc.collect.GetCollectInfo(userId, video.VideoId, consts.TYPE_VIDEO); collectInfo != nil {
			video.IsCollect = collectInfo.Status
		}

	}

	return list
}

// 获取app首页推荐的banner 默认取10条
func (svc *VideoModule) GetRecommendBanners() []*models.Banner {
	banners := svc.banner.GetRecommendBanners(int32(consts.HOMEPAGE_BANNERS), time.Now().Unix(), 0, 10)
	if len(banners) == 0 {
	  return []*models.Banner{}
  }

  return banners
}

// 获取关注的用户发布的视频列表
func (svc *VideoModule) GetAttentionVideos(userId string, page, size int) []*mvideo.VideoDetailInfo {
	// 用户未登录
	if userId == "" {
		log.Log.Error("video_trace: no login")
		return []*mvideo.VideoDetailInfo{}
	}

	userIds := svc.attention.GetAttentionList(userId)
	if len(userIds) == 0 {
		log.Log.Errorf("video_trace: not following any users")
		return []*mvideo.VideoDetailInfo{}
	}

	offset := (page - 1) * size
	uids := strings.Join(userIds, ",")
	list := svc.video.GetAttentionVideos(uids, offset, size)
	if len(list) == 0 {
		return []*mvideo.VideoDetailInfo{}
	}

	// 重新组装数据
	for _, video := range list {
		// 获取视频标签信息
    video.Labels = svc.video.GetVideoLabels(fmt.Sprint(video.VideoId))
    if video.Labels == nil {
      video.Labels = []*models.VideoLabels{}
    }

		// 查询用户信息
		userInfo := svc.user.FindUserByUserid(video.UserId)
		if userInfo == nil {
			log.Log.Errorf("video_trace: user not found, uid:%s", video.UserId)
			continue
		}

    video.Describe = util.TrimHtml(video.Describe)
    video.Title = util.TrimHtml(video.Title)

		video.Avatar = userInfo.Avatar
		video.Nickname = userInfo.NickName
    video.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)

		if userId == "" {
			log.Log.Error("video_trace: user no login")
			continue
		}

		video.IsAttention = consts.ALREADY_ATTENTION

		// 获取点赞的信息
		if likeInfo := svc.like.GetLikeInfo(userId, video.VideoId, consts.TYPE_VIDEO); likeInfo != nil {
			video.IsLike = likeInfo.Status
		}

		// 获取收藏的信息
		if collectInfo := svc.collect.GetCollectInfo(userId, video.VideoId, consts.TYPE_VIDEO); collectInfo != nil {
			video.IsCollect = collectInfo.Status
		}

	}

	return list
}

// 获取视频详情页数据
func (svc *VideoModule) GetVideoDetail(userId, videoId string) *mvideo.VideoDetailInfo {
	if videoId == "" {
		log.Log.Error("video_trace: videoId can't empty")
		return nil
	}

	video := svc.video.FindVideoById(videoId)
	if video == nil {
		log.Log.Error("video_trace: video not found, videoId:%s", videoId)
		return nil
	}

	resp := new(mvideo.VideoDetailInfo)
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
	resp.Labels = svc.video.GetVideoLabels(fmt.Sprint(video.VideoId))
  if resp.Labels == nil {
    resp.Labels = []*models.VideoLabels{}
  }

  // 获取转码后的视频数据
  if err := util.JsonFast.UnmarshalFromString(video.PlayInfo, &resp.PlayInfo); err != nil {
    log.Log.Errorf("video_trace: jsonFast unmarshal err:%s", err)
    resp.PlayInfo = []*mvideo.PlayInfo{}
  }

  if len(resp.PlayInfo) > 0 {
    for _, v := range resp.PlayInfo {
      // 添加防盗链
      v.Url = svc.video.AntiStealingLink(v.Url)
    }
  }

	// 获取视频相关统计数据
	info := svc.video.GetVideoStatistic(fmt.Sprint(video.VideoId))
	resp.BrowseNum = info.BrowseNum
	resp.CommentNum = info.CommentNum
	resp.FabulousNum = info.FabulousNum
	resp.ShareNum = info.ShareNum
	resp.BarrageNum = info.BarrageNum
	resp.CollectNum = info.CollectNum
	// 粉丝数
	resp.FansNum = svc.attention.GetTotalFans(fmt.Sprint(video.UserId))
  now := int(time.Now().Unix())
  // 增加视频浏览总数
  if err := svc.video.UpdateVideoBrowseNum(video.VideoId, now, 1); err != nil {
    log.Log.Errorf("video_trace: update video browse num err:%s", err)
  }

  if user := svc.user.FindUserByUserid(video.UserId); user != nil {
    resp.Avatar = user.Avatar
    resp.Nickname = user.NickName
  }

	if userId == "" {
		log.Log.Error("video_trace: user no login")
		return resp
	}

  // 获取用户信息
  if user := svc.user.FindUserByUserid(userId); user != nil {

    // 用户是否浏览过
    browse := svc.video.GetUserBrowseVideo(userId, consts.TYPE_VIDEO, video.VideoId)
    if browse != nil {
      svc.video.Browse.CreateAt = now
      svc.video.Browse.UpdateAt = now
      // 已有浏览记录 更新用户浏览的时间
      if err := svc.video.UpdateUserBrowseVideo(userId, consts.TYPE_VIDEO, video.VideoId); err != nil {
        log.Log.Errorf("video_trace: update user browse video err:%s", err)
      }
    } else {
      svc.video.Browse.CreateAt = now
      svc.video.Browse.UpdateAt = now
      svc.video.Browse.UserId = userId
      svc.video.Browse.ComposeId = video.VideoId
      svc.video.Browse.ComposeType = consts.TYPE_VIDEO
      // 添加用户浏览的视频记录
      if err := svc.video.RecordUserBrowseVideo(); err != nil {
        log.Log.Errorf("video_trace: record user browse video err:%s", err)
      }
    }
  }

	// 是否关注
	if attentionInfo := svc.attention.GetAttentionInfo(userId, video.UserId); attentionInfo != nil {
		resp.IsAttention = attentionInfo.Status
	}

	// 获取点赞的信息
	if likeInfo := svc.like.GetLikeInfo(userId, video.VideoId, consts.TYPE_VIDEO); likeInfo != nil {
		resp.IsLike = likeInfo.Status
	}

	// 获取收藏的信息
	if collectInfo := svc.collect.GetCollectInfo(userId, video.VideoId, consts.TYPE_VIDEO); collectInfo != nil {
		resp.IsCollect = collectInfo.Status
	}

	return resp

}

// 获取详情页推荐视频（根据同标签推荐）
func (svc *VideoModule) GetDetailRecommend(userId, videoId string, page, size int) []*mvideo.VideoDetailInfo {
	if videoId == "" {
		log.Log.Error("video_trace: videoId can't empty")
		return []*mvideo.VideoDetailInfo{}
	}

	video := svc.video.FindVideoById(videoId)
	if video == nil {
		log.Log.Error("video_trace: video not found, videoId:%s", videoId)
		return []*mvideo.VideoDetailInfo{}
	}

	// 获取视频所有标签
	labels := svc.video.GetVideoLabels(fmt.Sprint(video.VideoId))
	ids := make([]string, len(labels))
	for index, label := range labels {
		ids[index] = label.LabelId
	}

	labelIds := strings.Join(ids, ",")
	offset := (page - 1) * size
	// 通过标签列表 获取拥有该标签的视频们
	videoIds := svc.video.FindVideoIdsByLabelIds(labelIds, offset, size)
	if len(videoIds) == 0 {
		log.Log.Errorf("search_trace: not found videos by label ids, labelIds:%s", labelIds)
		return []*mvideo.VideoDetailInfo{}
	}

	vids := strings.Join(videoIds, ",")
	videos := svc.video.FindVideoListByIds(vids)
	if len(videos) == 0 {
		log.Log.Errorf("search_trace: not found videos, vids:%s", vids)
		return []*mvideo.VideoDetailInfo{}
	}

	// 重新组装返回数据
	res := make([]*mvideo.VideoDetailInfo, len(videos))
	for index, video := range videos {
		resp := new(mvideo.VideoDetailInfo)
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
		resp.Labels = svc.video.GetVideoLabels(fmt.Sprint(video.VideoId))
    if resp.Labels == nil {
      resp.Labels = []*models.VideoLabels{}
    }
		// 获取用户信息
		if user := svc.user.FindUserByUserid(video.UserId); user != nil {
			resp.Avatar = user.Avatar
			resp.Nickname = user.NickName
		}

		// 获取视频相关统计数据
		info := svc.video.GetVideoStatistic(fmt.Sprint(video.VideoId))
		resp.BrowseNum = info.BrowseNum
		resp.CommentNum = info.CommentNum
		resp.FabulousNum = info.FabulousNum
		resp.ShareNum = info.ShareNum
		resp.BarrageNum = info.BarrageNum
		// 粉丝数
		resp.FansNum = svc.attention.GetTotalFans(fmt.Sprint(video.UserId))

		if userId != "" {
			log.Log.Error("video_trace: user no login")
      // 是否关注
      if attentionInfo := svc.attention.GetAttentionInfo(userId, video.UserId); attentionInfo != nil {
        resp.IsAttention = attentionInfo.Status
      }

      // 获取点赞的信息
      if likeInfo := svc.like.GetLikeInfo(userId, video.VideoId, consts.TYPE_VIDEO); likeInfo != nil {
        resp.IsLike = likeInfo.Status
      }

      // 获取收藏的信息
      if collectInfo := svc.collect.GetCollectInfo(userId, video.VideoId, consts.TYPE_VIDEO); collectInfo != nil {
        resp.IsCollect = collectInfo.Status
      }
		}

		resp.PlayInfo = []*mvideo.PlayInfo{}

		res[index] = resp
	}

	return res
}

// 获取上传签名
func (svc *VideoModule) GetUploadSign(userId string) (int, string, int64) {
	// 用户未登录
	if userId == "" {
		log.Log.Error("video_trace: no login")
		return errdef.USER_NO_LOGIN, "", 0
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS, "", 0
	}

	client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
	taskId := util.GetXID()
	sign := client.GenerateSign(userId, consts.VOD_PROCEDURE_NAME, taskId)

	if err := svc.video.RecordUploadTaskId(userId, taskId); err != nil {
		log.Log.Errorf("video_trace: record upload taskid err:%s", err)
		return errdef.VIDEO_UPLOAD_GEN_SIGN_FAIL, "", 0
	}

	return errdef.SUCCESS, sign, taskId
}

// 事件回调(被动)
func (svc *VideoModule) EventCallback(params *vod.EventNotify) int {
	switch params.EventType {
	// 上传事件
	case consts.EVENT_TYPE_UPLOAD:
		context := new(cloud.SourceContext)
		if err := util.JsonFast.Unmarshal([]byte(params.FileUploadEvent.MediaBasicInfo.SourceInfo.SourceContext), context); err != nil {
			log.Log.Errorf("video_trace: jsonfast unmarshal sourceContext err:%s", err)
			return errdef.INVALID_PARAMS
		}

		// 通过任务id 获取 用户id
		userId, err := svc.video.GetUploadUserIdByTaskId(context.TaskId)
		if err != nil || userId == "" {
			log.Log.Errorf("video_trace: invalid taskId, taskId:%d", context.TaskId)
			return errdef.ERROR
		}

		// 查询用户是否存在
		if user := svc.user.FindUserByUserid(userId); user == nil {
			log.Log.Errorf("video_trace: user not found, userId:%s", userId)
			return errdef.USER_NOT_EXISTS
		}


	}

	return errdef.SUCCESS
}

// 检测自定义标签
func (svc *VideoModule) CheckCustomLabel(userId string, params *mvideo.CustomLabelParams) int {
  // 用户未登录
  if userId == "" {
    log.Log.Error("video_trace: no login")
    return errdef.USER_NO_LOGIN
  }

  // 查询用户是否存在
  if user := svc.user.FindUserByUserid(userId); user == nil {
    log.Log.Errorf("video_trace: user not found, userId:%s", userId)
    return errdef.USER_NOT_EXISTS
  }

  client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
  // 检测视频描述
  isPass, err := client.TextModeration(params.CustomLabel)
  if !isPass {
    log.Log.Errorf("video_trace: validate custom label err: %s，pass: %v", err, isPass)
    return errdef.VIDEO_INVALID_CUSTOM_LABEL
  }

  return errdef.SUCCESS
}

// 获取视频标签列表
func (svc *VideoModule) GetVideoLabelList() []*mlabel.VideoLabel {
  list := svc.label.GetVideoLabelList()
  if len(list) == 0 {
    return []*mlabel.VideoLabel{}
  }

  return list
}

// 添加视频举报
func (svc *VideoModule) AddVideoReport(params *mvideo.VideoReportParam) int {
  video := svc.video.FindVideoById(fmt.Sprint(params.VideoId))
  if video == nil {
    log.Log.Error("video_trace: video not found, videoId:%s", params.VideoId)
    return errdef.VIDEO_NOT_EXISTS
  }

  svc.video.Report.UserId = params.UserId
  svc.video.Report.VideoId = params.VideoId
  if _, err := svc.video.AddVideoReport(); err != nil {
    log.Log.Errorf("video_trace: add video report err:%s", err)
    return errdef.VIDEO_REPORT_FAIL
  }

  return errdef.SUCCESS
}

