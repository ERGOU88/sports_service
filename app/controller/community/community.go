package community

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/mposting"
	"fmt"
)

type CommunityModule struct {
	context     *gin.Context
	engine      *xorm.Session
	community   *mcommunity.CommunityModel
	post        *mposting.PostingModel
}

func New(c *gin.Context) CommunityModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return CommunityModule{
		context: c,
		community: mcommunity.NewCommunityModel(socket),
		post: mposting.NewPostingModel(socket),
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
