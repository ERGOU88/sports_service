package community

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/mlike"
	"sports_service/server/models/mposting"
	"fmt"
	"sports_service/server/models/muser"
	"sports_service/server/util"
)

type CommunityModule struct {
	context     *gin.Context
	engine      *xorm.Session
	community   *mcommunity.CommunityModel
	post        *mposting.PostingModel
	user        *muser.UserModel
	attention   *mattention.AttentionModel
	like        *mlike.LikeModel
}

func New(c *gin.Context) CommunityModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return CommunityModule{
		context: c,
		community: mcommunity.NewCommunityModel(socket),
		post: mposting.NewPostingModel(socket),
		user: muser.NewUserModel(socket),
		attention: mattention.NewAttentionModel(socket),
		like: mlike.NewLikeModel(socket),
		engine: socket,
	}
}

// 获取社区话题列表
func (svc *CommunityModule) GetCommunityTopics(isHot string, page, size int) (int, []*mcommunity.CommunityTopicInfo) {
	offset := (page - 1) * size
	list, err := svc.community.GetCommunityTopics(isHot, offset, size)
	if list == nil || err != nil {
		log.Log.Errorf("community_trace: get topics fail, err:%s", err)
		return errdef.COMMUNITY_TOPICS_FAIL, []*mcommunity.CommunityTopicInfo{}
	}

	res := make([]*mcommunity.CommunityTopicInfo, len(list))
	for index, item := range list {
		info := &mcommunity.CommunityTopicInfo{
			Id: item.Id,
			TopicName: item.TopicName,
			IsHot: item.IsHot,
		}

		num, err :=  svc.post.GetPostNumByTopic(fmt.Sprint(item.Id))
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
	list, err := svc.community.GetAllSection()
	if list == nil || err != nil {
		return errdef.COMMUNITY_SECTIONS_FAIL, []*mcommunity.CommunitySectionInfo{}
	}

	res := make([]*mcommunity.CommunitySectionInfo, len(list))
	for index, item := range list {
		info := &mcommunity.CommunitySectionInfo{
			Id: item.Id,
			SectionName: item.SectionName,
		}

		num, err :=  svc.post.GetPostNumBySection(fmt.Sprint(item.Id))
		if err != nil {
			log.Log.Errorf("community_trace: get post num by section fail, err:%s", err)
		}

		info.PostNum = num

		res[index] = info
	}

	return errdef.SUCCESS, res
}

// 通过id获取社区话题
func (svc *CommunityModule) GetCommunityTopicById(id string) (int, *mcommunity.CommunityTopicInfo) {
	info, err := svc.community.GetTopicInfo(id)
	if info == nil || err != nil {
		return errdef.COMMUNITY_TOPIC_FAIL, nil
	}

	res := &mcommunity.CommunityTopicInfo{
		Id: info.Id,
		TopicName: info.TopicName,
		IsHot: info.IsHot,
		Cover: info.Cover,
		Describe: info.Describe,
	}

	num, err :=  svc.post.GetPostNumByTopic(fmt.Sprint(info.Id))
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


	for _, item := range list {
		item.Topics, err = svc.post.GetPostTopic(fmt.Sprint(item.Id))
		if item.Topics == nil || err != nil  {
			item.Topics = []*models.PostingTopic{}
		}


		if user := svc.user.FindUserByUserid(item.UserId); user != nil {
			item.Avatar = user.Avatar
			item.Nickname = user.NickName
		}

		// 如果是转发的视频数据
		if item.ContentType == consts.COMMUNITY_FORWARD_VIDEO {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardVideo); err != nil {
				log.Log.Errorf("community_trace: get forward video info err:%s", err)
				return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
			}

		}

		// 如果是转发的帖子
		if item.PostingType == consts.POST_TYPE_TEXT && item.ContentType == consts.COMMUNITY_FORWARD_POST {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardPost); err != nil {
				log.Log.Errorf("community_trace: get forward post info err:%s", err)
				return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
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

	for _, item := range list {
		item.Topics, err = svc.post.GetPostTopic(fmt.Sprint(item.Id))
		if item.Topics == nil || err != nil  {
			item.Topics = []*models.PostingTopic{}
		}


		if user := svc.user.FindUserByUserid(item.UserId); user != nil {
			item.Avatar = user.Avatar
			item.Nickname = user.NickName
		}

		// 如果是转发的视频数据
		if item.ContentType == consts.COMMUNITY_FORWARD_VIDEO {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardVideo); err != nil {
				log.Log.Errorf("community_trace: get forward video info err:%s", err)
				return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
			}

		}

		// 如果是转发的帖子
		if item.PostingType == consts.POST_TYPE_TEXT && item.ContentType == consts.COMMUNITY_FORWARD_POST {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardPost); err != nil {
				log.Log.Errorf("community_trace: get forward post info err:%s", err)
				return errdef.COMMUNITY_POSTS_BY_TOPIC, []*mposting.PostDetailInfo{}
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


	return errdef.SUCCESS, list
}
