package job

import (
  "sports_service/server/global/app/log"
  "sports_service/server/models/mnotify"
  "sports_service/server/dao"
  "time"
)

// 检测定时推送 是否已发送（每2分钟）
func CheckTimedNotify() {
  ticker := time.NewTicker(time.Minute * 2)
  defer ticker.Stop()

  for {
    select {
    case <- ticker.C:
      log.Log.Debugf("开始检测定时广播推送 是否已发送")
      // 检测定时推送 是否已发送
      checkTimedNotify()
      log.Log.Debugf("检测完毕")
    }
  }

}

// 执行检测任务 检测定时广播推送 是否已发送 修改消息发送状态
func checkTimedNotify() {
  session := dao.AppEngine.NewSession()
  defer session.Close()
  nmodel := mnotify.NewNotifyModel(session)
  list := nmodel.GetAllSystemNotify()
  if list == nil {
    return
  }

  now := int(time.Now().Unix())
  for _, notify := range list {
    // 发送时间 > 当前时间 或 状态不是未发送
    if notify.SendTime > now || notify.SendStatus != 1 {
      continue
    }

    // 发送时间 <= 当前时间 表示已发送
    if notify.SendTime <= now {
      notify.SendStatus = 0
    }

    nmodel.UpdateSendStatus(0, notify.SystemId)
  }
}

