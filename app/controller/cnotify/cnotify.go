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

  mp := make(map[string]string)
  videoMp := make(map[int64]*models.Videos)
  commentMp := make(map[int64]*models.VideoComment)
  userMp := make(map[string]*models.User)
  for _, liked := range list {
    switch liked.ZanType {
    // 被点赞的视频
    case consts.TYPE_VIDEOS:
      // 视频作品
      if video := svc.video.FindVideoById(fmt.Sprint(liked.TypeId)); video != nil {
        videoMp[video.VideoId] = video
        if user := svc.user.FindUserByUserid(liked.UserId); user != nil {
          log.Log.Errorf("user.Id:%s", user.UserId)
          userMp[user.UserId] = user
          nickNames, ok := mp[fmt.Sprintf("%d_%d", video.VideoId, liked.ZanType)]
          // 如果点赞的是同一视频  昵称整合为一条数据
          if ok {
            if find := strings.Contains(nickNames, user.NickName); !find {
              nickNames += "," + user.NickName
              //if len(strings.Split(nickNames, ","))  <= 3 {
                // 存储评论id_点赞类型 -> 点赞的用户昵称
                mp[fmt.Sprintf("%d_%d", video.VideoId, liked.ZanType)] = nickNames
              //}
            }

          } else {
            // 存储 视频id_点赞类型 -> 点赞的用户昵称
            mp[fmt.Sprintf("%d_%d", video.VideoId, liked.ZanType)] = user.NickName
          }
        }
      }
    // 被点赞的帖子
    case consts.TYPE_POSTS:
    // 被点赞的评论
    case consts.TYPE_COMMENT:
      // 获取评论信息
      comment := svc.comment.GetVideoCommentById(fmt.Sprint(liked.TypeId))
      if comment != nil {
        commentMp[comment.Id] = comment
        // 点赞用户
        if user := svc.user.FindUserByUserid(liked.UserId); user != nil {
          userMp[user.UserId] = user
          nickNames, ok := mp[fmt.Sprintf("%d_%d", comment.Id, liked.ZanType)]
          // 如果点赞的是同一评论  整合为一条数据
          if ok {
            if find := strings.Contains(nickNames, user.NickName); !find {
              nickNames += "," + user.NickName
              //if len(strings.Split(nickNames, ","))  <= 3 {
                // 存储评论id_点赞类型 -> 点赞的用户昵称
              mp[fmt.Sprintf("%d_%d", comment.Id, liked.ZanType)] = nickNames
              //}
            }

            continue
          } else {
            // 存储评论id_点赞类型 -> 点赞的用户昵称
            mp[fmt.Sprintf("%d_%d", comment.Id, liked.ZanType)] = user.NickName
          }
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

  // 是否已记录读取的位置
  //var b bool
  // 上次已读取的数据下标（默认-1未读取）
  //var readIndex int32 = -1
  res := make([]interface{}, 0)
  for _, liked := range list {
    switch liked.ZanType {
    // 被点赞的视频
    case consts.TYPE_VIDEOS:
      info := new(mlike.BeLikedInfo)
      nickNames, ok := mp[fmt.Sprintf("%d_%d", liked.TypeId, liked.ZanType)]
      if ok && nickNames != "" {
        log.Log.Errorf("ok:%b, nicknames:%s", ok, nickNames)
        info.OpTime = liked.CreateAt
        info.Type = consts.TYPE_VIDEOS
        // 视频作品
        video, ok := videoMp[liked.TypeId]
        if ok {
          mp[fmt.Sprintf("%d_%d", liked.TypeId, liked.ZanType)] = ""
          info.ComposeId = video.VideoId
          info.Title = video.Title
          info.Describe = video.Describe
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

          user, ok := userMp[liked.UserId]
          if ok {
            log.Log.Errorf("user.Id:%s", user.UserId)
            info.Avatar = user.Avatar
          }

          // 同一视频点赞的用户昵称（多个）
          info.Nicknames = strings.Split(nickNames, ",")
          if len(info.Nicknames) > 3 {
            info.Nicknames = info.Nicknames[0:3]
          }

          info.TotalLikeNum = len(info.Nicknames)
          res = append(res, info)
        }
      }
      // 被点赞的帖子
    case consts.TYPE_POSTS:
    // 被点赞的评论
    case consts.TYPE_COMMENT:
      info := new(mlike.BeLikedInfo)
      info.OpTime = liked.CreateAt
      info.Type = consts.TYPE_COMMENT

      nickNames, ok := mp[fmt.Sprintf("%d_%d", liked.TypeId, liked.ZanType)]
      if ok && nickNames != "" {

        // 获取评论信息
        comment, ok := commentMp[liked.TypeId]
        if ok {
          mp[fmt.Sprintf("%d_%d", liked.TypeId, liked.ZanType)] = ""
          // 被点赞的信息
          info.Content = comment.Content
          info.ComposeId = comment.Id

          // 获取评论对应的视频信息
          video, ok := videoMp[comment.VideoId]
          if ok {
            info.Title = video.Title
            info.Describe = video.Describe
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

          user, ok := userMp[liked.UserId]
          if ok {
            info.Avatar = user.Avatar
          }

          info.Nicknames = strings.Split(nickNames, ",")
          if len(info.Nicknames) > 3 {
            info.Nicknames = info.Nicknames[0:3]
          }
          info.TotalLikeNum = len(info.Nicknames)

          res = append(res, info)
        }

      }
    }

    // 未记录读取的下标
    //if !b {
    //  // 用户上次读取的数据下标
    //  if lastRead >= liked.CreateAt {
    //    readIndex = int32(len(res)-1)
    //    b = true
    //  }
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

  // 是否已记录读取的位置
  var b bool
  // 上次已读取的数据下标（默认-1未读取）
  var readIndex = -1
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
			  // 评论id
			  info.CommentId = receiveAt.CommentId
        // 被@的用户信息
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

				// 默认1级评论
				info.CommentType = 1
        info.Content = comment.Content
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
            }
					}

					// 获取最上级的评论内容
					parent := svc.comment.GetVideoCommentById(fmt.Sprint(comment.ParentCommentId))
					if parent != nil {
					  info.ParentComment = parent.Content
          }
				}

        if userId != "" {
          // 获取点赞的信息
          if likeInfo := svc.like.GetLikeInfo(userId, comment.Id, consts.TYPE_COMMENT); likeInfo != nil {
            info.IsLike = likeInfo.Status
          }
        }

				res[index] = info

			}
		}

    // 未记录读取的下标
    if !b {
      log.Log.Errorf("receiveId: %d", receiveAt.Id)
      if lastRead < receiveAt.CreateAt {
        // 用户上次读取的数据下标
        readIndex = index
        b = true
      }

      //if lastRead >= receiveAt.CreateAt && id != int(receiveAt.Id) {
      //  id = int(receiveAt.Id)
      //  log.Log.Errorf("lastRead:%d, createAt:%d, index:%d, len(res):%d", lastRead, receiveAt.CreateAt, index, len(res)-1)
      //  readIndex = index
      //  // 如果数据长度 - 1 == 已读取的下标 表示当前页数据读取完毕 返回-2
      //  if len(res) == readIndex  {
      //    readIndex = -2
      //  }
      //
      //  b = true
      //}
    }
	}

	// 记录读取@通知消息的时间
	if err := svc.notify.RecordReadAtTime(userId); err != nil {
		log.Log.Errorf("notify_trace: record read at notify time err:%s", err)
	}

	return res, readIndex
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

  ids := make([]string, len(list))
  var msgIds string
  for index, msg := range list {
    // 0 未读 1 已读
    if msg.Status == 0 {
      ids[index] = fmt.Sprint(msg.SystemId)
    }
  }


  msgIds = strings.Join(ids, ",")
  // 更新为已读
  if err := svc.notify.UpdateSystemNotifyStatus(msgIds); err != nil {
    log.Log.Errorf("notify_trace: update system notify status err:%s", err)
  }

  return list
}
