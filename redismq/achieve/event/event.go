package event

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/global/rdskey"
	"sports_service/server/models"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mlike"
	"sports_service/server/models/mnotify"
	"sports_service/server/models/muser"
	"sports_service/server/models/umeng"
	"sports_service/server/nsqlx/protocol"
	"sports_service/server/util"
	"time"
	producer "sports_service/server/redismq/event"
)

func LoopPopStatEvent() {
	for !closing {
		conn := dao.RedisPool().Get()
		values, err := redis.Values(conn.Do("BRPOP", rdskey.MSG_EVENT_KEY, 0))
		conn.Close()
		if err != nil {
			log.Log.Errorf("redisMq_trace: loopPop event fail, err:%s", err)
			// 防止出现错误时 频繁刷日志
			time.Sleep(time.Second)
			continue
		}

		if len(values) < 2 {
			log.Log.Errorf("redisMq_trace: invalid values, len:%d, values:%+v", len(values), values)
		}


		bts, ok := values[1].([]byte)
		if !ok {
			log.Log.Errorf("redisMq_trace: value[1] unSupport type")
			continue
		}

		if err := EventConsumer(bts); err != nil {
			log.Log.Errorf("redisMq_trace: event consumer fail, err:%s, msg:%s", err, string(bts))
			// 重新投递消息
			producer.PushEventMsg(bts)
		}

	}
}

func EventConsumer(bts []byte) error {
	event := protocol.Event{}
	if err := util.JsonFast.Unmarshal(bts, &event); err != nil {
		log.Log.Errorf("redisMq_trace: proto unmarshal event err:%s", err)
		return err
	}

	if err := handleEvent(event); err != nil {
		log.Log.Errorf("handleEvent err:%s", err)
		return err
	}

	return nil
}

func handleEvent(event protocol.Event) error {
	info := &protocol.Data{}
	if err := util.JsonFast.Unmarshal(event.Data, info); err != nil {
		log.Log.Errorf("redisMq_trace: proto unmarshal data err:%s", err)
		return nil
	}

	session := dao.Engine.NewSession()
	defer session.Close()
	umodel := muser.NewUserModel(session)
	user := umodel.FindUserByUserid(event.UserId)
	if user == nil {
		log.Log.Errorf("redisMq_trace: user not found, userId:%s", event.UserId)
		return nil
	}

	nmodel := mnotify.NewNotifyModel(session)
	// 系统消息未读数
	sysNum := nmodel.GetUnreadSystemMsgNum(event.UserId)

	// likedNum 未读的被点赞的数量 atNum 未读的被@的数量
	var (
		likedNum, atNum int64
	)
	// 获取用户上次读取被点赞列表的时间
	readTm, err := nmodel.GetReadBeLikedTime(event.UserId)
	if err == nil || err == redis.ErrNil {
		if readTm == "" {
			readTm = "0"
		}

		lmodel := mlike.NewLikeModel(session)
		// 获取未读的被点赞的数量
		likedNum = lmodel.GetUnreadBeLikedCount(event.UserId, readTm)
	}

	// 获取用户上次读取被@列表数据的时间
	readAt, err := nmodel.GetReadAtTime(event.UserId)
	if err == nil || err == redis.ErrNil {
		if readAt == "" {
			readAt = "0"
		}

		cmodel := mcomment.NewCommentModel(session)
		// 获取未读的被@的数量
		atNum = cmodel.GetUnreadAtCount(event.UserId, readAt)
	}

	// 总未读数
	unReadNum := sysNum + likedNum + atNum

	var (
		content string
		msgType int32
	)

	log.Log.Infof("redisMq_trace: event:%+v", event)
	setting := nmodel.GetUserNotifySetting(event.UserId)

	var pushSet int
	switch event.EventType {
	// 系统类
	case consts.SYSTEM_MSG:

	// 活动类
	case consts.ACTIVITY_MSG:

	// 视频点赞
	case consts.VIDEO_LIKE_MSG:
		content = fmt.Sprintf("%s 赞了你的作品", info.NickName)
		msgType = int32(consts.MSG_TYPE_VIDEO_LIKE_NOTIFY)
		pushSet = setting.ThumbUpPushSet
	// 评论/回复 点赞
	case consts.COMMENT_LIKE_MSG:
		content = fmt.Sprintf("%s 赞了你的评论 @%s", info.NickName, info.Content)
		msgType = int32(consts.MSG_TYPE_COMMENT_LIKE_NOTIFY)
		pushSet = setting.ThumbUpPushSet
	// 收藏视频
	case consts.COLLECT_VIDEO_MSG:
		content = fmt.Sprintf("%s 收藏了你的作品", info.NickName)
		msgType = int32(consts.MSG_TYPE_VIDEO_COLLECT_NOTIFY)
	// 关注用户
	case consts.FOCUS_USER_MSG:
		content = fmt.Sprintf("%s 关注了你", info.NickName)
		msgType = int32(consts.MSG_TYPE_FOCUS_NOTIFY)
		pushSet = setting.AttentionPushSet
	// 关注的用户发布视频
	case consts.FOCUS_USER_PUBLISH_MSG:
		content = fmt.Sprintf("你关注的 %s 发布了新视频", info.NickName)
		msgType = int32(consts.MSG_TYPE_FOCUS_USER_PUBLISH_NOTIFY)
		pushSet = setting.AttentionPushSet
	// 视频评论
	case consts.VIDEO_COMMENT_MSG:
		content = fmt.Sprintf("%s 评论了你的作品", info.NickName)
		msgType = int32(consts.MSG_TYPE_VIDEO_COMMENT_NOTIFY)
		pushSet = setting.CommentPushSet
	// 视频回复
	case consts.VIDEO_REPLY_MSG:
		// @ 用户发布的评论
		content = fmt.Sprintf("%s 回复了你的评论 @%s", info.NickName, info.Content)
		msgType = int32(consts.MSG_TYPE_VIDEO_REPLY_NOTIFY)
		pushSet = setting.CommentPushSet
	default:
		log.Log.Errorf("redisMq_trace: unsupported eventType, eventType:%d", event.EventType)
		return nil

	}

	// 0为接收推送 1为拒绝接收
	if pushSet == 0 {
		// 推送通知
		PushNotify(user, "", content, info.Cover, msgType, unReadNum)
	}

	return nil
}

// 推送通知
func PushNotify(user *models.User, title, content, cover string, msgType int32, unreadNum int64) {
	extra := make(map[string]interface{}, 0)
	extra["unread_num"] = unreadNum
	title = "X-FLY官方"
	// android推送
	if user.DeviceType == int(consts.ANDROID_PLATFORM) && user.DeviceToken != "" {
		client := umeng.New()
		if err := client.PushUnicastNotify(msgType, umeng.FPV_ANDROID, user.DeviceToken, title, content, cover, extra, nil); err != nil {
			log.Log.Errorf("redisMq_trace: push notify by user err:%s", err)
		}
	}

	// iOS推送
	if user.DeviceType == int(consts.IOS_PLATFORM) && user.DeviceToken != "" {
		client := umeng.New()
		if err := client.PushUnicastNotify(msgType, umeng.FPV_IOS, user.DeviceToken, title, content, cover, extra, nil); err != nil {
			log.Log.Errorf("redisMq_trace: push notify by user err:%s", err)
		}
	}
}


