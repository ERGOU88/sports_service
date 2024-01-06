package community

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/dao"
	"sports_service/global/app/errdef"
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/models"
	"sports_service/models/mattention"
	"sports_service/models/mcommunity"
	"sports_service/models/mlike"
	"sports_service/models/mposting"
	"sports_service/models/muser"
	"sports_service/models/mvideo"
	"sports_service/tools/tencentCloud"
	"sports_service/util"
	"strings"
)

type CommunityModule struct {
	context   *gin.Context
	engine    *xorm.Session
	community *mcommunity.CommunityModel
	post      *mposting.PostingModel
	user      *muser.UserModel
	attention *mattention.AttentionModel
	like      *mlike.LikeModel
	video     *mvideo.VideoModel
}

func New(c *gin.Context) CommunityModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return CommunityModule{
		context:   c,
		community: mcommunity.NewCommunityModel(socket),
		post:      mposting.NewPostingModel(socket),
		user:      muser.NewUserModel(socket),
		attention: mattention.NewAttentionModel(socket),
		like:      mlike.NewLikeModel(socket),
		video:     mvideo.NewVideoModel(socket),
		engine:    socket,
	}
}

// 获取板块下的置顶帖
func (svc *CommunityModule) GetTopPostBySectionId(page, size int, sectionId string) (int, []*mposting.TopPost) {
	offset := (page - 1) * size
	list, err := svc.post.GetTopPostBySectionId(offset, size, sectionId)
	if err != nil {
		log.Log.Errorf("community_trace: get top post by section fail, err:%s", err)
		return errdef.ERROR, []*mposting.TopPost{}
	}

	return errdef.SUCCESS, list
}

// 获取社区话题列表 [按话题帖子数量排序]
func (svc *CommunityModule) GetTopicListOrderByPostNum(page, size int) (int, []*mcommunity.CommunityTopicInfo) {
	offset := (page - 1) * size
	list, err := svc.community.GetTopicListOrderByPostNum(offset, size)
	if err != nil {
		log.Log.Errorf("community_trace: get topics fail, err:%s", err)
		return errdef.COMMUNITY_TOPICS_FAIL, []*mcommunity.CommunityTopicInfo{}
	}

	if list == nil {
		return errdef.SUCCESS, []*mcommunity.CommunityTopicInfo{}
	}

	return errdef.SUCCESS, list
}

// 获取社区话题列表
func (svc *CommunityModule) GetCommunityTopics(sectionId, isHot string, page, size int) (int, []*mcommunity.CommunityTopicInfo) {
	offset := (page - 1) * size
	list, err := svc.community.GetCommunityTopics(sectionId, isHot, offset, size)
	if list == nil || err != nil {
		log.Log.Errorf("community_trace: get topics fail, err:%s", err)
		return errdef.COMMUNITY_TOPICS_FAIL, []*mcommunity.CommunityTopicInfo{}
	}

	res := make([]*mcommunity.CommunityTopicInfo, len(list))
	for index, item := range list {
		info := &mcommunity.CommunityTopicInfo{
			Id:        item.Id,
			TopicName: item.TopicName,
			IsHot:     item.IsHot,
		}

		num, err := svc.post.GetPostNumByTopic(fmt.Sprint(item.Id))
		if err != nil {
			log.Log.Errorf("community_trace: get post num by topic fail, err:%s", err)
		}

		info.PostNum = num

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 获取社区板块
func (svc *CommunityModule) GetCommunitySections() (int, []*mcommunity.CommunitySectionInfo) {
	condition := "status=1"
	list, err := svc.community.GetAllSection(condition)
	if list == nil || err != nil {
		return errdef.COMMUNITY_SECTIONS_FAIL, []*mcommunity.CommunitySectionInfo{}
	}

	res := make([]*mcommunity.CommunitySectionInfo, len(list))
	for index, item := range list {
		info := &mcommunity.CommunitySectionInfo{
			Id:          item.Id,
			SectionName: item.SectionName,
		}

		num := svc.GetPostNumBySection(fmt.Sprint(item.Id))

		info.PostNum = num

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 获取板块下的帖子数
func (svc *CommunityModule) GetPostNumBySection(sectionId string) int64 {
	num, err := svc.post.GetPostNumBySection(sectionId)
	if err != nil {
		log.Log.Errorf("community_trace: get post num by section fail, err:%s", err)
		return 0
	}

	return num
}

// 通过id获取社区话题
func (svc *CommunityModule) GetCommunityTopicById(id string) (int, *mcommunity.CommunityTopicInfo) {
	info, err := svc.community.GetTopicInfo(id)
	if info == nil || err != nil {
		return errdef.COMMUNITY_TOPIC_FAIL, nil
	}

	res := &mcommunity.CommunityTopicInfo{
		Id:        info.Id,
		TopicName: info.TopicName,
		IsHot:     info.IsHot,
		Cover:     tencentCloud.BucketURI(info.Cover),
		Describe:  info.Describe,
	}

	num, err := svc.post.GetPostNumByTopic(fmt.Sprint(info.Id))
	if err != nil {
		log.Log.Errorf("community_trace: get post num by topic fail, err:%s", err)
	}

	res.PostNum = num

	return errdef.SUCCESS, res
}

// 获取板块下的帖子列表
func (svc *CommunityModule) GetPostListBySection(page, size int, userId, sectionId string) (int, []*mposting.PostDetailInfo) {
	section, err := svc.community.GetSectionInfo(sectionId)
	if section == nil || err != nil {
		log.Log.Errorf("community_trace: section not found, sectionId:%s", sectionId)
		return errdef.COMMUNITY_SECTION_NOT_EXISTS, []*mposting.PostDetailInfo{}
	}

	offset := (page - 1) * size
	list, err := svc.post.GetPostListBySectionId(sectionId, offset, size)
	if err != nil {
		log.Log.Errorf("community_trace: get post list by section fail, err:%s", err)
		return errdef.COMMUNITY_POSTS_BY_SECTION, []*mposting.PostDetailInfo{}
	}

	if list == nil {
		return errdef.SUCCESS, []*mposting.PostDetailInfo{}
	}

	svc.GetPostDetailByList(userId, list)

	return errdef.SUCCESS, list

}

// 获取话题下的帖子列表 sortHot 1 按热度排序 0 按发布时间排序
func (svc *CommunityModule) GetPostListByTopic(page, size int, userId, topicId, sortHot string) (int, []*mposting.PostDetailInfo) {
	topic, err := svc.community.GetTopicInfo(topicId)
	if topic == nil || err != nil {
		log.Log.Errorf("community_trace: topic not found, topicId:%s", topicId)
		return errdef.COMMUNITY_TOPIC_NOT_EXISTS, []*mposting.PostDetailInfo{}
	}

	offset := (page - 1) * size
	list, err := svc.post.GetPostListByTopicId(topicId, sortHot, offset, size)
	if err != nil {
		log.Log.Errorf("community_trace: get post list by topic fail, err:%s", err)
		return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
	}

	if list == nil {
		return errdef.SUCCESS, []*mposting.PostDetailInfo{}
	}

	svc.GetPostDetailByList(userId, list)

	return errdef.SUCCESS, list
}

// 获取帖子列表 详情数据
func (svc *CommunityModule) GetPostDetailByList(userId string, list []*mposting.PostDetailInfo) []*mposting.PostDetailInfo {
	if len(list) == 0 {
		return []*mposting.PostDetailInfo{}
	}

	for _, item := range list {
		var err error
		item.Topics, err = svc.post.GetPostTopic(fmt.Sprint(item.Id))
		if item.Topics == nil || err != nil {
			item.Topics = []*models.PostingTopic{}
		}

		user := svc.user.FindUserByUserid(item.UserId)
		if user != nil {
			item.Avatar = tencentCloud.BucketURI(user.Avatar)
			item.Nickname = user.NickName
		}

		// 如果是转发的视频数据
		if item.ContentType == consts.COMMUNITY_FORWARD_VIDEO {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardVideo); err != nil {
				log.Log.Errorf("community_trace: get forward video info err:%s", err)
				//return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
			} else {
				item.ForwardVideo.VideoAddr = svc.video.AntiStealingLink(item.ForwardVideo.VideoAddr)
			}

		}

		// 如果是转发的帖子
		if item.PostingType == consts.POST_TYPE_TEXT && item.ContentType == consts.COMMUNITY_FORWARD_POST {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardPost); err != nil {
				log.Log.Errorf("community_trace: get forward post info err:%s", err)
				//return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
			}

			// 如果转发的是图文类型 需要展示图文
			if item.ForwardPost.PostingType == consts.POST_TYPE_IMAGE {
				if err := util.JsonFast.UnmarshalFromString(item.ForwardPost.Content, &item.ForwardPost.ImagesAddr); err != nil {
					log.Log.Errorf("community_trace: get images by forward post fail, err:%s", err)
				}
			}

		}

		// 图文帖
		if item.PostingType == consts.POST_TYPE_IMAGE {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ImagesAddr); err != nil {
				log.Log.Errorf("community_trace: get image info err:%s", err)
				//return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
			}
		}

		// 如果视频+文 的帖子 且 为社区发布 查询关联的视频信息
		if item.PostingType == consts.POST_TYPE_VIDEO && item.ContentType == consts.COMMUNITY_PUB_POST {
			video := svc.video.FindVideoById(fmt.Sprint(item.VideoId))
			if video == nil {
				log.Log.Errorf("community_trace: get video info err:%s, videoId:%s", err, item.VideoId)
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
					item.RelatedVideo.Avatar = tencentCloud.BucketURI(user.Avatar)
				}

				// 是否点赞
				if likeInfo := svc.like.GetLikeInfo(userId, video.VideoId, consts.TYPE_VIDEOS); likeInfo != nil {
					item.RelatedVideo.IsLike = likeInfo.Status
				}

				subarea, err := svc.video.GetSubAreaById(fmt.Sprint(video.Subarea))
				if err != nil || subarea == nil {
					log.Log.Errorf("community_trace: get subarea by id fail, err:%s", err)
				} else {
					item.RelatedVideo.Subarea = subarea
				}

			}
		}

		item.Content = ""

		if userId == "" {
			continue
		}

		// 是否关注
		if attentionInfo := svc.attention.GetAttentionInfo(userId, item.UserId); attentionInfo != nil {
			item.IsAttention = attentionInfo.Status
		}

		// 是否点赞
		if likeInfo := svc.like.GetLikeInfo(userId, item.Id, consts.TYPE_POSTS); likeInfo != nil {
			item.IsLike = likeInfo.Status
		}

	}

	return list
}

// 获取关注的人发布的帖子
func (svc *CommunityModule) GetPostListByAttention(userId string, page, size int) []*mposting.PostDetailInfo {
	// 用户未登录
	if userId == "" {
		log.Log.Error("post_trace: no login")
		return []*mposting.PostDetailInfo{}
	}

	userIds := svc.attention.GetAttentionList(userId)
	if len(userIds) == 0 {
		log.Log.Errorf("post_trace: not following any users")
		return []*mposting.PostDetailInfo{}
	}

	offset := (page - 1) * size
	uids := strings.Join(userIds, ",")
	list, err := svc.post.GetPostListByAttention(uids, offset, size)
	if err != nil {
		log.Log.Errorf("post_trace: get post by attention fail, err:%s", err)
		return []*mposting.PostDetailInfo{}
	}

	if len(list) == 0 {
		return []*mposting.PostDetailInfo{}
	}

	return svc.GetPostDetailByList(userId, list)
}
