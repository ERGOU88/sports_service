package umeng

import (
  "sports_service/server/global/consts"
  "sports_service/server/util"
  "time"
  "errors"
  "fmt"
  "sports_service/server/tools/umeng"
)

type UmengModel struct {
  Data    *umeng.Data
}

type PushMessage struct {
  MsgId       string                     `json:"msg_id"`   // 消息id
  Data        map[string]interface{}     `json:"data"`     // 具体数据
  MsgType     int32                      `json:"msg_type"` // 消息类型
  SendTime    int64                      `json:"send_time"`// 发送时间
  Display     bool                       `json:"display"`  // 是否展示 false不展示
}

// 后台系统推送通知请求参数
type SystemNotifyParams struct {
  SendType        int32       `json:"send_type" binding:"required"`   // -1 全部 1 指定用户发送
  UserIds         string      `json:"user_ids"`                       // 如果指定用户id 多个用逗号分隔
  Content         string      `json:"content" binding:"required"`     // 推送内容
  Topic           string      `json:"topic"`                          // 推送标题
  SendTm          int         `json:"send_tm"`                        // 指定发送时间
}

// 撤回系统推送请求参数
type CancelSystemNotifyParam struct {
  SystemId        int64        `json:"system_id" binding:"required"`   // 系统消息id
}

// 删除系统推送请求参数
type DelSystemNotifyParam struct {
  SystemId        int64        `json:"system_id" binding:"required"`   // 系统消息id
}

const (
  FPV_ANDROID = 1
  FPV_IOS     = 2
)

func New() (umodel *UmengModel) {
  umodel = new(UmengModel)
  //if pf == FPV_ANDROID {
  //  umodel.Data = umeng.NewData(umeng.APP_ANDROID)
  //  return
  //}
  //
  //if pf == FPV_IOS {
  //  umodel.Data = umeng.NewData(umeng.APP_IOS)
  //  return
  //}

  return umodel
}

// 推送通知(单播)
func (m *UmengModel) PushUnicastNotify(msgType, pf int32, deviceToken, title, content, cover string, extra map[string]interface{}, policy interface{}) error {
  data := umeng.NewData(umeng.Platform(pf))
  data.Type = "unicast"
  data.TimeStamp = time.Now().Unix()
  data.DeviceTokens = deviceToken
  //m.Data.Description = ""
  data.ProductionMode = false
  body := PushMessage{
    MsgId: fmt.Sprint(util.GetSnowId()),
    SendTime: time.Now().Unix(),
    MsgType: msgType,
    Display: true,
  }

  body.Data = make(map[string]interface{}, 0)
  if len(extra) > 0 {
    for key, val := range extra {
      body.Data[key] = val
    }
  }

  bts, _ := util.JsonFast.Marshal(body)

  if pf == FPV_ANDROID {
    body := umeng.AndroidBody{}
    // android: notification 通知栏推送  message 自定义推送
    body.DisplayType = consts.ANDROID_PUSH_TYPE_NOTIFICATION

    body.Custom = string(bts)
    body.Text = content
    body.Title = title
    if cover != "" {
      body.Img = cover
    }

    resp := data.Push(body, nil, policy, nil)
    if resp.Code != "SUCCESS" {
      return errors.New("push notify fail, error_msg" + resp.Data["error_msg"])
    }
  }

  if pf == FPV_IOS {
    extras := make(map[string]string, 0)
    extras["extra"] = string(bts)

    info := umeng.Alert{
      Title: title,
      SubTitle: title,
      Body: content,
    }

    aps := &umeng.IOSAps{
      Alert: info,
      Sound: "default",
    }

    resp := data.Push(nil, aps, policy, extras)
    if resp.Code != "SUCCESS" {
      return errors.New("push unicast notify fail, error_msg" + resp.Data["error_msg"])
    }

  }

  return nil
}

// 推送广播通知(android端)
func (m *UmengModel) PushBroadcastNotifyByAndroid(title, content string, extra map[string]interface{}, policy interface{}) (string, error) {
  // android
  data := umeng.NewData(umeng.APP_ANDROID)
  data.Type = "broadcast"
  data.TimeStamp = time.Now().Unix()
  data.ProductionMode = false

  body := PushMessage{
    MsgId:    fmt.Sprint(util.GetSnowId()),
    SendTime: time.Now().Unix(),
    MsgType:  int32(consts.MSG_TYPE_SYSTEM_NOTIFY),
    Display:  true,
  }

  body.Data = make(map[string]interface{}, 0)
  if len(extra) > 0 {
    for key, val := range extra {
      body.Data[key] = val
    }
  }

  bts, _ := util.JsonFast.Marshal(body)
  msgInfo := umeng.AndroidBody{}
  // android: notification 通知栏推送  message 自定义推送
  msgInfo.DisplayType = consts.ANDROID_PUSH_TYPE_NOTIFICATION

  msgInfo.Custom = string(bts)
  msgInfo.Text = content
  msgInfo.Title = title

  resp := data.Push(msgInfo, nil, policy, nil)
  if resp.Code != "SUCCESS" {
    return "", errors.New("push broadcast notify, error_msg" + resp.Data["error_msg"])
  }

  return resp.Data["task_id"], nil
}

// 推送广播通知(ios端)
func (m *UmengModel) PushBroadcastNotifyByIos(title, content string, extra map[string]interface{}, policy interface{}) (string, error) {
  // ios
  data := umeng.NewData(umeng.APP_IOS)
  data.Type = "broadcast"
  data.TimeStamp = time.Now().Unix()
  data.ProductionMode = false
  extras := make(map[string]string, 0)
  body := PushMessage{
    MsgId:    fmt.Sprint(util.GetSnowId()),
    SendTime: time.Now().Unix(),
    MsgType:  int32(consts.MSG_TYPE_SYSTEM_NOTIFY),
    Display:  true,
  }

  body.Data = make(map[string]interface{}, 0)
  if len(extra) > 0 {
    for key, val := range extra {
      body.Data[key] = val
    }
  }

  bts, _ := util.JsonFast.Marshal(body)
  extras["extra"] = string(bts)

  info := umeng.Alert{
    Title: title,
    SubTitle: title,
    Body: content,
  }

  aps := &umeng.IOSAps{
    Alert: info,
    Sound: "default",
  }

  resp := data.Push(nil, aps, policy, extras)
  if resp.Code != "SUCCESS" {
    return "", errors.New("push broadcast notify, error_msg" + resp.Data["error_msg"])
  }

  return resp.Data["task_id"], nil
}

// 撤回通知
func (m *UmengModel) CancelNotify(taskId string, pf int) error {
  data := umeng.NewData(umeng.Platform(pf))
  data.TimeStamp = time.Now().Unix()
  data.TaskId = taskId
  resp := data.Cancel()
  if resp.Code != "SUCCESS" {
    return errors.New("cancel notify fail, error_msg" + resp.Data["error_msg"])
  }

  return nil
}
