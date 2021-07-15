package cposting

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
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	cloud "sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"time"
	"github.com/microcosm-cc/bluemonday"
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
	client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测帖子标题
	isPass, err := client.TextModeration(params.Title)
	if !isPass || err != nil {
		log.Log.Errorf("post_trace: validate title err: %s，pass: %v", err, isPass)
		return errdef.POST_INVALID_TITLE
	}

	// 检测帖子内容
	isPass, err = client.TextModeration(params.Describe)
	if !isPass || err != nil {
		log.Log.Errorf("post_trace: validate content err: %s，pass: %v", err, isPass)
		return errdef.POST_INVALID_CONTENT
	}


	postType := svc.GetPostingType(params)
	if b := svc.VerifyContentLen(postType, params.Describe, params.Title); !b {
		log.Log.Error("post_trace: invalid content len")
		return errdef.POST_INVALID_CONTENT_LEN
	}

	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("post_trace: session begin err:%s", err)
		return errdef.ERROR
	}
	//
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("post_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 获取板块信息
	section, err := svc.community.GetSectionInfo(params.SectionId)
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

	// 默认为审核中的状态
	status := 0
	// 不带视频的帖子 只需通过图文检测
	if postType != consts.POST_TYPE_VIDEO {
		status = 1
	}

	now := int(time.Now().Unix())
	svc.posting.Posting.Title = params.Title
	svc.posting.Posting.Describe = params.Describe
	svc.posting.Posting.PostingType = postType
	svc.posting.Posting.Status = status
	// 社区发布
	svc.posting.Posting.ContentType = consts.COMMUNITY_PUB_POST
	svc.posting.Posting.CreateAt = now
	svc.posting.Posting.UpdateAt = now
	svc.posting.Posting.SectionId = section.Id
	svc.posting.Posting.UserId = userId
	if len(params.ImagesAddr) > 0 {
		bts, _ := util.JsonFast.Marshal(params.ImagesAddr)
		svc.posting.Posting.Content = string(bts)
	}

	if _, err := svc.posting.AddPost(); err != nil {
		log.Log.Errorf("post_trace: publish post fail, err:%s", err)
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
			log.Log.Errorf("post_trace: add post topics fail, err:%s", err)
			return errdef.POST_PUBLISH_FAIL
		}
	}

	svc.posting.Statistic.PostingId = svc.posting.Posting.Id
	svc.posting.Statistic.CreateAt = now
	svc.posting.Statistic.UpdateAt = now
	// 初始化帖子统计数据
	if err := svc.posting.AddPostStatistic(); err != nil {
		log.Log.Errorf("post_trace: add post statistic err:%s", err)
		return errdef.POST_PUBLISH_FAIL
	}


	svc.engine.Commit()

	return errdef.SUCCESS
}

// 获取帖子类型 todo: 常量
func (svc *PostingModule) GetPostingType(params *mposting.PostPublishParam) (postType int) {
	//if params.Describe != "" && params.VideoId != "" {
	//	// 视频 + 文字
	//	postType = consts.POST_TYPE_VIDEO
	//	return
	//}

	if params.Describe != "" && len(params.ImagesAddr) > 0 {
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
func (svc *PostingModule) VerifyContentLen(postType int, content, title string) bool {
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

// 获取帖子详情
func (svc *PostingModule) GetPostingDetail(postId string) (*models.PostingInfo, int) {
	detail, err := svc.posting.GetPostById(postId)
	if err != nil {
		return nil, errdef.POST_DETAIL_FAIL
	}

	return detail, errdef.SUCCESS
}

// 获取帖子详情页数据
func (svc *PostingModule) GetPostDetail(userId, postId string) (*mposting.PostDetailInfo, int) {
	if postId == "" {
		log.Log.Error("post_trace: postId can't empty")
		return nil, errdef.POST_NOT_EXISTS
	}

	post, err := svc.posting.GetPostById(postId)
	if err != nil {
		return nil, errdef.POST_NOT_EXISTS
	}

	if fmt.Sprint(post.Status) != consts.POST_AUDIT_SUCCESS {
		log.Log.Error("post_trace: post not audit, postId:%s", postId)
		return nil, errdef.POST_NOT_EXISTS
	}

	// todo: 完善返回数据
	resp := new(mposting.PostDetailInfo)
	resp.PostId = post.Id
	resp.Title = post.Title
	resp.Describe = post.Describe
	resp.IsRecommend = post.IsRecommend
	resp.IsTop = post.IsTop
	resp.CreateAt = post.CreateAt
	resp.UserId = post.UserId
	resp.ContentType = post.ContentType
	resp.PostingType = post.PostingType
	resp.Topics, err = svc.posting.GetPostTopic(postId)
	if resp.Topics == nil || err != nil  {
		resp.Topics = []*models.PostingTopic{}
	}

	now := int(time.Now().Unix())
	// 增加帖子浏览总数
	if err := svc.posting.UpdatePostBrowseNum(post.Id, now, 1); err != nil {
		log.Log.Errorf("post_trace: update post browse num err:%s", err)
	}

	if user := svc.user.FindUserByUserid(post.UserId); user != nil {
		resp.Avatar = user.Avatar
		resp.Nickname = user.NickName
	}

	// 如果是转发的视频数据
	if resp.ContentType == consts.COMMUNITY_FORWARD_VIDEO {
		if err = util.JsonFast.UnmarshalFromString(post.Content, &resp.ForwardVideo); err != nil {
			log.Log.Errorf("post_trace: get forward video info err:%s", err)
			return nil, errdef.POST_DETAIL_FAIL
		}

	}

	// 如果是转发的帖子
	if resp.PostingType == consts.POST_TYPE_TEXT && resp.ContentType == consts.COMMUNITY_FORWARD_POST {
		if err = util.JsonFast.UnmarshalFromString(post.Content, &resp.ForwardPost); err != nil {
			log.Log.Errorf("post_trace: get forward post info err:%s", err)
			return nil, errdef.POST_DETAIL_FAIL
		}


	}

	if userId == "" {
		log.Log.Error("post_trace: user no login")
		return resp, errdef.SUCCESS
	}

	// 获取用户信息
	if user := svc.user.FindUserByUserid(userId); user != nil {
		// 用户是否浏览过
		browse := svc.posting.GetUserBrowsePost(userId, consts.TYPE_POST, post.Id)
		if browse != nil {
			svc.video.Browse.CreateAt = now
			svc.video.Browse.UpdateAt = now
			// 已有浏览记录 更新用户浏览的时间
			if err := svc.posting.UpdateUserBrowsePost(userId, consts.TYPE_POST, post.Id); err != nil {
				log.Log.Errorf("post_trace: update user browse post err:%s", err)
			}
		} else {
			svc.video.Browse.CreateAt = now
			svc.video.Browse.UpdateAt = now
			svc.video.Browse.UserId = userId
			svc.video.Browse.ComposeId = post.Id
			svc.video.Browse.ComposeType = consts.TYPE_POST
			// 添加用户浏览的帖子记录
			if err := svc.posting.RecordUserBrowsePost(); err != nil {
				log.Log.Errorf("post_trace: record user browse post err:%s", err)
			}
		}
	}
	// 获取视频相关统计数据
	info := svc.video.GetVideoStatistic(fmt.Sprint(post.Id))
	resp.BrowseNum = info.BrowseNum
	resp.CommentNum = info.CommentNum
	resp.FabulousNum = info.FabulousNum
	resp.ShareNum = info.ShareNum


	return resp, errdef.SUCCESS
}

// 过滤富文本 todo：和客户端确认最终的策略
func (svc *PostingModule) SanitizeHtml(content string) string {
	p := bluemonday.NewPolicy()
	p.AllowStandardURLs()

	// 只允许<body> <p> 和 <a href="">
	p.AllowAttrs("href").OnElements("a")
	p.AllowElements("p")
	p.AllowElements("body")

	return p.Sanitize(content)

}


//func (svc *PostingModule) SanitizeHtml(content string) (string, error) {
//	config := `
//	{
//		"stripWhitespace": true,
//		"elements": {
//			"body": [],
//			"i": [],
//            "p": [],
//            "a":[]
//		}
//	}`
//
//	whitelist, err := sanitize.NewWhitelist([]byte(config))
//	if err != nil {
//		return "", err
//	}
//
//	fmt.Println(whitelist)
//
//	readStr := strings.NewReader(content)
//	fmt.Println(readStr)
//	sanitized, err := whitelist.SanitizeUnwrap(readStr)
//	if err != nil {
//		return "", err
//	}
//
//	log.Log.Debugf("sanitized html: %s", sanitized)
//	return sanitized, nil
//}
