package cshare

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/mposting"
	"sports_service/server/models/mshare"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"fmt"
	cloud "sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"time"
)

type ShareModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	posting     *mposting.PostingModel
	video       *mvideo.VideoModel
	community   *mcommunity.CommunityModel
	share       *mshare.ShareModel
}

func New(c *gin.Context) ShareModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return ShareModule{
		context: c,
		user: muser.NewUserModel(socket),
		posting: mposting.NewPostingModel(socket),
		video: mvideo.NewVideoModel(socket),
		community: mcommunity.NewCommunityModel(socket),
		share: mshare.NewShareModel(socket),
		engine: socket,
	}
}

// 分享/转发数据
func (svc *ShareModule) ShareData(params *mshare.ShareParams) int {
	now := int(time.Now().Unix())
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("post_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	switch params.SharePlatform {
	// 分享/转发 到微信、微博、qq todo: 记录即可
	case consts.SHARE_PLATFORM_WECHAT,consts.SHARE_PLATFORM_WEIBO,consts.SHARE_PLATFORM_QQ:

	// 分享到社区 则需发布一条新帖子
	case consts.SHARE_PLATFORM_COMMUNITY:
		client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
		// 检测帖子标题
		isPass, err := client.TextModeration(params.Title)
		if !isPass || err != nil {
			log.Log.Errorf("share_trace: validate title err: %s，pass: %v", err, isPass)
			return errdef.POST_INVALID_TITLE
		}

		// 检测帖子内容
		isPass, err = client.TextModeration(params.Describe)
		if !isPass || err != nil {
			log.Log.Errorf("share_trace: validate content err: %s，pass: %v", err, isPass)
			return errdef.POST_INVALID_CONTENT
		}

		user := svc.user.FindUserByUserid(params.UserId)
		if user == nil {
			log.Log.Errorf("share_trace: user not found, userId:%s", params.UserId)
			svc.engine.Rollback()
			return errdef.USER_NOT_EXISTS
		}

		section, err := svc.community.GetSectionInfo(fmt.Sprint(params.SectionId))
		if section == nil || err != nil {
			log.Log.Errorf("share_trace: section not found, id:%d", params.SectionId)
			svc.engine.Rollback()
			return errdef.COMMUNITY_SECTION_NOT_EXISTS
		}


		// 获取话题信息（多个）
		topics, err := svc.community.GetTopicByIds(params.TopicIds)
		if len(params.TopicIds) != len(topics) {
			log.Log.Errorf("share_trace: topic not found, topic_ids:%v, topics:%+v", params.TopicIds, topics)
			svc.engine.Rollback()
			return errdef.POST_TOPIC_NOT_EXISTS
		}

		switch params.ShareType {

		// 分享视频
		case consts.SHARE_VIDEO:
			video := svc.video.FindVideoById(fmt.Sprint(params.ComposeId))
			if video == nil {
				log.Log.Errorf("share_trace: video not found, videoId:%s", params.ComposeId)
				svc.engine.Rollback()
				return errdef.VIDEO_NOT_EXISTS
			}

			shareInfo := &mshare.ShareVideoInfo{
				VideoId: video.VideoId,
				Title: video.Title,
				Describe: video.Describe,
				Cover: video.Cover,
				VideoAddr: video.VideoAddr,
				VideoDuration: video.VideoDuration,
				CreateAt: video.CreateAt,
				UserId: video.UserId,
				Size: video.Size,
				//Labels:	svc.video.GetVideoLabels(fmt.Sprint(video.VideoId)),
			}

			shareInfo.Subarea, err = svc.video.GetSubAreaById(fmt.Sprint(video.Subarea))
			if err != nil {
				log.Log.Errorf("share_trace: get subarea by id fail, err:%s", err)
			}

			user = svc.user.FindUserByUserid(video.UserId)
			if user != nil {
				shareInfo.Nickname = user.NickName
				shareInfo.Avatar = user.Avatar
			}

			statistic := svc.video.GetVideoStatistic(fmt.Sprint(video.VideoId))
			shareInfo.BarrageNum = statistic.BarrageNum
			shareInfo.BrowseNum = statistic.BrowseNum

			// todo: 也可自己查
			videoInfo, _ := util.JsonFast.MarshalToString(shareInfo)
			svc.posting.Posting.Content = videoInfo
			// 记录类型为分享的视频
			svc.posting.Posting.ContentType = consts.COMMUNITY_FORWARD_VIDEO
			// 分享视频 则类型为视频+文本
			svc.posting.Posting.PostingType = consts.POST_TYPE_VIDEO
			// 关联的视频id
			svc.posting.Posting.VideoId = video.VideoId

		// 分享帖子
		case consts.SHARE_POST:
			post, err := svc.posting.GetPostById(fmt.Sprint(params.ComposeId))
			if post == nil || err != nil {
				log.Log.Errorf("share_trace: post not found, postId:%s", params.ComposeId)
				svc.engine.Rollback()
				return errdef.POST_NOT_EXISTS
			}

			shareInfo := &mshare.SharePostInfo{
				PostId: post.Id,
				PostingType: post.PostingType,
				Topics: topics,
				ContentType: post.ContentType,
				Title: post.Title,
				Describe: post.Describe,
				Content: post.Content,
				UserId: post.UserId,
			}

			user = svc.user.FindUserByUserid(post.UserId)
			if user != nil {
				shareInfo.Nickname = user.NickName
				shareInfo.Avatar = user.Avatar
			}

			statistic, err := svc.posting.GetPostStatistic(fmt.Sprint(post.Id))
			if err == nil && statistic != nil {
				shareInfo.BrowseNum = statistic.BrowseNum
				shareInfo.CommentNum = statistic.CommentNum
			} else {
				log.Log.Errorf("share_trace: get post statistic fail, err:%s", err)
			}

			postInfo, _ := util.JsonFast.MarshalToString(shareInfo)
			svc.posting.Posting.Content = postInfo
			// 记录类型为分享的帖子
			svc.posting.Posting.ContentType = consts.COMMUNITY_FORWARD_POST
			// 产品需求： 分享的帖子 皆为文本
			svc.posting.Posting.PostingType = consts.POST_TYPE_TEXT
		}


		svc.posting.Posting.Id = 0
		svc.posting.Posting.SectionId = section.Id
		svc.posting.Posting.UserId = params.UserId
		svc.posting.Posting.CreateAt = now
		svc.posting.Posting.Describe = params.Describe
		svc.posting.Posting.Title = params.Title
		svc.posting.Posting.CreateAt = now
		svc.posting.Posting.UpdateAt = now
		svc.posting.Posting.Status = 1
		if _, err := svc.posting.AddPost(); err != nil {
			svc.engine.Rollback()
			log.Log.Errorf("share_trace: add post fail, err:%s", err)
			return errdef.SHARE_DATA_FAIL
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
				log.Log.Errorf("share_trace: add post topics fail, err:%s", err)
				return errdef.SHARE_DATA_FAIL
			}
		}

		svc.posting.Statistic.PostingId = svc.posting.Posting.Id
		svc.posting.Statistic.CreateAt = now
		svc.posting.Statistic.UpdateAt = now
		// 初始化帖子统计数据
		if err := svc.posting.AddPostStatistic(); err != nil {
			log.Log.Errorf("share_trace: add post statistic err:%s", err)
			return errdef.SHARE_DATA_FAIL
		}

	}

	info, _ := util.JsonFast.Marshal(params)
	svc.share.Share.Content = string(info)
	svc.share.Share.UserId = params.UserId
	svc.share.Share.ComposeId = int64(params.ComposeId)
	svc.share.Share.ShareType = params.ShareType
	svc.share.Share.SharePlatform = params.SharePlatform
	svc.share.Share.CreateAt = now
	svc.share.Share.UpdateAt = now
	if _, err := svc.share.AddShareRecord(); err != nil {
		log.Log.Errorf("share_trace: record share err:%s", err)
		svc.engine.Rollback()
		return errdef.SHARE_DATA_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

