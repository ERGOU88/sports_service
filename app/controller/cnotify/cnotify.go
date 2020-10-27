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
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mlike"
	"sports_service/server/models/mnotify"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
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
}

func New(c *gin.Context) NotifyModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return NotifyModule{
		context: c,
		notify: mnotify.NewNotifyModel(socket),
		like: mlike.NewLikeModel(socket),
		collect: mcollect.NewCollectModel(socket),
		video: mvideo.NewVideoModel(socket),
		user: muser.NewUserModel(socket),
		comment: mcomment.NewCommentModel(socket),
		engine: socket,
	}
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

	res := make([]interface{}, len(list))
	for index, liked := range list {
		switch liked.ZanType {
		// 被点赞的视频
		case consts.TYPE_VIDEOS:
			info := new(mlike.BeLikedVideoInfo)
			info.OpTime = liked.CreateAt
			info.Type = consts.TYPE_VIDEOS
			// 视频作品
			if video := svc.video.FindVideoById(fmt.Sprint(liked.TypeId)); video != nil {
				info.ComposeId = video.VideoId
				info.Title = video.Title
				info.Describe = video.Describe
				info.Cover = video.Cover
				info.VideoAddr = video.VideoAddr
				info.VideoDuration = video.VideoDuration
				info.VideoWidth = video.VideoWidth
				info.VideoHeight = video.VideoHeight
				info.CreateAt = video.CreateAt
				info.Status = int32(video.Status)
			}

			// 视频统计数据
			if statistic := svc.video.GetVideoStatistic(fmt.Sprint(liked.TypeId)); statistic != nil {
				info.BarrageNum = statistic.BarrageNum
				info.BrowseNum = statistic.BrowseNum
			}

			// 点赞用户
			if user := svc.user.FindUserByUserid(liked.UserId); user != nil {
				info.UserId = user.UserId
				info.Avatar = user.Avatar
				info.Nickname = user.NickName
			}

			// 被点赞用户
			if toUser := svc.user.FindUserByUserid(userId); toUser != nil {
				info.ToUserId = toUser.UserId
				info.ToUserAvatar = toUser.Avatar
				info.ToUserName = toUser.NickName
			}

			res[index] = info


		// 被点赞的帖子
		case consts.TYPE_POSTS:
		// 被点赞的评论
		case consts.TYPE_COMMENT:
			info := new(mlike.BeLikedCommentInfo)
			info.OpTime = liked.CreateAt
			info.Type = consts.TYPE_COMMENT
			// 获取评论信息
			comment := svc.comment.GetVideoCommentById(fmt.Sprint(liked.TypeId))
			if comment != nil {
				// 被点赞的信息
				info.Content = comment.Content
				info.ToUserId = comment.UserId
				info.ToUserAvatar = comment.Avatar
				info.ToUserName = comment.UserName

				// 获取评论对应的视频信息
				if video := svc.video.FindVideoById(fmt.Sprint(comment.VideoId)); video != nil {
					info.ComposeId = video.VideoId
					info.Title = video.Title
					info.Describe = video.Describe
					info.Cover = video.Cover
					info.VideoAddr = video.VideoAddr
					info.VideoDuration = video.VideoDuration
					info.VideoWidth = video.VideoWidth
					info.VideoHeight = video.VideoHeight
					info.CreateAt = video.CreateAt
					info.Status = int32(video.Status)
				}

				// 视频统计数据
				if statistic := svc.video.GetVideoStatistic(fmt.Sprint(liked.TypeId)); statistic != nil {
					info.BarrageNum = statistic.BarrageNum
					info.BrowseNum = statistic.BrowseNum
				}

				// 点赞用户
				if user := svc.user.FindUserByUserid(liked.UserId); user != nil {
					info.UserId = user.UserId
					info.Avatar = user.Avatar
					info.Nickname = user.NickName
				}

				res[index] = info

			}

		}
	}

	// 记录读取被点赞通知消息的时间
	if err := svc.notify.RecordReadBeLikedTime(userId); err != nil {
		log.Log.Errorf("notify_trace: record read liked notify time err:%s", err)
	}

	return res
}

// 获取用户 @ 通知
func (svc *NotifyModule) GetReceiveAtNotify(userId string, page, size int) []interface{} {
	if userId == "" {
		log.Log.Error("notify_trace: need login")
		return []interface{}{}
	}

	if info := svc.user.FindUserByUserid(userId); info == nil {
		log.Log.Errorf("notify_trace: user not found, userId:%s", userId)
		return []interface{}{}
	}

	offset := (page - 1) * size
	// 被@的评论列表
	list := svc.comment.GetReceiveAtList(userId, offset, size)
	if len(list) == 0 {
		log.Log.Error("notify_trace: receive at list empty")
		return []interface{}{}
	}

	res := make([]interface{}, len(list))
	for index, receiveAt := range list {
		switch receiveAt.TopicType {
		case consts.TYPE_POSTS:

		case consts.TYPE_COMMENT:
			info := new(mnotify.ReceiveCommentAtInfo)
			info.AtTime = receiveAt.CreateAt
			info.Type = consts.TYPE_COMMENT
			// 获取评论信息
			comment := svc.comment.GetVideoCommentById(fmt.Sprint(receiveAt.CommentId))
			if comment != nil {
				// 执行@的用户信息
				info.Reply = comment.Content
				info.UserId = comment.UserId
				info.Avatar = comment.Avatar
				info.Nickname = comment.UserName

				// 获取评论对应的视频信息
				if video := svc.video.FindVideoById(fmt.Sprint(comment.VideoId)); video != nil {
					info.ComposeId = video.VideoId
					info.Title = video.Title
					info.Describe = video.Describe
					info.Cover = video.Cover
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

				info.CommentType = 1
				// 如果父评论id为0 则表示 是1级评论 不为0 则表示是回复
				if comment.ParentCommentId != 0 {
					// 获取父级回复
					if parent := svc.comment.GetVideoCommentById(fmt.Sprint(comment.ReplyCommentId)); parent != nil {
						info.CommentType = 2
						info.Content = parent.Content
					}
				}

				res[index] = info

			}
		}
	}

	// 记录读取@通知消息的时间
	if err := svc.notify.RecordReadAtTime(userId); err != nil {
		log.Log.Errorf("notify_trace: record read at notify time err:%s", err)
	}

	return res
}

// 获取未读消息数量
func (svc *NotifyModule) GetUnreadNum(userId string) map[string]int64 {
	mp := make(map[string]int64, 0)
	mp["sys_notify"] = 0
	mp["liked_notify"] = 0
	mp["comment_notify"] = 0

	if userId == "" {
		log.Log.Error("notify_trace: need login")
		return mp
	}

	if info := svc.user.FindUserByUserid(userId); info == nil {
		log.Log.Errorf("notify_trace: user not found, userId:%s", userId)
		return mp
	}

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

