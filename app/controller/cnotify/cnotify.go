package cnotify

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mlike"
	"sports_service/server/models/mnotify"
	"sports_service/server/models/mposting"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"sports_service/server/util"
	"strconv"
	"strings"
	"time"
)

type NotifyModule struct {
	context    *gin.Context
	engine     *xorm.Session
	notify     *mnotify.NotifyModel
	like       *mlike.LikeModel
	collect    *mcollect.CollectModel
	video      *mvideo.VideoModel
	user       *muser.UserModel
	comment    *mcomment.CommentModel
	attention  *mattention.AttentionModel
	post       *mposting.PostingModel
}

func New(c *gin.Context) NotifyModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return NotifyModule{
		context: c,
		notify: mnotify.NewNotifyModel(socket),
		like: mlike.NewLikeModel(socket),
		collect: mcollect.NewCollectModel(socket),
		video: mvideo.NewVideoModel(socket),
		user: muser.NewUserModel(socket),
		comment: mcomment.NewCommentModel(socket),
		attention: mattention.NewAttentionModel(socket),
		post: mposting.NewPostingModel(socket),
		engine: socket,
	}
}

// 所有数据 同一视频/评论点赞 整合为一条数据
func (svc *NotifyModule) GetNewBeLikedList(userId string, page, size int) []interface{} {
	if userId == "" {
		log.Log.Error("notify_trace: need login")
		return []interface{}{}
	}

	if info := svc.user.FindUserByUserid(userId); info == nil {
		log.Log.Errorf("notify_trace: user not found, userId:%s", userId)
		return []interface{}{}
	}

	offset := (page - 1) * size
	// 被点赞的作品列表
	list := svc.like.GetNewBeLikedList(userId, offset, size)
	if len(list) == 0 {
		log.Log.Error("notify_trace: be liked list empty")
		return []interface{}{}
	}

	res := make([]interface{}, len(list))
	log.Log.Debugf("notify_trace: length:%d", len(res))
	for index, liked := range list {
		info := new(mlike.BeLikedInfo)
		info.OpTime = liked.CreateAt
		switch liked.ZanType {
		// 被点赞的视频
		case consts.TYPE_VIDEOS:
			info.Type = consts.TYPE_VIDEOS
			info.JumpVideoId = liked.TypeId
			video := svc.video.FindVideoById(fmt.Sprint(liked.TypeId))
			if video != nil {
				info.ComposeId = video.VideoId
				info.Describe = util.TrimHtml(video.Describe)
				info.Title = util.TrimHtml(video.Title)
				//info.Title = video.Title
				//info.Describe = video.Describe
				info.Cover = video.Cover
				info.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
				info.VideoDuration = video.VideoDuration
				info.VideoWidth = video.VideoWidth
				info.VideoHeight = video.VideoHeight
				info.CreateAt = video.CreateAt
				// 视频统计数据
				if statistic := svc.video.GetVideoStatistic(fmt.Sprint(liked.TypeId)); statistic != nil {
					info.BarrageNum = statistic.BarrageNum
					info.BrowseNum = statistic.BrowseNum
				}
			}


		// 被点赞的帖子
		case consts.TYPE_POSTS:
			info.Type = consts.TYPE_POSTS
			info.JumpPostId = liked.TypeId

			post, err := svc.post.GetPostById(fmt.Sprint(liked.TypeId))
			if post != nil && err == nil {
				info.ComposeId = post.Id
				info.Title = post.Title
				info.Describe = post.Describe
				info.CreateAt = post.CreateAt
				// 图文帖
				if post.PostingType == consts.POST_TYPE_IMAGE {
					var images []string
					if err = util.JsonFast.UnmarshalFromString(post.Content, &images); err != nil {
						log.Log.Errorf("post_trace: get image info err:%s", err)
					}

					// 使用第一张图片展示
					if len(images) > 0 {
						info.Cover = images[0]
					}
				}
			}


		// 被点赞的视频评论
		case consts.TYPE_VIDEO_COMMENT:
			info.Type = consts.TYPE_VIDEO_COMMENT
			info.JumpVideoId = liked.TypeId

			// 获取评论信息
			comment := svc.comment.GetVideoCommentById(fmt.Sprint(liked.TypeId))
			if comment != nil {
				// 被点赞的信息
				info.Content = comment.Content
				info.ComposeId = comment.Id
				// 顶级评论id
				info.ParentCommentId = comment.ParentCommentId
				if info.ParentCommentId == 0 {
					// 当前评论即顶级评论
					info.ParentCommentId = comment.Id
				}

				// 获取评论对应的视频信息
				video := svc.video.FindVideoById(fmt.Sprint(comment.VideoId))
				if video != nil {
					//info.Title = video.Title
					//info.Describe = video.Describe
					info.Describe = util.TrimHtml(video.Describe)
					info.Title = util.TrimHtml(video.Title)
					info.Cover = video.Cover
					info.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
					info.VideoDuration = video.VideoDuration
					info.VideoWidth = video.VideoWidth
					info.VideoHeight = video.VideoHeight
					info.CreateAt = video.CreateAt
					// 视频统计数据
					if statistic := svc.video.GetVideoStatistic(fmt.Sprint(liked.TypeId)); statistic != nil {
						info.BarrageNum = statistic.BarrageNum
						info.BrowseNum = statistic.BrowseNum
					}
				}


			}

			case consts.TYPE_POST_COMMENT:
				info.Type = consts.TYPE_POST_COMMENT
				info.JumpPostId = liked.TypeId

				// 获取视频评论信息
				comment := svc.comment.GetPostCommentById(fmt.Sprint(liked.TypeId))
				if comment != nil {
					// 被点赞的信息
					info.Content = comment.Content
					info.ComposeId = comment.Id
					// 顶级评论id
					info.ParentCommentId = comment.ParentCommentId
					if info.ParentCommentId == 0 {
						// 当前评论即顶级评论
						info.ParentCommentId = comment.Id
					}

					post, err := svc.post.GetPostById(fmt.Sprint(liked.TypeId))
					if post != nil && err == nil {
						info.ComposeId = post.Id
						info.Title = post.Title
						info.Describe = post.Describe
						info.CreateAt = post.CreateAt
						// 图文帖
						if post.PostingType == consts.POST_TYPE_IMAGE {
							var images []string
							if err = util.JsonFast.UnmarshalFromString(post.Content, &images); err != nil {
								log.Log.Errorf("post_trace: get image info err:%s", err)
							}

							// 使用第一张图片展示
							if len(images) > 0 {
								info.Cover = images[0]
							}
						}
					}

				}
		}


		var userList []*models.User
		userIds := strings.Split(liked.UserId, ",")
		lenth := len(userIds)
		if lenth >= 2 {
			// 最多取两个 取最新
			userList = svc.user.FindUserByUserids(strings.Join(userIds[lenth-2:], ","), 0, 2)

		} else {
			userList = svc.user.FindUserByUserids(strings.Join(userIds, ""), 0, 1)
		}

		for _, user := range userList {
			info.UserList = append(info.UserList, &mlike.LikedUserInfo{
				UserId: user.UserId,
				NickName: user.NickName,
				Avatar: user.Avatar,
				OpTm:  liked.CreateAt,
			})
		}

		info.TotalLikeNum = lenth
		log.Log.Debugf("notify_trace: beLiked info:%+v", info)
		res[index] = info
	}

	// 记录读取被点赞通知消息的时间
	if err := svc.notify.RecordReadBeLikedTime(userId); err != nil {
		log.Log.Errorf("notify_trace: record read liked notify time err:%s", err)
	}

	return res
}

// 保存用户通知设置
func (svc *NotifyModule) SaveUserNotifySetting(userId string, params *mnotify.NotifySettingParams) int {
	if info := svc.user.FindUserByUserid(userId); info == nil {
		log.Log.Errorf("notify_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	svc.notify.NofitySetting.UserId = userId
	svc.notify.NofitySetting.AttentionPushSet = params.AttentionPushSet
	svc.notify.NofitySetting.CommentPushSet = params.CommentPushSet
	svc.notify.NofitySetting.SharePushSet = params.SharePushSet
	svc.notify.NofitySetting.SlotPushSet = params.SlotPushSet
	svc.notify.NofitySetting.ThumbUpPushSet = params.ThumbUpPushSet
	svc.notify.NofitySetting.UpdateAt = int(time.Now().Unix())
	// 更新用户设置
	if err := svc.notify.UpdateUserNotifySetting(); err != nil {
		log.Log.Errorf("notify_trace: update user notify setting err:%s", err)
		return errdef.NOTIFY_SETTING_FAIL
	}

	return errdef.SUCCESS
}

// todo: review
// 获取用户被点赞的作品列表 返回值：map中key表示类型 1 点赞视频 2 点赞帖子 3 点赞评论
func (svc *NotifyModule) GetBeLikedList(userId string, page, size int) []interface{} {
	if userId == "" {
		log.Log.Error("notify_trace: need login")
		return []interface{}{}
	}

	if info := svc.user.FindUserByUserid(userId); info == nil {
		log.Log.Errorf("notify_trace: user not found, userId:%s", userId)
		return []interface{}{}
	}

	offset := (page - 1) * size
	// 被点赞的作品列表
	list := svc.like.GetBeLikedList(userId, offset, size)
	if len(list) == 0 {
		log.Log.Error("notify_trace: be liked list empty")
		return []interface{}{}
	}

	videoMp := make(map[int64]*models.Videos)
	commentMp := make(map[int64]*models.VideoComment)

	// 视频点赞的用户列表
	vlikeMap := make(map[int64][]*mlike.LikedUserInfo, 0)
	// 评论点赞的用户列表
	clikeMap := make(map[int64][]*mlike.LikedUserInfo, 0)
	for _, liked := range list {
		switch liked.ZanType {
		// 被点赞的视频
		case consts.TYPE_VIDEOS:
			// 视频作品
			if video := svc.video.FindVideoById(fmt.Sprint(liked.TypeId)); video != nil {
				videoMp[video.VideoId] = video
				if user := svc.user.FindUserByUserid(liked.UserId); user != nil {
					log.Log.Errorf("user.Id:%s", user.UserId)
					//userMp[fmt.Sprintf("%s_%d", user.UserId, liked.TypeId)] = user

					vlikeMap[video.VideoId] = append(vlikeMap[video.VideoId], &mlike.LikedUserInfo{
						UserId: user.UserId,
						Avatar: user.Avatar,
						NickName: user.NickName,
						OpTm: liked.CreateAt,
					})

				}
			}
		// 被点赞的帖子
		case consts.TYPE_POSTS:
		// 被点赞的评论
		case consts.TYPE_VIDEO_COMMENT:
			// 获取评论信息
			comment := svc.comment.GetVideoCommentById(fmt.Sprint(liked.TypeId))
			if comment != nil {
				commentMp[comment.Id] = comment
				// 点赞用户
				if user := svc.user.FindUserByUserid(liked.UserId); user != nil {
					clikeMap[comment.Id] = append(clikeMap[comment.Id], &mlike.LikedUserInfo{
						UserId: user.UserId,
						Avatar: user.Avatar,
						NickName: user.NickName,
						OpTm: liked.CreateAt,
					})

				}
			}
		}
	}

	//var lastRead int
	// 获取用户上次读取被点赞通知消息的时间
	//readTm, err := svc.notify.GetReadBeLikedTime(userId)
	//if err == nil {
	//  tm, err := strconv.Atoi(readTm)
	//  if err != nil {
	//    log.Log.Errorf("notify_trace: strconv atoi err:%s", err)
	//  }
	//
	//  lastRead = tm
	//}

	// 上次已读取的数据下标（默认-1未读取）
	//var readIndex int32 = -1
	res := make([]interface{}, 0)
	log.Log.Debugf("notify_trace: length:%d", len(res))
	for _, liked := range list {
		switch liked.ZanType {
		// 被点赞的视频
		case consts.TYPE_VIDEOS:
			info := new(mlike.BeLikedInfo)
			info.OpTime = liked.CreateAt
			info.Type = consts.TYPE_VIDEOS
			info.JumpVideoId = liked.TypeId
			// 视频作品
			video, ok := videoMp[liked.TypeId]
			if ok && video != nil {
				videoMp[liked.TypeId] = nil
				info.ComposeId = video.VideoId
				info.Title = video.Title
				info.Describe = util.TrimHtml(video.Describe)
				info.Cover = util.TrimHtml(video.Cover)
				info.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
				info.VideoDuration = video.VideoDuration
				info.VideoWidth = video.VideoWidth
				info.VideoHeight = video.VideoHeight
				info.CreateAt = video.CreateAt
				// 视频统计数据
				if statistic := svc.video.GetVideoStatistic(fmt.Sprint(liked.TypeId)); statistic != nil {
					info.BarrageNum = statistic.BarrageNum
					info.BrowseNum = statistic.BrowseNum
				}

				users, ok := vlikeMap[liked.TypeId]
				if ok && users != nil {
					lenth := len(users)
					if lenth >= 2 {
						// 最多取两个 取最新
						info.UserList = users[lenth - 2:]
						// 返回最新的点赞时间
						info.OpTime = users[0].OpTm
					} else {
						info.UserList = users
					}

					info.TotalLikeNum = lenth
				}

				res = append(res, info)
			}
			// 被点赞的帖子
		case consts.TYPE_POSTS:
		// 被点赞的评论
		case consts.TYPE_VIDEO_COMMENT:
			info := new(mlike.BeLikedInfo)
			info.OpTime = liked.CreateAt
			info.Type = consts.TYPE_VIDEO_COMMENT

			// 获取评论信息
			comment, ok := commentMp[liked.TypeId]
			if ok && comment != nil {
				info.JumpVideoId = comment.VideoId
				// 评论/视频点赞 多条记录整合为一条
				commentMp[liked.TypeId] = nil
				// 被点赞的信息
				info.Content = comment.Content
				info.ComposeId = comment.Id
				// 顶级评论id
				info.ParentCommentId = comment.ParentCommentId
				if info.ParentCommentId == 0 {
					// 当前评论即顶级评论
					info.ParentCommentId = comment.Id
				}

				// 获取评论对应的视频信息
				video, ok := videoMp[comment.VideoId]
				log.Log.Debugf("notify_trace: get video by comment id, videoId:%d", comment.VideoId)
				if ok && video != nil {
					info.Title = video.Title
					info.Describe = util.TrimHtml(video.Describe)
					info.Cover = util.TrimHtml(video.Cover)
					info.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
					info.VideoDuration = video.VideoDuration
					info.VideoWidth = video.VideoWidth
					info.VideoHeight = video.VideoHeight
					info.CreateAt = video.CreateAt
					// 视频统计数据
					if statistic := svc.video.GetVideoStatistic(fmt.Sprint(liked.TypeId)); statistic != nil {
						info.BarrageNum = statistic.BarrageNum
						info.BrowseNum = statistic.BrowseNum
					}
				}

				users, ok := clikeMap[liked.TypeId]
				if ok && users != nil {
					lenth := len(users)
					if lenth >= 2 {
						// 最多取两个 取最新
						info.UserList = users[lenth - 2:]
						// 返回最新的点赞时间
						info.OpTime = users[0].OpTm
					} else {
						info.UserList = users
					}

					info.TotalLikeNum = lenth
				}

				res = append(res, info)
			}
		}

		// 未记录读取的下标
		// 用户上次读取的数据下标
		//if lastRead < liked.UpdateAt {
		//  readIndex = int32(index)
		//}
	}

	// 记录读取被点赞通知消息的时间
	if err := svc.notify.RecordReadBeLikedTime(userId); err != nil {
		log.Log.Errorf("notify_trace: record read liked notify time err:%s", err)
	}

	return res
}

// 获取用户 @ 通知
func (svc *NotifyModule) GetReceiveAtNotify(userId string, page, size int) ([]interface{}, int) {
	if userId == "" {
		log.Log.Error("notify_trace: need login")
		return []interface{}{}, -1
	}

	if info := svc.user.FindUserByUserid(userId); info == nil {
		log.Log.Errorf("notify_trace: user not found, userId:%s", userId)
		return []interface{}{}, -1
	}

	offset := (page - 1) * size
	// 被@的评论列表
	list := svc.comment.GetReceiveAtList(userId, offset, size)
	if len(list) == 0 {
		log.Log.Error("notify_trace: receive at list empty")
		return []interface{}{}, -1
	}

	// 用户上次读取被@列表数据的时间
	var lastRead int
	// 获取用户上次读取被@通知消息的时间
	readAt, err := svc.notify.GetReadAtTime(userId)
	if err == nil {
		tm, err := strconv.Atoi(readAt)
		if err != nil {
			log.Log.Errorf("notify_trace: strconv atoi err:%s", err)
		}

		lastRead = tm

	}

	// 上次已读取的数据下标（默认-1未读取）
	var readIndex = -1
	res := make([]interface{}, len(list))
	for index, receiveAt := range list {
		switch receiveAt.TopicType {
		// 视频评论里@ 则展示视频内容
		case consts.TYPE_VIDEOS:
			info := new(mnotify.ReceiveAtInfo)
			info.AtTime = receiveAt.UpdateAt
			info.Type = consts.TYPE_VIDEOS
			// 获取评论信息
			comment := svc.comment.GetVideoCommentById(fmt.Sprint(receiveAt.ComposeId))
			if comment != nil {
				video := svc.video.FindVideoById(fmt.Sprint(comment.VideoId))
				if video != nil {
					info.ComposeId = video.VideoId
					info.Title = video.Title
					info.Describe = util.TrimHtml(video.Describe)
					info.Cover = util.TrimHtml(video.Cover)
					info.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
					info.VideoDuration = video.VideoDuration
					info.VideoWidth = video.VideoWidth
					info.VideoHeight = video.VideoHeight
					info.CreateAt = video.CreateAt
					// 视频统计数据
					if statistic := svc.video.GetVideoStatistic(fmt.Sprint(receiveAt.ComposeId)); statistic != nil {
						info.BarrageNum = statistic.BarrageNum
						info.BrowseNum = statistic.BrowseNum
					}
				}
			}
		// 帖子评论里@ 或 发布帖子时 内容@用户 则展示帖子内容 todo: 展示待确认
		case consts.TYPE_POSTS, consts.TYPE_PUBLISH_POST:
			info := new(mnotify.ReceiveAtInfo)
			info.AtTime = receiveAt.UpdateAt
			info.Type = receiveAt.TopicType
			comment := svc.comment.GetPostCommentById(fmt.Sprint(receiveAt.ComposeId))
			if comment != nil {
				post, err := svc.post.GetPostById(fmt.Sprint(comment.PostId))
				if post != nil && err == nil {
					info.ComposeId = post.Id
					info.Title = post.Title
					info.Describe = post.Describe
					info.CreateAt = post.CreateAt

					// 图文帖
					if post.PostingType == consts.POST_TYPE_IMAGE {
						var images []string
						if err = util.JsonFast.UnmarshalFromString(post.Content, &images); err != nil {
							log.Log.Errorf("post_trace: get image info err:%s", err)
						}

						// 使用第一张图片展示
						if len(images) > 0 {
							info.Cover = images[0]
						}
					}

				}

				// 执行@的用户信息
				if user := svc.user.FindUserByUserid(receiveAt.UserId); user != nil {
					info.UserId = user.UserId
					info.Avatar = user.Avatar
					info.Nickname = user.NickName
				}

				res[index] = info
			}

		// 视频直接评论/回复
		case consts.TYPE_VIDEO_COMMENT:
			info := new(mnotify.ReceiveAtInfo)
			info.AtTime = receiveAt.UpdateAt
			info.Type = consts.TYPE_VIDEO_COMMENT
			// 获取评论信息
			comment := svc.comment.GetVideoCommentById(fmt.Sprint(receiveAt.ComposeId))
			if comment != nil {
				// 执行@的用户信息
				if user := svc.user.FindUserByUserid(receiveAt.UserId); user != nil {
					// 执行@的用户信息
					info.UserId = user.UserId
					info.Avatar = user.Avatar
					info.Nickname = user.NickName
				}

				// 获取评论对应的视频信息
				if video := svc.video.FindVideoById(fmt.Sprint(comment.VideoId)); video != nil {
					info.ComposeId = video.VideoId
					info.Title = video.Title
					info.Describe = util.TrimHtml(video.Describe)
					info.Cover = util.TrimHtml(video.Cover)
					info.VideoAddr = video.VideoAddr
					info.VideoDuration = video.VideoDuration
					info.VideoWidth = video.VideoWidth
					info.VideoHeight = video.VideoHeight
					info.CreateAt = video.CreateAt
					info.Status = int32(video.Status)
				}

				// 视频统计数据
				if statistic := svc.video.GetVideoStatistic(fmt.Sprint(comment.VideoId)); statistic != nil {
					info.BarrageNum = statistic.BarrageNum
					info.BrowseNum = statistic.BrowseNum
				}

				// 被@的用户信息
				if user := svc.user.FindUserByUserid(receiveAt.ToUserId); user != nil {
					info.ToUserId = user.UserId
					info.ToUserAvatar = user.Avatar
					info.ToUserName = user.NickName
				}

				// 默认1级评论
				info.CommentType = 1
				// 评论id（1级评论id）
				info.CommentId = receiveAt.ComposeId
				// 进行回复使用的id
				info.ReplyCommentId = receiveAt.ComposeId
				info.Content = comment.Content
				// 获取当前评论 / 回复的被点赞数
				info.TotalLikeNum = svc.like.GetLikeNumByType(receiveAt.ComposeId, consts.TYPE_VIDEO_COMMENT)
				// 如果父评论id为0 则表示 是1级评论 不为0 则表示是回复
				if comment.ParentCommentId != 0 {
					// 获取被回复的内容
					beReply := svc.comment.GetVideoCommentById(fmt.Sprint(comment.ReplyCommentId))
					if beReply != nil {
						info.CommentType = 2
						info.Content = beReply.Content
						info.Reply = comment.Content

						// 被回复的不是1级评论 则@消息 为1
						if beReply.CommentLevel != 1 {
							info.IsAt = 1
						} else {
							// 1级评论id
							info.CommentId = beReply.Id
						}
					}

					// 获取最上级的评论内容
					parent := svc.comment.GetVideoCommentById(fmt.Sprint(comment.ParentCommentId))
					if parent != nil {
						info.ParentComment = parent.Content
						// 1级评论id
						info.CommentId = parent.Id
					}

				}

				if userId != "" {
					// 获取点赞的信息
					if likeInfo := svc.like.GetLikeInfo(userId, comment.Id, consts.TYPE_VIDEO_COMMENT); likeInfo != nil {
						info.IsLike = likeInfo.Status
					}
				}

				res[index] = info

			}

		// 帖子直接评论/回复
		case consts.TYPE_POST_COMMENT:
			info := new(mnotify.ReceiveAtInfo)
			info.AtTime = receiveAt.UpdateAt
			info.Type = consts.TYPE_VIDEO_COMMENT
			// 获取评论信息
			comment := svc.comment.GetPostCommentById(fmt.Sprint(receiveAt.ComposeId))
			if comment != nil {
				// 执行@的用户信息
				if user := svc.user.FindUserByUserid(receiveAt.UserId); user != nil {
					// 执行@的用户信息
					info.UserId = user.UserId
					info.Avatar = user.Avatar
					info.Nickname = user.NickName
				}

				// 获取评论对应的帖子信息
				post, err := svc.post.GetPostById(fmt.Sprint(receiveAt.ComposeId))
				if post != nil && err == nil {
					info.ComposeId = post.Id
					info.Title = post.Title
					info.Describe = post.Describe
					info.CreateAt = post.CreateAt
					// 图文帖
					if post.PostingType == consts.POST_TYPE_IMAGE {
						var images []string
						if err = util.JsonFast.UnmarshalFromString(post.Content, &images); err != nil {
							log.Log.Errorf("post_trace: get image info err:%s", err)
						}

						// 使用第一张图片展示
						if len(images) > 0 {
							info.Cover = images[0]
						}
					}
				}

				// 被@的用户信息
				if user := svc.user.FindUserByUserid(receiveAt.ToUserId); user != nil {
					info.ToUserId = user.UserId
					info.ToUserAvatar = user.Avatar
					info.ToUserName = user.NickName
				}

				// 默认1级评论
				info.CommentType = 1
				// 评论id（1级评论id）
				info.CommentId = receiveAt.ComposeId
				// 进行回复使用的id
				info.ReplyCommentId = receiveAt.ComposeId
				info.Content = comment.Content
				// 获取当前评论 / 回复的被点赞数
				info.TotalLikeNum = svc.like.GetLikeNumByType(receiveAt.ComposeId, consts.TYPE_POST_COMMENT)
				// 如果父评论id为0 则表示 是1级评论 不为0 则表示是回复
				if comment.ParentCommentId != 0 {
					// 获取被回复的内容
					beReply := svc.comment.GetPostCommentById(fmt.Sprint(comment.ReplyCommentId))
					if beReply != nil {
						info.CommentType = 2
						info.Content = beReply.Content
						info.Reply = comment.Content

						// 被回复的不是1级评论 则@消息 为1
						if beReply.CommentLevel != 1 {
							info.IsAt = 1
						} else {
							// 1级评论id
							info.CommentId = beReply.Id
						}
					}

					// 获取最上级的评论内容
					parent := svc.comment.GetPostCommentById(fmt.Sprint(comment.ParentCommentId))
					if parent != nil {
						info.ParentComment = parent.Content
						// 1级评论id
						info.CommentId = parent.Id
					}

				}

				if userId != "" {
					// 获取点赞的信息
					if likeInfo := svc.like.GetLikeInfo(userId, comment.Id, consts.TYPE_POST_COMMENT); likeInfo != nil {
						info.IsLike = likeInfo.Status
					}
				}

				res[index] = info
			}
		}

		log.Log.Debugf("receiveId: %d", receiveAt.Id)
		if lastRead < receiveAt.UpdateAt {
			// 用户上次读取的数据下标
			readIndex = index
		}
	}

	// 记录读取@通知消息的时间
	if err := svc.notify.RecordReadAtTime(userId); err != nil {
		log.Log.Errorf("notify_trace: record read at notify time err:%s", err)
	}

	return res, readIndex
}

// 获取未读消息总数 及 未浏览视频数[关注用户发布的视频]（首页展示）
func (svc *NotifyModule) GetUnreadTotalNum(userId string) *mnotify.HomePageNotify {
	resp := &mnotify.HomePageNotify{
		UnBrowsedNum: 0,
		UnreadNum: 0,
	}

	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("notify_trace: user not found, userId:%s", userId)
		return resp
	}

	// 获取未读的系统消息数
	resp.UnreadNum = svc.notify.GetUnreadSystemMsgNum(userId)
	// 获取用户上次读取被点赞列表的时间
	readTm, err := svc.notify.GetReadBeLikedTime(userId)
	if err == nil || err == redis.ErrNil {
		if readTm == "" {
			readTm = "0"
		}
		// 获取未读的被点赞的数量
		resp.UnreadNum += svc.like.GetUnreadBeLikedCount(userId, readTm)
	}

	// 获取用户上次读取被@列表数据的时间
	readAt, err := svc.notify.GetReadAtTime(userId)
	if err == nil || err == redis.ErrNil {
		if readAt == "" {
			readAt = "0"
		}
		// 获取未读的被@的数量
		resp.UnreadNum += svc.comment.GetUnreadAtCount(userId, readAt)
	}

	// 用户上次浏览时间（关注用户发布的视频列表 ）
	browseTm, err := svc.notify.GetReadAttentionPubVideo(userId)
	if err == nil || err == redis.ErrNil {
		if browseTm == "" {
			browseTm = "0"
		}

		userIds := svc.attention.GetAttentionList(userId)
		if len(userIds) == 0 {
			log.Log.Errorf("video_trace: not following any users")
			return resp
		}

		uids := strings.Join(userIds, ",")
		resp.UnBrowsedNum = svc.video.GetUnBrowsedAttentionVideos(uids, browseTm)
	}

	return resp
}

// 获取未读消息数量
func (svc *NotifyModule) GetUnreadNum(userId string) map[string]int64 {
	mp := make(map[string]int64, 0)
	mp["sys_notify"] = 0
	mp["liked_notify"] = 0
	mp["comment_notify"] = 0
	mp["at_notify"] = 0

	if userId == "" {
		log.Log.Errorf("notify_trace: user need login")
		return mp
	}

	if info := svc.user.FindUserByUserid(userId); info == nil {
		log.Log.Errorf("notify_trace: user not found, userId:%s", userId)
		return mp
	}

	mp["sys_notify"] = svc.notify.GetUnreadSystemMsgNum(userId)

	// 获取用户上次读取被点赞列表的时间
	readTm, err := svc.notify.GetReadBeLikedTime(userId)
	if err == nil || err == redis.ErrNil {
		if readTm == "" {
			readTm = "0"
		}
		// 获取未读的被点赞的数量
		mp["liked_notify"] = svc.like.GetUnreadBeLikedCount(userId, readTm)
	}

	// 获取用户上次读取被@列表数据的时间
	readAt, err := svc.notify.GetReadAtTime(userId)
	log.Log.Errorf("comment_trace: readAt:%s, err:%s", readAt, err)
	if err == nil || err == redis.ErrNil {
		if readAt == "" {
			readAt = "0"
		}
		// 获取未读的被@的数量
		mp["at_notify"] = svc.comment.GetUnreadAtCount(userId, readAt)
	}

	return mp
}

// 获取用户通知设置
func (svc *NotifyModule) GetUserNotifySetting(userId string) *models.SystemNoticeSettings {
	if userId == "" {
		log.Log.Error("notify_trace: need login")
		return nil
	}

	if info := svc.user.FindUserByUserid(userId); info == nil {
		log.Log.Errorf("notify_trace: user not found, userId:%s", userId)
		return nil
	}

	return svc.notify.GetUserNotifySetting(userId)
}

// 获取系统通知列表
func (svc *NotifyModule) GetSystemNotify(userId string, page, size int) []*models.SystemMessage {
	offset := (page - 1) * size
	list := svc.notify.GetSystemNotify(userId, offset, size)
	if list == nil {
		return []*models.SystemMessage{}
	}

	if userId == "" {
		return list
	}

	// 最后读取的消息id
	var lastReadId int64
	for _, msg := range list {
		// 0 未读 1 已读
		if msg.Status == 0 && lastReadId == 0 {
			lastReadId = msg.SystemId
		}

		// 剔除h5标签
		msg.SystemContent = util.TrimHtml(msg.SystemContent)
	}

	// 状态更新为已读
	svc.UpdateNotifyStatus(lastReadId)

	return list
}

// 状态更新为已读
func (svc *NotifyModule) UpdateNotifyStatus(lastReadId int64) {
	go func() {
		if err := svc.notify.UpdateSystemNotifyStatus(lastReadId); err != nil {
			log.Log.Errorf("notify_trace: update system notify status err:%s", err)
		}
	}()
}

// 获取系统消息详情
func (svc *NotifyModule) GetSystemMsgById(systemId string) *models.SystemMessage {
	return svc.notify.GetSystemNotifyById(systemId)
}
