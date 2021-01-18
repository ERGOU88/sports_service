package cnotify

import (
  "github.com/gin-gonic/gin"
  "github.com/go-xorm/xorm"
  "sports_service/server/global/backend/errdef"
  "sports_service/server/global/consts"
  "sports_service/server/models"
  "sports_service/server/models/umeng"
  "sports_service/server/tools/tencentCloud"
  "strings"
  "time"
  "fmt"
  "sports_service/server/models/muser"
  "sports_service/server/models/mnotify"
  "sports_service/server/dao"
  "sports_service/server/global/backend/log"
  umengClient "sports_service/server/tools/umeng"
)

type NotifyModule struct {
  context     *gin.Context
  engine      *xorm.Session
  user        *muser.UserModel
  notify      *mnotify.NotifyModel
}

func New(c *gin.Context) NotifyModule {
  socket := dao.Engine.NewSession()
  defer socket.Close()
  return NotifyModule{
    context: c,
    user: muser.NewUserModel(socket),
    notify: mnotify.NewNotifyModel(socket),
    engine: socket,
  }
}

// 获取系统推送列表
func (svc *NotifyModule) GetSystemNotifyList(page, size int) []*models.SystemMessage {
  offset := (page - 1) * size
  list := svc.notify.GetSystemNotifyList(offset, size)
  if list == nil {
    return []*models.SystemMessage{}
  }

  return list
}

// 后台推送通知
func (svc *NotifyModule) PushSystemNotify(param *umeng.SystemNotifyParams) int {
  var (
    duration, sendStatus int
    policy   interface{}
  )

  now := int(time.Now().Unix())
  if param.SendTm <= 0 {
    // 0 已发送
    sendStatus = 0
    param.SendTm = now
  } else {
    if param.SendTm < now {
      log.Log.Errorf("notify_trace: invalid send time, sendTm:%d, now:%d", param.SendTm, now)
      return errdef.NOTIFY_INVALID_SEND_TM
    }

    duration = param.SendTm - now
    if duration >= 7 * 3600 * 24 {
      log.Log.Error("notify_trace: invalid send time, max: 7 days")
      return errdef.NOTIFY_INVALID_START_TM
    }

    // 1 未发送 （定时发送）
    sendStatus = 1
    policy = umengClient.Policy{StartTime: time.Now().Add(time.Second * time.Duration(duration)).Format(consts.FORMAT_TM)}
  }

  client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
  // 检测推送内容
  isPass, err := client.TextModeration(param.Content)
  if !isPass || err != nil {
    log.Log.Errorf("notify_trace: invalid send content, content:%s, err:%s", param.Content, err)
    return errdef.NOTIFY_INVALID_CONTENT
  }

  umodel := umeng.New()
  switch param.SendType {
  // -1 全部用户推送
  case -1:
    // android端广播推送
    taskId, err := umodel.PushBroadcastNotifyByAndroid(param.Topic, param.Content, nil, policy)
    if err != nil {
      log.Log.Errorf("notify_trace: push android broadcast notify err:%s", err)
      return errdef.NOTIFY_PUSH_FAIL
    }

    notify := new(models.SystemMessage)
    notify.SendType = 0
    notify.Status = 0
    notify.SendTime = param.SendTm
    notify.ReceiveId = ""
    notify.SystemContent = param.Content
    notify.SystemTopic = param.Topic
    notify.SendDefault = 1
    notify.SendId = "admin"
    notify.CreateAt = now
    notify.AndroidTaskId = taskId
    notify.SendStatus = sendStatus

    // ios端广播推送
    taskId, err = umodel.PushBroadcastNotifyByIos(param.Topic, param.Content, nil, policy)
    if err != nil {
      log.Log.Errorf("notify_trace: push ios broadcast notify err:%s", err)
    }
    notify.IosTaskId = taskId

    affected, err := svc.notify.AddSystemNotify(notify)
    if affected != 1 || err != nil {
      log.Log.Errorf("notify_trace: add system notify err:%s", err)
      return errdef.NOTIFY_PUSH_FAIL
    }

  // 1 指定用户推送
  case 1:
    userIds := strings.Split(param.UserIds, ",")
    if len(userIds) == 0 {
      log.Log.Errorf("notify_trace: invalid userId:%s", param.UserIds)
      return errdef.NOTIFY_INVALID_USER_IDS
    }

    // 查询用户列表
    list := svc.user.FindUserByUserids(param.UserIds, 0, len(userIds))
    if list == nil {
      log.Log.Errorf("notify_trace: user not found, userIds:%s", param.UserIds)
      return errdef.NOTIFY_USER_NOT_FOUND
    }

    notifyList := make([]*models.SystemMessage, 0)
    for _, user := range list {
      if user.DeviceToken == "" {
        log.Log.Errorf("notify_trace: user device token empty, userId:%s, token:%s", user.UserId, user.DeviceToken)
        continue
      }

      notify := new(models.SystemMessage)
      notify.SendType = 0
      notify.Status = 0
      notify.SendTime = param.SendTm
      notify.ReceiveId = user.UserId
      notify.SystemContent = param.Content
      notify.SystemTopic = param.Topic
      notify.SendDefault = 0
      notify.SendId = "admin"
      notify.CreateAt = now
      notify.SendStatus = sendStatus
      // android端
      if user.DeviceType == int(consts.ANDROID_PLATFORM) {
        if err := umodel.PushUnicastNotify(int32(consts.MSG_TYPE_SYSTEM_NOTIFY), umeng.FPV_ANDROID, user.DeviceToken,
          param.Topic, "您有一条新消息，点我查看～", "", nil, policy); err != nil {
          log.Log.Errorf("notify_trace: push unicast notify err:%s, userId:%s", err, user.UserId)
        }

        notify.UmengPlatform = umeng.FPV_ANDROID
      }

      // ios端
      if user.DeviceType == int(consts.IOS_PLATFORM) {
        if err := umodel.PushUnicastNotify(int32(consts.MSG_TYPE_SYSTEM_NOTIFY), umeng.FPV_IOS, user.DeviceToken,
          param.Topic, "您有一条新消息，点我查看～","", nil, policy); err != nil {
          log.Log.Errorf("notify_trace: push unicast notify err:%s, userId:%s", err, user.UserId)
        }

        notify.UmengPlatform = umeng.FPV_IOS
      }

      notifyList = append(notifyList, notify)
    }

    if len(notifyList) == 0 {
      return errdef.NOTIFY_PUSH_FAIL
    }

    affected, err := svc.notify.AddMultiSystemNotify(notifyList)
    if int(affected) != len(userIds) || err != nil {
      log.Log.Errorf("notify_trace: add multi system notify err:%s, affected:%d, len:%d", err, affected, len(userIds))
      return errdef.NOTIFY_PUSH_FAIL
    }
  }

  return errdef.SUCCESS
}

// 撤回系统推送
func (svc *NotifyModule) CancelSystemNotify(systemId int64) int {
  msg := svc.notify.GetSystemNotifyById(fmt.Sprint(systemId))
  if msg == nil {
    log.Log.Errorf("notify_trace: system notify not found, id:%d", systemId)
    return errdef.NOTIFY_MSG_NOT_EXISTS
  }

  if msg.AndroidTaskId == "" && msg.IosTaskId == "" {
    log.Log.Error("notify_trace: can not cancel, taskId is empty")
    return errdef.NOTIFY_CAN_NOT_CANCEL
  }

  // 只有未发送状态的消息才可撤回
  if msg.SendStatus != 1 {
    log.Log.Error("notify_trace: can not cancel, send status not 0")
    return errdef.NOTIFY_CAN_NOT_CANCEL
  }

  umodel := umeng.New()
  if err := umodel.CancelNotify(msg.AndroidTaskId, umeng.FPV_ANDROID); err != nil {
    log.Log.Errorf("notify_trace: cancel android notify err:%s", err)
    return errdef.NOTIFY_CANCEL_FAIL
  }

  if err := umodel.CancelNotify(msg.IosTaskId, umeng.FPV_IOS); err != nil {
    log.Log.Errorf("notify_trace: cancel ios notify err:%s", err)
    return errdef.NOTIFY_CANCEL_FAIL
  }

  // 将消息状态设置为已撤回(send_status 2)
  if err := svc.notify.UpdateSendStatus(2, systemId); err != nil {
    log.Log.Errorf("notify_trace: update send status err:%s", err)
    return errdef.NOTIFY_CANCEL_FAIL
  }

  return errdef.SUCCESS
}

// 撤回系统推送
func (svc *NotifyModule) DelSystemNotify(systemId int64) int {
  msg := svc.notify.GetSystemNotifyById(fmt.Sprint(systemId))
  if msg == nil {
    log.Log.Errorf("notify_trace: system notify not found, id:%d", systemId)
    return errdef.NOTIFY_MSG_NOT_EXISTS
  }

  // 未发送状态的消息 需要先撤回 方可删除
  if msg.SendStatus == 1 {
    log.Log.Error("notify_trace: can not del, send status is 1")
    return errdef.NOTIFY_CAN_NOT_DEL
  }

  affected, err := svc.notify.DelSystemNotify(systemId)
  if affected != 1 || err != nil {
    log.Log.Errorf("notify_trace: del system notify err:%s", err)
    return errdef.NOTIFY_DEL_FAIL
  }

  return errdef.SUCCESS
}
