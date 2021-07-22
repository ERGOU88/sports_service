package cposting

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
	attention   *mattention.AttentionModel
	like        *mlike.LikeModel
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
		attention: mattention.NewAttentionModel(socket),
		like: mlike.NewLikeModel(socket),
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
		svc.engine.Rollback()
		return errdef.POST_PUBLISH_FAIL
	}

	// 添加@
	if len(params.AtInfo) > 0 {
		var atList []*models.ReceivedAt
		for _, val := range params.AtInfo {
			user := svc.user.FindUserByUserid(val)
			if user == nil {
				log.Log.Errorf("post_trace: at user not found, userId:%s", val)
				continue
			}

			at := &models.ReceivedAt{
				ToUserId: val,
				UserId: userId,
				ComposeId: svc.posting.Posting.Id,
				TopicType: consts.TYPE_POST_COMMENT,
				CreateAt: now,
			}

			atList = append(atList, at)
		}

		affected, err := svc.posting.AddReceiveAtList(atList)
		if err != nil || int(affected) != len(atList) {
			log.Log.Errorf("post_trace: add receive at list fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.POST_PUBLISH_FAIL
		}
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
//func (svc *PostingModule) GetPostingDetail(postId string) (*models.PostingInfo, int) {
//	detail, err := svc.posting.GetPostById(postId)
//	if err != nil {
//		return nil, errdef.POST_DETAIL_FAIL
//	}
//
//	return detail, errdef.SUCCESS
//}

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
	resp.Id = post.Id
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

	// 图文帖
	if resp.PostingType == consts.POST_TYPE_IMAGE {
		if err = util.JsonFast.UnmarshalFromString(post.Content, &resp.ImagesAddr); err != nil {
			log.Log.Errorf("post_trace: get image info err:%s", err)
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
			svc.posting.Browse.CreateAt = now
			svc.posting.Browse.UpdateAt = now
			// 已有浏览记录 更新用户浏览的时间
			if err := svc.posting.UpdateUserBrowsePost(userId, consts.TYPE_POST, post.Id); err != nil {
				log.Log.Errorf("post_trace: update user browse post err:%s", err)
			}
		} else {
			svc.posting.Browse.CreateAt = now
			svc.posting.Browse.UpdateAt = now
			svc.posting.Browse.UserId = userId
			svc.posting.Browse.ComposeId = post.Id
			svc.posting.Browse.ComposeType = consts.TYPE_POST
			// 添加用户浏览的帖子记录
			if err := svc.posting.RecordUserBrowsePost(); err != nil {
				log.Log.Errorf("post_trace: record user browse post err:%s", err)
			}
		}

		// 是否关注
		if attentionInfo := svc.attention.GetAttentionInfo(userId, post.UserId); attentionInfo != nil {
			resp.IsAttention = attentionInfo.Status
		}

		// 是否点赞
		if likeInfo := svc.like.GetLikeInfo(userId, post.Id, consts.TYPE_POSTS); likeInfo != nil {
			resp.IsLike = likeInfo.Status
		}
	}
	// 获取视频相关统计数据
	info, err := svc.posting.GetPostStatistic(fmt.Sprint(post.Id))
	if err == nil && info != nil {
		resp.BrowseNum = info.BrowseNum
		resp.CommentNum = info.CommentNum
		resp.FabulousNum = info.FabulousNum
		resp.ShareNum = info.ShareNum
	}


	return resp, errdef.SUCCESS
}

// 获取用户发布的帖子列表
func (svc *PostingModule) GetPostPublishListByUser(userId string, page, size int) []*mposting.PostDetailInfo {
	// 查询用户是否存在
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("post_trace: user not found, userId:%s", userId)
		return []*mposting.PostDetailInfo{}
	}

	offset := (page - 1) * size
	// 获取用户发布的帖子列表
	list, err := svc.posting.GetPublishPostByUser(userId, offset, size)
	if err != nil {
		log.Log.Errorf("post_trace: get publish post by user fail, userId:%s", userId)
		return []*mposting.PostDetailInfo{}
	}

	if len(list) == 0 {
		return []*mposting.PostDetailInfo{}
	}

	for _, item := range list {
		//item.Topics, err = svc.posting.GetPostTopic(fmt.Sprint(item.Id))
		//if item.Topics == nil || err != nil  {
		//	log.Log.Errorf("post_trace: get post topic fail, err:%s, item.Topics:%v", err, item.Topics)
		//	item.Topics = []*models.PostingTopic{}
		//}


		item.Avatar = user.Avatar
		item.Nickname = user.NickName


		// 如果是转发的视频数据
		if item.ContentType == consts.COMMUNITY_FORWARD_VIDEO {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardVideo); err != nil {
				log.Log.Errorf("community_trace: get forward video info err:%s", err)
				//return errdef.COMMUNITY_POSTS_BY_SECTION, []*mposting.PostDetailInfo{}
			} else {
				item.ForwardVideo.VideoAddr = svc.video.AntiStealingLink(item.ForwardVideo.VideoAddr)
			}

		}

		// 如果是转发的帖子
		if item.PostingType == consts.POST_TYPE_TEXT && item.ContentType == consts.COMMUNITY_FORWARD_POST {
			if err = util.JsonFast.UnmarshalFromString(item.Content, &item.ForwardPost); err != nil {
				log.Log.Errorf("community_trace: get forward post info err:%s", err)
				//return errdef.COMMUNITY_POSTS_BY_SECTION, []*mposting.PostDetailInfo{}
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
				//return errdef.COMMUNITY_POSTS_BY_SECTION, []*mposting.PostDetailInfo{}
			}
		}

		item.Content = ""
	}

	return list
}

// 用户删除发布的帖子
func (svc *PostingModule) DeletePublishPost(userId, postId string) int {
	// 开启事务
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("post_trace: session begin err:%s", err)
		return errdef.ERROR
	}

	// 查询用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("post_trace: user not found, userId:%s", userId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	// 查询帖子信息
	post, err := svc.posting.GetPostById(postId)
	if post == nil || err != nil {
		log.Log.Errorf("post_trace: post not found, postId:%d, err:%s", postId, err)
		svc.engine.Rollback()
		return errdef.POST_NOT_EXISTS
	}


	// 物理删除发布的帖子、帖子所属话题、帖子统计数据
	if err := svc.posting.DelPublishPostById(postId); err != nil {
		log.Log.Errorf("post_trace: delete publish post by id err:%s", err)
		svc.engine.Rollback()
		return errdef.POST_DELETE_PUBLISH_FAIL
	}

	// 删除帖子所属话题
	if err := svc.posting.DelPostTopics(postId); err != nil {
		log.Log.Errorf("post_trace: delete post topics err:%s", err)
		svc.engine.Rollback()
		return errdef.POST_DELETE_TOPIC_FAIL
	}

	// 删除帖子统计数据
	if err := svc.posting.DelPostStatistic(postId); err != nil {
		log.Log.Errorf("post_trace: delete post statistic err:%s", err)
		svc.engine.Rollback()
		return errdef.POST_DELETE_STATISTIC_FAIL
	}

	svc.engine.Commit()
	return errdef.SUCCESS
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

// 用户申请精华帖
func (svc *PostingModule) ApplyPostCream(userId string, param *mposting.ApplyCreamParam) int {
	if userId == "" {
		return errdef.USER_NO_LOGIN
	}

	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		return errdef.USER_NOT_EXISTS
	}

	post, err := svc.posting.GetPostById(fmt.Sprint(param.PostId))
	if err != nil || post == nil {
		log.Log.Errorf("post_trace: post not found, postId:%d", param.PostId)
		return errdef.POST_NOT_EXISTS
	}

	if post.UserId != userId {
		log.Log.Errorf("post_trace: userId not match, post.UserId:%s, userId:%s", post.UserId, userId)
		return errdef.POST_AUTHOR_NOT_MATCH
	}

	record, err := svc.posting.GetApplyCreamRecord(fmt.Sprint(param.PostId))
	if err != nil {
		log.Log.Errorf("post_trace: get apply cream record fail, err:%s, postId:%d", err, param.PostId)
		return errdef.POST_APPLY_CREAM_FAIL
	}

	if record != nil {
		log.Log.Errorf("post_trace: apply already exists, postId:%d", param.PostId)
		return errdef.POST_APPLY_ALREADY_EXISTS
	}

	now := int(time.Now().Unix())
	svc.posting.ApplyCream.PostId = param.PostId
	svc.posting.ApplyCream.UserId = userId
	svc.posting.ApplyCream.CreateAt = now
	svc.posting.ApplyCream.UpdateAt = now
	if _, err := svc.posting.AddApplyCreamRecord(); err != nil {
		log.Log.Errorf("post_trace: add apply cream record fail, err:%s", err)
		return errdef.POST_APPLY_CREAM_FAIL
	}

	return errdef.SUCCESS
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
