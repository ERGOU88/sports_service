package cposting

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/mposting"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"strings"
	"time"
	"sports_service/server/tools/sanitize"
	"fmt"
)

type PostingModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	posting     *mposting.PostingModel
	video       *mvideo.VideoModel
	community   *mcommunity.CommunityModel
}

func New(c *gin.Context) PostingModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return PostingModule{
		context: c,
		user: muser.NewUserModel(socket),
		posting: mposting.NewPostingModel(socket),
		video: mvideo.NewVideoModel(socket),
		community: mcommunity.NewCommunityModel(socket),
		engine: socket,
	}
}

// 发布帖子
func (svc *PostingModule) PublishPosting(userId string, params *mposting.PostPublishParam) int {
	//client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测帖子标题
	//isPass, err := client.TextModeration(params.Title)
	//if !isPass || err != nil {
	//	log.Log.Errorf("post_trace: validate title err: %s，pass: %v", err, isPass)
	//	return errdef.POST_INVALID_TITLE
	//}

	// 检测帖子内容
	//isPass, err = client.TextModeration(params.Content)
	//if !isPass || err != nil {
	//	log.Log.Errorf("post_trace: validate content err: %s，pass: %v", err, isPass)
	//	return errdef.POST_INVALID_CONTENT
	//}

	//postType := svc.GetPostingType(params)
	//if b := svc.VerifyContentLen(postType, params.Content); !b {
	//	log.Log.Error("post_trace: invalid content len")
	//	return errdef.POST_INVALID_CONTENT_LEN
	//}

	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("post_trace: session begin err:%s", err)
		return errdef.ERROR
	}
	//
	//if user := svc.user.FindUserByUserid(userId); user == nil {
	//	log.Log.Errorf("post_trace: user not found, userId:%s", userId)
	//	svc.engine.Rollback()
	//	return errdef.USER_NOT_EXISTS
	//}

	// 获取板块信息
	//section, err := svc.community.GetSectionInfo(params.SectionId)
	//if section == nil || err != nil {
	//	log.Log.Errorf("post_trace: section not found, sectionId:%d", params.SectionId)
	//	svc.engine.Rollback()
	//	return errdef.POST_SECTION_NOT_EXISTS
	//}

	// 获取话题信息（多个）
	//topics, err := svc.community.GetTopicByIds(strings.Join(params.TopicIds, ","))
	//if len(params.TopicIds) != len(topics) {
	//	log.Log.Errorf("post_trace: topic not found, topic_ids:%v, topics:%+v", params.TopicIds, topics)
	//	svc.engine.Rollback()
	//	return errdef.POST_TOPIC_NOT_EXISTS
	//}

	now := int(time.Now().Unix())
	svc.posting.Posting.Title = params.Title
	svc.posting.Posting.Content = params.Content
	svc.posting.Posting.VideoAddr = params.VideoAddr
	svc.posting.Posting.VideoDuration = params.VideoDuration
	svc.posting.Posting.Cover = params.Cover
	//svc.posting.Posting.PostingType = postType
	svc.posting.Posting.ContentType = params.ContentType
	svc.posting.Posting.CreateAt = now
	svc.posting.Posting.UpdateAt = now
	// todo: 常量 暂时只有官方社区 考虑到扩展性 已设计社区表
	svc.posting.Posting.CommunityId = 1
	if svc.posting.Posting.ImagesAddr != "" {
		svc.posting.Posting.ImagesAddr = strings.Join(params.ImagesAddr, ",")
	}

	// todo: 如上传了视频的贴子 需等vod事件流审核完成 修改贴子状态为审核完成
	if _, err := svc.posting.AddPost(); err != nil {
		log.Log.Errorf("post_trace: publish post fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.POST_PUBLISH_FAIL
	}

	// 组装多条记录 写入帖子话题表
	//topicInfos := make([]*models.PostingTopic, len(topics))
	//for _, val := range topics {
	//	info := new(models.PostingTopic)
	//	info.TopicId = val.TopicId
	//	info.TopicName = val.TopicName
	//	info.CreateAt = now
	//	info.PostingId = svc.posting.Posting.Id
	//	topicInfos = append(topicInfos, info)
	//}
	//
	//if len(topicInfos) > 0 {
	//	// 添加帖子所属话题（多条）
	//	affected, err := svc.posting.AddPostingTopics(topicInfos)
	//	if err != nil || int(affected) != len(topicInfos) {
	//		svc.engine.Rollback()
	//		log.Log.Errorf("post_trace: add post topics fail, err:%s", err)
	//		return errdef.POST_PUBLISH_FAIL
	//	}
	//}


	svc.engine.Commit()

	return errdef.SUCCESS
}

// 获取帖子类型 todo: 常量
func (svc *PostingModule) GetPostingType(params *mposting.PostPublishParam) (postType int) {
	if params.Content != "" && params.VideoAddr != "" && params.VideoDuration > 0 {
		// 视频 + 文字
		postType = 2
		return
	}

	if params.Content != "" && len(params.ImagesAddr) > 0 {
		// 图文
		postType = 1
		return
	}

	if params.Content != "" {
		// 纯文本
		postType = 0
	}

	return
}

// 验证正文长度
func (svc *PostingModule) VerifyContentLen(postType int, content string) bool {
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

// 获取帖子详情
func (svc *PostingModule) GetPostingDetail(postId string) (*models.PostingInfo, int) {
	detail, err := svc.posting.GetPostById(postId)
	if err != nil {
		return nil, errdef.POST_DETAIL_FAIL
	}

	return detail, errdef.SUCCESS
}

// 过滤富文本
func (svc *PostingModule) SanitizeHtml(content string) (string, error) {
	whitelist, err := sanitize.WhitelistFromFile("../../config/white_dom.json")
	if err != nil {
		return "", err
	}

	fmt.Println(whitelist)

	readStr := strings.NewReader(content)
	fmt.Println(readStr)
	sanitized, err := whitelist.SanitizeUnwrap(readStr)
	if err != nil {
		return "", err
	}

	log.Log.Debugf("sanitized html: %s", sanitized)
	return sanitized, nil
}
