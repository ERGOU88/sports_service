package cpost

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/mposting"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	redismq "sports_service/server/redismq/event"
	"sports_service/server/util"
	"sports_service/server/global/backend/log"
	"time"
)

type PostModule struct {
	context      *gin.Context
	engine       *xorm.Session
	post         *mposting.PostingModel
	attention    *mattention.AttentionModel
	user         *muser.UserModel
	video        *mvideo.VideoModel
	community    *mcommunity.CommunityModel
}

func New(c *gin.Context) PostModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return PostModule{
		context: c,
		post: mposting.NewPostingModel(socket),
		attention: mattention.NewAttentionModel(socket),
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		community: mcommunity.NewCommunityModel(socket),
		engine: socket,
	}
}

// 帖子审核
func (svc *PostModule) AudiPost(param *mposting.AudiPostParam) int {
	post, err := svc.post.GetPostById(param.Id)
	if post == nil || err != nil {
		return errdef.POST_NOT_FOUND
	}

	status := fmt.Sprint(post.Status)
	// 帖子已删除
	if status == consts.POST_DELETE_STATUS {
		svc.engine.Rollback()
		return errdef.POST_ALREADY_DELETE
	}

	// 帖子已通过审核 只能执行逻辑删除
	if status == consts.POST_AUDIT_SUCCESS && fmt.Sprint(param.Status) != consts.POST_DELETE_STATUS {
		svc.engine.Rollback()
		return errdef.POST_ALREADY_PASS
	}

	// 通过 / 不通过 / 执行删除操作 且 视频状态为审核通过 则只能逻辑删除/不通过 直接更新视频状态
	if fmt.Sprint(param.Status) == consts.POST_AUDIT_SUCCESS || fmt.Sprint(param.Status) == consts.POST_AUDIT_FAILURE ||
		(fmt.Sprint(param.Status) == consts.POST_DELETE_STATUS && status == consts.POST_AUDIT_SUCCESS) {
		post.Status = param.Status
		// 更新视频状态
		if err := svc.post.UpdateStatusByPost(); err != nil {
			svc.engine.Rollback()
			return errdef.POST_EDIT_STATUS_FAIL
		}

		// 如果是审核通过
		if fmt.Sprint(param.Status) == consts.POST_AUDIT_SUCCESS {
			// 获取发布者用户信息
			user := svc.user.FindUserByUserid(post.UserId)
			if user != nil {
				// 获取发布者粉丝们的userId
				userIds := svc.attention.GetFansList(user.UserId)
				for _, userId := range userIds {
					// 给发布者的粉丝 发送 发布新帖子推送
					//event.PushEventMsg(config.Global.AmqpDsn, userId, user.NickName, video.Cover, "", consts.FOCUS_USER_PUBLISH_VIDEO_MSG)
					redismq.PushEventMsg(redismq.NewEvent(userId, fmt.Sprint(post.Id), user.NickName, "", "", consts.FOCUS_USER_PUBLISH_POST_MSG))
				}
			}
		}

		svc.engine.Commit()
		return errdef.SUCCESS
	}

	// 如果执行删除操作 且 视频状态未审核通过 删除相关所有数据
	if fmt.Sprint(param.Status) == consts.VIDEO_DELETE_STATUS && status != consts.VIDEO_AUDIT_SUCCESS {
		// 物理删除发布的帖子、帖子所属话题、帖子统计数据
		if err := svc.post.DelPublishPostById(param.Id); err != nil {
			svc.engine.Rollback()
			return errdef.POST_DELETE_PUBLISH_FAIL
		}

		// 删除帖子所属话题
		if err := svc.post.DelPostTopics(param.Id); err != nil {
			svc.engine.Rollback()
			return errdef.POST_DELETE_TOPIC_FAIL
		}

		// 删除帖子统计数据
		if err := svc.post.DelPostStatistic(param.Id); err != nil {
			svc.engine.Rollback()
			return errdef.POST_DELETE_STATISTIC_FAIL
		}
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 管理后台获取帖子列表
func (svc *PostModule) GetPostList(page, size int) (int, []*mposting.PostDetailInfo) {
	offset := (page - 1) * size
	list, err := svc.post.GetPostList(offset, size)
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*mposting.PostDetailInfo{}
	}

	for _, item := range list {
		var err error
		item.Topics, err = svc.post.GetPostTopic(fmt.Sprint(item.Id))
		if item.Topics == nil || err != nil {
			item.Topics = []*models.PostingTopic{}
		}

		user := svc.user.FindUserByUserid(item.UserId)
		if user != nil {
			item.Avatar = user.Avatar
			item.Nickname = user.NickName
		}

		// 如果是转发的视频数据
		if item.ContentType == consts.COMMUNITY_FORWARD_VIDEO {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardVideo); err != nil {
				log.Log.Errorf("post_trace: get forward video info err:%s", err)
				//return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
			} else {
				item.ForwardVideo.VideoAddr = svc.video.AntiStealingLink(item.ForwardVideo.VideoAddr)
			}

		}

		// 如果是转发的帖子
		if item.PostingType == consts.POST_TYPE_TEXT && item.ContentType == consts.COMMUNITY_FORWARD_POST {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardPost); err != nil {
				log.Log.Errorf("post_trace: get forward post info err:%s", err)
				//return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
			}

			// 如果转发的是图文类型 需要展示图文
			if item.ForwardPost.PostingType == consts.POST_TYPE_IMAGE {
				if err := util.JsonFast.UnmarshalFromString(item.ForwardPost.Content, &item.ForwardPost.ImagesAddr); err != nil {
					log.Log.Errorf("post_trace: get images by forward post fail, err:%s", err)
				}
			}

		}

		// 图文帖
		if item.PostingType == consts.POST_TYPE_IMAGE {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ImagesAddr); err != nil {
				log.Log.Errorf("post_trace: get image info err:%s", err)
				//return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
			}
		}

		// 如果视频+文 的帖子 且 为社区发布 查询关联的视频信息
		if item.PostingType == consts.POST_TYPE_VIDEO && item.ContentType == consts.COMMUNITY_PUB_POST {
			video := svc.video.FindVideoById(fmt.Sprint(item.VideoId))
			if video == nil {
				log.Log.Errorf("post_trace: get video info err:%s, videoId:%s", err, item.VideoId)
			} else {
				item.RelatedVideo = new(mposting.RelatedVideo)
				item.RelatedVideo.VideoId = video.VideoId
				item.RelatedVideo.UserId = video.UserId
				item.RelatedVideo.CreateAt = video.CreateAt
				item.RelatedVideo.Describe = video.Describe
				item.RelatedVideo.Cover = video.Cover
				item.RelatedVideo.Title = video.Title
				item.RelatedVideo.VideoDuration = video.VideoDuration
				item.RelatedVideo.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
				item.RelatedVideo.Size = video.Size

				statistic := svc.video.GetVideoStatistic(fmt.Sprint(video.VideoId))
				if statistic != nil {
					item.RelatedVideo.FabulousNum = statistic.FabulousNum
					item.RelatedVideo.CommentNum = statistic.CommentNum
					item.RelatedVideo.ShareNum = statistic.ShareNum
				}

				if user != nil {
					item.RelatedVideo.Nickname = user.NickName
					item.RelatedVideo.Avatar = user.Avatar
				}

				subarea, err := svc.video.GetSubAreaById(fmt.Sprint(video.Subarea))
				if err != nil || subarea == nil {
					log.Log.Errorf("post_trace: get subarea by id fail, err:%s", err)
				} else {
					item.RelatedVideo.Subarea = subarea
				}

			}
		}

		sectionInfo, err := svc.community.GetSectionInfo(fmt.Sprint(item.SectionId))
		if err == nil {
			item.SectionName = sectionInfo.SectionName
		}

		item.Content = ""

	}

	return errdef.SUCCESS, list
}

func (svc *PostModule) AddSection(param *mcommunity.AddSection) int {
	svc.community.CommunitySection.SectionName = param.SectionName
	svc.community.CommunitySection.Sortorder = param.Sortorder
	svc.community.CommunitySection.CreateAt = int(time.Now().Unix())
	svc.community.CommunitySection.Status = 1
	if _, err := svc.community.AddCommunitySection(); err != nil {
		log.Log.Errorf("post_trace: add section fail, err:%s", err)
		return errdef.POST_ADD_SECTION_FAIL
	}

	return errdef.SUCCESS
}

// 软删除 将板块隐藏
func (svc *PostModule) DelSection(param *mcommunity.DelSection) int {
	svc.community.CommunitySection.Status = 2
	if _, err := svc.community.UpdateSectionStatus(param.Id); err != nil {
		log.Log.Errorf("post_trace: del section fail, err:%s", err)
		return errdef.POST_DEL_SECTION_FAIL
	}

	return errdef.SUCCESS
}

func (svc *PostModule) AddTopic(param *mcommunity.AddTopic) int {
	svc.community.CommunityTopic.Status = 1
	svc.community.CommunityTopic.Cover = param.Cover
	svc.community.CommunityTopic.CreateAt = int(time.Now().Unix())
	svc.community.CommunityTopic.Sortorder = param.Sortorder
	svc.community.CommunityTopic.TopicName = param.Name
	svc.community.CommunityTopic.Describe = param.Describe
	svc.community.CommunityTopic.SectionId = 1
	if _, err := svc.community.AddTopic(); err != nil {
		log.Log.Errorf("post_trace: add topic fail, err:%s", err)
		return errdef.POST_ADD_TOPIC_FAIL
	}

	return errdef.SUCCESS
}

func (svc *PostModule) DelTopic(param *mcommunity.DelTopic) int {
	svc.community.CommunityTopic.Status = 2
	if _, err := svc.community.UpdateTopicStatus(param.Id); err != nil {
		log.Log.Errorf("post_trace: update topic status fail, err:%s", err)
		return errdef.POST_DEL_TOPIC_FAIL
	}

	return errdef.SUCCESS
}

// 帖子设置 置顶/精华
func (svc *PostModule) PostSetting(param *mposting.SettingParam) int {
	var cols string
	switch param.SettingType {
	case 1:
		cols = "is_cream"
		svc.post.Posting.IsCream = param.ActionType
		svc.post.ApplyCream.Status = param.ActionType
		if param.ActionType == 0 {
			svc.post.ApplyCream.Status = 2
		}

		if _, err := svc.post.UpdateApplyCreamStatus(param.Id); err != nil {
			log.Log.Errorf("post_trace: update apply cream status fail, err:%s", err)
			return errdef.POST_SETTING_FAIL
		}
	case 2:
		cols = "is_top"
		svc.post.Posting.IsTop = param.ActionType
	}

	if _, err := svc.post.UpdatePostInfo(param.Id, cols); err != nil {
		log.Log.Errorf("post_trace: update post info fail, err:%s", err)
		return errdef.POST_SETTING_FAIL
	}

	return errdef.SUCCESS
}

// 申精列表
func (svc *PostModule) GetApplyCreamList(page, size int) (int, []*mposting.PostDetailInfo) {
	offset := (page - 1) * size
	list, err := svc.post.GetApplyCreamList(offset, size)
	if err != nil {
		return errdef.POST_APPLY_CREAM_LIST_FAIL, nil
	}

	for _, item := range list {
		item.Topics, err = svc.post.GetPostTopic(fmt.Sprint(item.Id))
		if item.Topics == nil || err != nil {
			item.Topics = []*models.PostingTopic{}
		}

		sectionInfo, err := svc.community.GetSectionInfo(fmt.Sprint(item.SectionId))
		if err == nil {
			item.SectionName = sectionInfo.SectionName
		}

	}

	return errdef.SUCCESS, list
}

// 板块列表
func (svc *PostModule) GetSectionList() (int, []*models.CommunitySection) {
	list, err := svc.community.GetAllSection()
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*models.CommunitySection{}
	}

	return errdef.SUCCESS, list
}

// 话题列表
func (svc *PostModule) GetTopicList() (int, []*models.CommunityTopic) {
	list, err := svc.community.GetAllTopic()
	if err != nil {
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*models.CommunityTopic{}
	}

	return errdef.SUCCESS, list
}
