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
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"strings"
	"time"
)

type VideoModule struct {
	context      *gin.Context
	engine       *xorm.Session
	video        *mvideo.VideoModel
	user         *muser.UserModel
	attention    *mattention.AttentionModel
}

func New(c *gin.Context) VideoModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return VideoModule{
		context: c,
		video: mvideo.NewVideoModel(socket),
		user: muser.NewUserModel(socket),
		attention: mattention.NewAttentionModel(socket),
		engine: socket,
	}
}

// 用户发布视频
// 事务处理
// 数据记录到视频审核表 同时 标签记录到 视频标签表（多条记录 同一个videoId对应N个labelId 生成N条记录）
func (svc *VideoModule) UserPublishVideo(userId string, params *mvideo.VideoPublishParams) error {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", userId)
		return nil
	}

	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("video_trace: session begin err:%s", err)
		return err
	}

	// todo: 查询上传的标签是否存在
	now := time.Now().Unix()
	svc.video.Videos.UserId = userId
	svc.video.Videos.Cover = params.Cover
	svc.video.Videos.Title = params.Title
	svc.video.Videos.Describe = params.Describe
	svc.video.Videos.VideoAddr = params.VideoAddr
	svc.video.Videos.VideoDuration = params.VideoDuration
	svc.video.Videos.CreateAt = int(now)
	svc.video.Videos.UpdateAt = int(now)
	svc.video.Videos.UserType = consts.PUBLISH_VIDEO_BY_USER

	// 视频发布
	affected, err := svc.video.VideoPublish()
	if err != nil || affected != 1 {
		log.Log.Errorf("video_trace: publish video err:%s, affected:%d", err, affected)
		svc.engine.Rollback()
		return err
	}

	labelIds := strings.Split(params.VideoLabels, ",")
	// 组装多条记录 写入视频标签表
	labelInfos := make([]*models.VideoLabels, len(labelIds))
	for index, labelId := range labelIds {
		info := new(models.VideoLabels)
		info.VideoId = svc.video.Videos.VideoId
		info.LabelId = labelId
		// todo: 通过标签id获取名称
		info.LabelName = "todo:通过标签id获取名称"
		info.CreateAt = int(now)
		labelInfos[index] = info
	}

	// 添加视频标签（多条）
	affected, err = svc.video.AddVideoLabels(labelInfos)
	if err != nil || int(affected) != len(labelInfos) {
		svc.engine.Rollback()
		log.Log.Errorf("video_trace: add video labels err:%s", err)
		return errors.New("add video labels error")
	}

	svc.video.Statistic.VideoId = svc.video.Videos.VideoId
	svc.video.Statistic.CreateAt = int(now)
	svc.video.Statistic.UpdateAt = int(now)
	// 初始化视频统计数据
	if err := svc.video.AddVideoStatistic(); err != nil {
		log.Log.Errorf("video_trace: add video statistic err:%s", err)
		//svc.engine.Rollback()
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
		return nil
	}

	offset := (page - 1) * size
	records := svc.video.GetBrowseVideosRecord(userId, consts.TYPE_BROWSE_VIDEOS, offset, size)
	if len(records) == 0 {
		return nil
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
	videoList := svc.video.FindVideoListByIds(vids, offset, size)
	if len(videoList) == 0 {
		log.Log.Errorf("video_trace: not found browse video list info, len:%d, videoIds:%s", len(videoList), vids)
		return nil
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
		resp.IsAttention = attentionInfo.Status

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
func (svc *VideoModule) GetUserPublishList(userId, status, condition string, page, size int) []*mvideo.PublishVideosInfo {
	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", userId)
		return nil
	}

	offset := (page - 1) * size
	field := svc.GetConditionFieldByPublish(condition)
	// 获取用户发布的视频列表[通过审核状态和条件查询]
	return svc.video.GetUserPublishVideos(offset, size, userId, status, field)
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
		return errdef.VIDEO_DELETE_PUBLISH_FAIL
	}

	// 删除视频统计数据
	if err := svc.video.DelVideoStatistic(videoId); err != nil {
		log.Log.Errorf("video_trace: delete video statistic err:%s", err)
		svc.engine.Rollback()
		return errdef.VIDEO_DELETE_PUBLISH_FAIL
	}

	svc.engine.Commit()
	return errdef.SUCCESS
}
