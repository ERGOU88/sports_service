package cpub

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/backend/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mbanner"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/minformation"
	"sports_service/server/models/mlabel"
	"sports_service/server/models/mlike"
	"sports_service/server/models/mnotify"
	"sports_service/server/models/mposting"
	"sports_service/server/models/msection"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	cloud "sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"strconv"
	"strings"
	"time"
)

type PubModule struct {
	context     *gin.Context
	engine      *xorm.Session
	video       *mvideo.VideoModel
	user        *muser.UserModel
	attention   *mattention.AttentionModel
	banner      *mbanner.BannerModel
	like        *mlike.LikeModel
	collect     *mcollect.CollectModel
	label       *mlabel.LabelModel
	notify      *mnotify.NotifyModel
	posting     *mposting.PostingModel
	section     *msection.SectionModel
	information *minformation.InformationModel
	community   *mcommunity.CommunityModel
}

func New(c *gin.Context) PubModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return PubModule{
		context:     c,
		video:       mvideo.NewVideoModel(socket),
		user:        muser.NewUserModel(socket),
		label:       mlabel.NewLabelModel(socket),
		posting:     mposting.NewPostingModel(socket),
		section:     msection.NewSectionModel(socket),
		information: minformation.NewInformationModel(socket),
		community:   mcommunity.NewCommunityModel(socket),
		engine:      socket,
	}
}

// 记录后台发布的视频信息
func (svc *PubModule) RecordPubVideoInfo(params *mvideo.VideoPublishParams) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("video_trace: session begin err:%s", err)
		return errdef.VIDEO_PUBLISH_FAIL
	}

	// 通过任务id获取用户id 是否为同一个用户
	uid, err := svc.video.GetUploadUserIdByTaskId(params.TaskId)
	if err != nil || strings.Compare(uid, params.UserId) != 0 {
		log.Log.Errorf("video_trace: user not match, cur userId:%s, uid:%s", params.UserId, uid)
		svc.engine.Rollback()
		return errdef.VIDEO_PUBLISH_FAIL
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(params.UserId); user == nil {
		log.Log.Errorf("video_trace: user not found, userId:%s", params.UserId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 用户发布视频
	if err := svc.UserPublishVideo(params.UserId, params); err != nil {
		log.Log.Errorf("video_trace: video publish failed, err:%s", err)
		svc.engine.Rollback()
		return errdef.VIDEO_PUBLISH_FAIL
	}

	info, _ := util.JsonFast.Marshal(params)

	// 记录到缓存 数据规则为 {videoId_info}
	if err := svc.video.RecordPublishInfo(params.UserId, svc.genVideoTag(svc.video.Videos.VideoId, string(info), svc.video.Videos.PubType), params.TaskId); err != nil {
		log.Log.Errorf("video_trace: record publish info by redis err:%s", err)
		svc.engine.Rollback()
		return errdef.VIDEO_PUBLISH_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 生成视频信息标签
func (svc *PubModule) genVideoTag(videoId int64, info string, pubType int) string {
	return fmt.Sprintf("%d__%s__%d",videoId, info, pubType)
}

// 用户发布视频
// 事务处理
// 标签记录到 视频标签表（多条记录 同一个videoId对应N个labelId 生成N条记录）
func (svc *PubModule) UserPublishVideo(userId string, params *mvideo.VideoPublishParams) error {
	var (
		subarea *models.VideoSubarea
		err error
	)

	// 视频所属分区
	if params.SubareaId != "" {
		subarea, err = svc.video.GetSubAreaById(params.SubareaId)
		if err != nil || subarea == nil {
			log.Log.Errorf("video_trace: get subarea by id fail, err:%s", err)
		} else {
			svc.video.Videos.Subarea = subarea.Id
		}
	}

	now := int(time.Now().Unix())
	// 视频所属专辑
	if params.AlbumId != "" {
		album, err := svc.video.GetVideoAlbumById(params.AlbumId)
		if err != nil || subarea == nil {
			log.Log.Errorf("video_trace: get video album by id fail, err:%s", err)
		} else {
			svc.video.Videos.Album = album.Id
		}
	}

	svc.video.Videos.UserId = userId
	svc.video.Videos.Cover = params.Cover
	svc.video.Videos.Title = params.Title
	svc.video.Videos.Describe = params.Describe
	svc.video.Videos.VideoAddr = params.VideoAddr
	svc.video.Videos.VideoDuration = params.VideoDuration
	svc.video.Videos.CreateAt = now
	svc.video.Videos.UpdateAt = now
	svc.video.Videos.UserType = consts.PUBLISH_VIDEO_BY_USER
	svc.video.Videos.VideoWidth = params.VideoWidth
	svc.video.Videos.VideoHeight = params.VideoHeight
	svc.video.Videos.Size = params.Size
	svc.video.Videos.SectionId = params.SectionId
	fileId, _ := strconv.Atoi(params.FileId)
	svc.video.Videos.FileId = int64(fileId)
	// 默认为首页发布
	svc.video.Videos.PubType = 1
	if params.PubType > 1 {
		svc.video.Videos.PubType = params.PubType
	}

	// 视频发布
	affected, err := svc.video.VideoPublish()
	if err != nil || affected != 1 {
		log.Log.Errorf("video_trace: publish video err:%s, affected:%d", err, affected)
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
		//info.Status = 1
		info.CreateAt = now
		labelInfos = append(labelInfos, info)
	}

	if len(labelInfos) > 0 {
		// 添加视频标签（多条）
		affected, err = svc.video.AddVideoLabels(labelInfos)
		if err != nil || int(affected) != len(labelInfos) {
			log.Log.Errorf("video_trace: add video labels err:%s", err)
			return errors.New("add video labels error")
		}
	}

	if params.PubType == 2 {
		// 同步到社区
		svc.posting.Posting.UserId = userId
		svc.posting.Posting.VideoId = svc.video.Videos.VideoId
		// 默认发布到综合
		svc.posting.Posting.SectionId = 1
		// 视频+文
		svc.posting.Posting.PostingType = consts.POST_TYPE_VIDEO
		// 社区发布
		svc.posting.Posting.ContentType = consts.COMMUNITY_PUB_POST

		svc.posting.Posting.CreateAt = now
		svc.posting.Posting.UpdateAt = now
		// 添加帖子
		if _, err := svc.posting.AddPost(); err != nil {
			log.Log.Errorf("video_trace: add posting fail, err:%s", err)
			return errors.New("add posting fail")
		}
	}

	svc.video.Statistic.VideoId = svc.video.Videos.VideoId
	svc.video.Statistic.CreateAt = now
	svc.video.Statistic.UpdateAt = now
	// 初始化视频统计数据
	if err := svc.video.AddVideoStatistic(); err != nil {
		log.Log.Errorf("video_trace: add video statistic err:%s", err)
		return err
	}

	return nil
}

// 获取上传签名
func (svc *PubModule) GetUploadSign(userId string, biteRate int64) (int, string, int64) {
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
	procedureName := svc.GetProcedureByBiteRate(biteRate)
	log.Log.Infof("procedure_trace: biteRate:%d, procedureName:%s", biteRate, procedureName)
	sign := client.GenerateSign(userId, procedureName, taskId)

	if err := svc.video.RecordUploadTaskId(userId, taskId); err != nil {
		log.Log.Errorf("video_trace: record upload taskid err:%s", err)
		return errdef.VIDEO_UPLOAD_GEN_SIGN_FAIL, "", 0
	}

	return errdef.SUCCESS, sign, taskId
}

// 通过码率 获取相应的任务模版
func (svc *PubModule) GetProcedureByBiteRate(biteRate int64) string {
	if biteRate >= 1800 {
		return consts.VOD_PROCEDURE_TRANSCODE_2
	}

	if biteRate > 1000 && biteRate < 1800 {
		return consts.VOD_PROCEDURE_TRANSCODE_1
	}


	return consts.VOD_PROCEDURE_NAME
}

// 发布帖子
func (svc *PubModule) PublishPosting(params *mposting.PostPublishParam) int {
	postType := svc.GetPostingType(params)
	if b := svc.VerifyContentLen(postType, params.Describe, params.Title); !b {
		log.Log.Error("post_trace: invalid content len")
		return errdef.POST_INVALID_CONTENT_LEN
	}

	b := util.IsSpace([]rune(params.Describe))
	if (!b || params.Describe == "") && len(params.ImagesAddr) == 0 {
		log.Log.Error("post_trace: describe and images empty")
		return errdef.POST_PARAMS_FAIL
	}

	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("post_trace: session begin err:%s", err)
		return errdef.ERROR
	}
	//
	if user := svc.user.FindUserByUserid(params.UserId); user == nil {
		log.Log.Errorf("post_trace: user not found, userId:%s", params.UserId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 获取板块信息
	section, err := svc.community.GetSectionInfo(fmt.Sprint(params.SectionId))
	if section == nil || err != nil {
		log.Log.Errorf("post_trace: section not found, sectionId:%d", params.SectionId)
		svc.engine.Rollback()
		return errdef.POST_SECTION_NOT_EXISTS
	}

	// 获取话题信息（多个）
	topics, err := svc.community.GetTopicByIds(params.TopicIds)
	if len(params.TopicIds) != len(topics) {
		log.Log.Errorf("post_trace: topic not found, topic_ids:%v, topics:%+v", params.TopicIds, topics)
		svc.engine.Rollback()
		return errdef.POST_TOPIC_NOT_EXISTS
	}

	now := int(time.Now().Unix())
	svc.posting.Posting.Title = params.Title
	svc.posting.Posting.Describe = params.Describe
	svc.posting.Posting.PostingType = postType
	// 社区发布
	svc.posting.Posting.ContentType = consts.COMMUNITY_PUB_POST
	svc.posting.Posting.CreateAt = now
	svc.posting.Posting.UpdateAt = now
	svc.posting.Posting.SectionId = section.Id
	svc.posting.Posting.UserId = params.UserId
	svc.posting.Posting.Status = 1
	if len(params.ImagesAddr) > 0 {
		bts, _ := util.JsonFast.Marshal(params.ImagesAddr)
		svc.posting.Posting.Content = string(bts)
	}

	if _, err := svc.posting.AddPost(); err != nil {
		log.Log.Errorf("post_trace: publish posting fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.POST_PUBLISH_FAIL
	}

	// 组装多条记录 写入帖子话题表
	topicInfos := make([]*models.PostingTopic, len(topics))
	for key, val := range topics {
		info := new(models.PostingTopic)
		info.TopicId = val.Id
		info.TopicName = val.TopicName
		info.CreateAt = now
		info.PostingId = svc.posting.Posting.Id
		info.Status = 1
		topicInfos[key] = info
	}

	if len(topicInfos) > 0 {
		// 添加帖子所属话题（多条）
		affected, err := svc.posting.AddPostingTopics(topicInfos)
		if err != nil || int(affected) != len(topicInfos) {
			svc.engine.Rollback()
			log.Log.Errorf("post_trace: add posting topics fail, err:%s", err)
			return errdef.POST_PUBLISH_FAIL
		}
	}

	svc.posting.Statistic.PostingId = svc.posting.Posting.Id
	svc.posting.Statistic.CreateAt = now
	svc.posting.Statistic.UpdateAt = now
	// 初始化帖子统计数据
	if err := svc.posting.AddPostStatistic(); err != nil {
		log.Log.Errorf("post_trace: add posting statistic err:%s", err)
		svc.engine.Rollback()
		return errdef.POST_PUBLISH_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 获取帖子类型
func (svc *PubModule) GetPostingType(params *mposting.PostPublishParam) (postType int) {
	if params.Describe != "" && len(params.ImagesAddr) > 0 || len(params.ImagesAddr) > 0 {
		// 图文
		postType = consts.POST_TYPE_IMAGE
		return
	}

	if params.Describe != "" {
		// 纯文本
		postType = consts.POST_TYPE_TEXT
	}

	return
}

// 验证正文长度
func (svc *PubModule) VerifyContentLen(postType int, content, title string) bool {
	if len(title) > 250 {
		return false
	}

	size := len(content)
	switch postType {
	// 纯文本、图文
	case 0, 1:
		if size > 10000 {
			return false
		}

	// 视频+文字
	case 2:
		if size > 250 {
			return false
		}
	}

	return true
}

// 发布资讯
func (svc *PubModule) PubInformation(param *models.Information) int {
	now := int(time.Now().Unix())
	param.CreateAt = now
	param.Status = 1
	if _, err := svc.information.AddInformation(param); err != nil {
		return errdef.ERROR
	}
	
	svc.information.Statistic.NewsId = param.Id
	svc.information.Statistic.CreateAt = now
	svc.information.Statistic.UpdateAt = now
	// 初始化资讯统计数据
	if _, err := svc.information.AddInformationStatistic(); err != nil {
		log.Log.Errorf("information_trace: add statistic id:%d, err:%s", param.Id, err)
		return errdef.ERROR
	}

	return errdef.SUCCESS
}

// 获取视频首页板块信息
func (svc *PubModule) GetHomepageSectionInfo(sectionType string) (int, []*models.RecommendInfoSection) {
	list, err := svc.section.GetRecommendSectionByType(sectionType)
	if err != nil {
		return errdef.ERROR, nil
	}
	
	if list == nil {
		return errdef.SUCCESS, []*models.RecommendInfoSection{}
	}
	
	return errdef.SUCCESS, list
}
