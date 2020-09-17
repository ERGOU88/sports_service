package cvideo

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"strings"
	"fmt"
	"time"
	"errors"
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

// 用户发布的视频列表
// 事务处理
// 数据记录到视频审核表 同时 标签记录到 视频标签表（多条记录 同一个videoId对应N个labelId 生成N条记录）
func (svc *VideoModule) UserPublishVideo(userId string, params *mvideo.VideoPublishParams) error {
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
	if err := svc.video.VideoPublish(); err != nil {
		log.Log.Errorf("video_trace: publish video err:%s", err)
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
	affected, err := svc.video.AddVideoLabels(labelInfos)
	if err != nil || int(affected) != len(labelInfos) {
		svc.engine.Rollback()
		return errors.New("add video labels error")
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

// 用户浏览过的视频记录 todo:视频标签
func (svc *VideoModule) UserBrowseVideosRecord(userId string, page, size int) []*mvideo.VideosInfoResp {
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
