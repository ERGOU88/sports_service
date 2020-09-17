package cnotify

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/models/mnotify"
	"time"
)

type NotifyModule struct {
	context    *gin.Context
	engine     *xorm.Session
	notify     *mnotify.NotifyModel
}

func New(c *gin.Context) NotifyModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return NotifyModule{
		context: c,
		notify: mnotify.NewNotifyModel(socket),
		engine: socket,
	}
}

// 保存用户通知设置
func (svc *NotifyModule) SaveUserNotifySetting(userId string, params *mnotify.NotifySettingParams) int {
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
