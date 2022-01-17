package umeng_test

import (
  "sports_service/server/tools/umeng"
  "testing"
  "time"
  "encoding/json"
  "fmt"
)

type PushMessage struct {
  MsgId       string                     `json:"msg_id"`   // 消息id
  Data        map[string]interface{}     `json:"data"`     // 具体数据
  MsgType     int32                      `json:"msg_type"` // 消息类型
  SendTime    int64                      `json:"send_time"`// 发送时间
  Display     bool                       `json:"display"`  // 是否展示 false不展示
}

//var data *umeng.Data
//
//func init() {
// data = umeng.NewData(umeng.APP_IOS)
// //data := umeng.NewData(umeng.APP_ANDROID)
// data.Type = "unicast"
// data.TimeStamp = time.Now().Unix()
// data.DeviceTokens = "3fa23e7d5afa9705f7a3d7161a64bb99531100bf818b2ea6d35a14d8dec6b6ce"
// data.Description = "风暴英雄"
// data.ProductionMode = false
//
//}
//
//func TestStatus(t *testing.T) {
// resp := data.Status()
// if len(resp.Code) > 0 {
//   t.Log("data.Status 测试通过")
// } else {
//   t.Error("data.Status 测试不通过")
// }
//
//}
//
//func TestCancel(t *testing.T) {
// resp := data.Cancel()
// if len(resp.Code) > 0 {
//   t.Log("data.Cancel 测试通过")
// } else {
//   t.Error("data.Cancel 测试不通过")
// }
//}
//
//func TestUpload(t *testing.T) {
// resp := data.Upload()
// if len(resp.Code) > 0 {
//   t.Log("data.Upload 测试通过")
// } else {
//   t.Error("data.Upload 测试不通过")
// }
//}
//

//
//func TestPush(t *testing.T) {
//  body := PushMessage{
//    MsgId: "123456",
//    SendTime: time.Now().Unix(),
//    MsgType: 100,
//    Display: false,
//  }
//
//  body.Data = make(map[string]interface{}, 0)
//  body.Data["content"] = "测试测试测试～～～～～"
//  bts, _ := util.JsonFast.Marshal(body)
//  info := umeng.Alert{
//   SkuName: "123456",
//   SubTitle: "23232",
//   Body: "测试测试",
//  }
//
//  extras := make(map[string]string, 0)
//  extras["extra"] = string(bts)
//
//  aps := &umeng.IOSAps{
//    Alert: info,
//    Sound: "default",
//  }
//
//  policy := umeng.Policy{ExpireTime: time.Now().Add(time.Hour * 100).Format("2006-01-02 15:04:05")}
//  resp := data.Push(nil, aps, policy, extras)
//  if len(resp.Code) > 0 {
//    t.Logf("data.Push 测试通过:%+v", resp)
//  } else {
//    t.Error("data.Push 测试不通过")
//  }
//}

var data *umeng.Data

func init() {
  data = umeng.NewData(umeng.APP_ANDROID)
  data.Type = "unicast"
  data.TimeStamp = time.Now().Unix()
  data.DeviceTokens = "AqWqyJjw6UvTTE9eRXYuYjHu20A8oX6iW0J-_zE2YjY0"
  //data.Description = "威猛先生"
  data.ProductionMode = false
}

func TestStatus(t *testing.T) {
 resp := data.Status()
 if len(resp.Code) > 0 {
   t.Log("data.Status 测试通过")
 } else {
   t.Error("data.Status 测试不通过")
 }

}

func TestCancel(t *testing.T) {
  resp := data.Cancel()
  if len(resp.Code) > 0 {
    t.Log("data.Cancel 测试通过")
  } else {
    t.Error("data.Cancel 测试不通过")
  }
}

func TestUpload(t *testing.T) {
  resp := data.Upload()
  if len(resp.Code) > 0 {
    t.Log("data.Upload 测试通过")
  } else {
    t.Error("data.Upload 测试不通过")
  }
}

func TestPush(t *testing.T) {
  body := umeng.AndroidBody{}
  // android: notification 通知栏推送  message 自定义推送
  body.DisplayType = "message"
  body.Title = "订单通知"
  body.Text = "计时10分钟付款倒计时10分钟付款倒计时10分钟付款倒计时10分钟"

  msg := PushMessage{
    MsgId: "123456",
    SendTime: time.Now().Unix(),
    MsgType: 100,
    Display: true,
  }

  msg.Data = make(map[string]interface{}, 0)
  msg.Data["unread_num"] = 100

  bts, _ := json.Marshal(msg)
  fmt.Println("\nbts", string(bts))
  body.Custom = string(bts)
  //policy := umeng.Policy{ExpireTime: time.Now().Add(time.Hour * 100).Format("2006-01-02 15:04:05")}

  //extra := make(map[string]string, 0)
  //extra["extra"] = string(bts)
  resp := data.Push(body, nil, nil, nil)
  if len(resp.Code) > 0 {
    t.Logf("data.Push 测试通过, resp:%+v", resp)
  } else {
    t.Error("data.Push 测试不通过")
  }
}

