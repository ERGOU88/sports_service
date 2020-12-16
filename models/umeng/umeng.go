package umeng

import (
  "errors"
  "fmt"
  "sports_service/server/global/app/log"
  "sports_service/server/global/consts"
  "sports_service/server/tools/umeng"
  "sports_service/server/util"
  "time"
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

const (
  FPV_ANDROID = 1
  FPV_IOS     = 2
)

func New(pf int32) (umodel *UmengModel) {
  umodel = new(UmengModel)
  if pf == FPV_ANDROID {
    umodel.Data = umeng.NewData(umeng.APP_ANDROID)
    return
  }

  if pf == FPV_IOS {
    umodel.Data = umeng.NewData(umeng.APP_IOS)
    return
  }

  return
}

// 推送消息(单播)
func (m *UmengModel) PushUnicastNotify(msgType, pf int32, deviceToken, title, content, cover string, extra map[string]interface{}) error {
  m.Data.Type = "unicast"
  m.Data.TimeStamp = time.Now().Unix()
  m.Data.DeviceTokens = deviceToken
  //m.Data.Description = ""
  m.Data.ProductionMode = false

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
    body.DisplayType = consts.ANDROID_PUSH_TYPE_CUSTOM

    body.Custom = string(bts)
    body.Text = content
    body.Title = title
    body.Img = cover

    log.Log.Errorf("event_trace: msg:%+v", m.Data)
    resp := m.Data.Push(body, nil, nil, extras)
    if resp.Code != "SUCCESS" {
      log.Log.Errorf("event_trace: umeng push errCode:%s", resp.Code)
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

    resp := m.Data.Push(nil, aps, nil, extras)
    if resp.Code != "SUCCESS" {
      return errors.New("push order timeOut fail, error_msg" + resp.Data["error_msg"])
    }
  }

  log.Log.Error("event_trace: push notify success")
  return nil

}
