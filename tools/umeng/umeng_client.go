package umeng

import (
  "bytes"
  "crypto/md5"
  "encoding/hex"
  "fmt"
  "io/ioutil"
  "net/http"
)

var (
  Host                          = "http://msg.umeng.com"
  UploadPath                    = "/upload"
  StatusPath                    = "/api/status"
  CancelPath                    = "/api/cancel"
  PostPath                      = "/api/send"
  AndroidAppKey                 = "5fab8a2545b2b751a929328c"
  IOSAppKey                     = "5fc60a125a31bc1b28ee9cb7"
  AndroidAppMasterSecret        = "ultl2ip6u5x75vljp9zcwht9qbgdsinv"
  IOSAppMasterSecret            = "zynivoaugy5nja5tsui2dhrsuerndzmp"
  DataByte               []byte = nil
)

type Platform int32

const (
  APP_ANDROID Platform = 1
  APP_IOS     Platform = 2
)

/*
可选 发送策略
"policy":{
	"start_time":"xx",   // 可选 定时发送时间，若不填写表示立即发送。定时发送时间不能小于当前时间 格式: "YYYY-MM-DD HH:mm:ss"。注意, start_time只对任务生效。
	"expire_time":"xx",  // 可选 消息过期时间,其值不可小于发送时间或者 start_time(如果填写了的话),如果不填写此参数，默认为3天后过期。格式同start_time
	"max_send_num": xx   // 可选 发送限速，每秒发送的最大条数。开发者发送的消息如果有请求自己服务器的资源，可以考虑此参数。
	"out_biz_no": "xx"   // 可选 开发者对消息的唯一标识，服务器会根据这个标识避免重复发送。有些情况下（例如网络异常）开发者可能会重复调用API导致消息多次下发到客户端。如果需要处理这种情况，可以考虑此参数。注意, out_biz_no只对任务生效。
}
*/
type Policy struct {
  StartTime  string `json:"start_time,omitempty"`
  ExpireTime string `json:"expire_time,omitempty"`
  MaxSendNum int    `json:"max_send_num,omitempty"`
  OutBizNo   string `json:"out_biz_no,omitempty"`
}

/*
必填 消息内容(Android最大为1840B), 包含参数说明如下(JSON格式)
"payload":{
    "display_type":"xx",  // 必填 消息类型，值可以为:
    "body":               // 必填 消息体。display_type=message时,body的内容只需填写custom字段。display_type=notification时, body包含如下参数:
    {
        // 通知展现内容:
        "ticker":"xx",     // 必填 通知栏提示文字
        "title":"xx",      // 必填 通知标题
        "text":"xx",       // 必填 通知文字描述
        // 自定义通知图标:
        "icon":"xx",       // 可选 状态栏图标ID, R.drawable.[smallIcon],
                                   如果没有, 默认使用应用图标。
                                   图片要求为24
*24dp的图标,或24*24px放在drawable-mdpi下。
                                   注意四周各留1个dp的空白像素
        "largeIcon":"xx",  // 可选 通知栏拉开后左侧图标ID, R.drawable.[largeIcon].
                                   图片要求为64*64dp的图标,
                                   可设计一张64*64px放在drawable-mdpi下,
                                   注意图片四周留空，不至于显示太拥挤
        "img": "xx",       // 可选 通知栏大图标的URL链接。该字段的优先级大于largeIcon。
                                   该字段要求以http或者https开头。
        // 自定义通知声音:
        "sound": "xx",     // 可选 通知声音，R.raw.[sound].
                                   如果该字段为空，采用SDK默认的声音, 即res/raw/下的
                                       umeng_push_notification_default_sound声音文件
                                   如果SDK默认声音文件不存在，
                                       则使用系统默认的Notification提示音。
        // 自定义通知样式:
        "builder_id": xx   // 可选 默认为0，用于标识该通知采用的样式。使用该参数时,
                       开发者必须在SDK里面实现自定义通知栏样式。
        // 通知到达设备后的提醒方式
        "play_vibrate":"true/false", // 可选 收到通知是否震动,默认为"true".
                       注意，"true/false"为字符串
        "play_lights":"true/false",  // 可选 收到通知是否闪灯,默认为"true"
        "play_sound":"true/false",   // 可选 收到通知是否发出声音,默认为"true"
        // 点击"通知"的后续行为，默认为打开app。
        "after_open": "xx" // 必填 值可以为:
                                   "go_app": 打开应用
                                   "go_url": 跳转到URL
                                   "go_activity": 打开特定的activity
                                   "go_custom": 用户自定义内容。
        "url": "xx",       // 可选 当"after_open"为"go_url"时，必填。
                                   通知栏点击后跳转的URL，要求以http或者https开头
        "activity":"xx",   // 可选 当"after_open"为"go_activity"时，必填。
                                   通知栏点击后打开的Activity
        "custom":"xx"/{}   // 可选 display_type=message, 或者
                                   display_type=notification且
                                   "after_open"为"go_custom"时，
                                   该字段必填。用户自定义内容, 可以为字符串或者JSON格式。
    },
    extra:                 // 可选 用户自定义key-value。只对"通知"(display_type=notification)生效。可以配合通知到达后,打开App,打开URL,打开Activity使用。
    {
        "key1": "value1",
        "key2": "value2",
        ...
    }
}
*/
type AndroidBody struct {
  DisplayType string `json:"-"`                      // notification 通知栏推送  message 自定义推送
  Ticker      string `json:"ticker,omitempty"`
  Title       string `json:"title,omitempty"`
  Text        string `json:"text,omitempty"`
  Icon        string `json:"icon,omitempty"`
  LargeIcon   string `json:"largeIcon,omitempty"`
  Img         string `json:"img,omitempty"`
  Sound       string `json:"sound,omitempty"`
  BuilderId   string `json:"builder_id,omitempty"`
  PlayVibrate bool   `json:"play_vibrate,omitempty"`
  PlayLights  bool   `json:"play_lights,omitempty"`
  PlaySound   bool   `json:"play_sound,omitempty"`
  AfterOpen   string `json:"after_open,omitempty"`
  Url         string `json:"url,omitempty"`
  Activity    string `json:"activity,omitempty"`
  Custom      string `json:"custom,omitempty"`
}

type AndroidPayload struct {
  DisplayType string            `json:"display_type,omitempty"`
  Body        AndroidBody       `json:"body,omitempty"`
  Extra       map[string]string `json:"extra,omitempty"`
}

/*
 必填 消息内容(iOS最大为2012B), 包含参数说明如下(JSON格式):
 "payload":{
    "aps":                 // 必填 严格按照APNs定义来填写
    {
        "alert": "xx"          // 必填
        "badge": xx,           // 可选
        "sound": "xx",         // 可选
        "content-available":xx // 可选
        "category": "xx",      // 可选, 注意: ios8才支持该字段。
    },
    "key1":"value1",       // 可选 用户自定义内容, "d","p"为友盟保留字段，key不可以是"d","p"
    "key2":"value2",
    ...
   }
*/
type IOSPayload map[string]interface{}

type Alert struct {
  Title    string  `json:"title"`
  SubTitle string `json:"subtitle"`
  Body     string  `json:"body"`
}

type IOSAps struct {
  Alert            Alert  `json:"alert,omitempty"`
  Badge            string `json:"badge,omitempty"`
  Sound            string `json:"sound,omitempty"`
  ContentAvailable string `json:"content-available,omitempty"`
  Category         string `json:"category,omitempty"`
}

/*
{
	  "appkey":"xx",          // 必填 应用唯一标识
	  "timestamp":"xx",       // 必填 时间戳，10位或者13位均可，时间戳有效期为10分钟
	  "type":"xx",            // 必填 消息发送类型,其值可以为:
					  unicast-单播
					  listcast-列播(要求不超过500个device_token)
					  filecast-文件播
					    (多个device_token可通过文件形式批量发送）
					  broadcast-广播
					  groupcast-组播
					    (按照filter条件筛选特定用户群, 具体请参照filter参数)
					  customizedcast(通过开发者自有的alias进行推送),
					  包括以下两种case:
					     - alias: 对单个或者多个alias进行推送
					     - file_id: 将alias存放到文件后，根据file_id来推送
	  "device_tokens":"xx",   // 可选 设备唯一表示
					  当type=unicast时,必填, 表示指定的单个设备
					  当type=listcast时,必填,要求不超过500个,
					  多个device_token以英文逗号间隔
	  "alias_type": "xx"      // 可选 当type=customizedcast时，必填，alias的类型,
					  alias_type可由开发者自定义,开发者在SDK中
					  调用setAlias(alias, alias_type)时所设置的alias_type
	  "alias":"xx",           // 可选 当type=customizedcast时, 开发者填写自己的alias。
					  要求不超过50个alias,多个alias以英文逗号间隔。
					  在SDK中调用setAlias(alias, alias_type)时所设置的alias
	  "file_id":"xx",         // 可选 当type=filecast时，file内容为多条device_token,
					    device_token以回车符分隔
					  当type=customizedcast时，file内容为多条alias，
					    alias以回车符分隔，注意同一个文件内的alias所对应
					    的alias_type必须和接口参数alias_type一致。
					  注意，使用文件播前需要先调用文件上传接口获取file_id,
					     具体请参照"2.4文件上传接口"
	  "filter":{},            // 可选 终端用户筛选条件,如用户标签、地域、应用版本以及渠道等,
	  "payload":{},
	  "policy":{},
	  "production_mode":"true/false" // 可选 正式/测试模式。测试模式下，广播/组播只会将消息发给测试设备。测试设备需要到web上添加。Android: 测试设备属于正式设备的一个子集。
	  "description": "xx"      // 可选 发送消息描述，建议填写。
	  "thirdparty_id": "xx"    // 可选 开发者自定义消息标识ID, 开发者可以为同一批发送的多条消息
					   提供同一个thirdparty_id, 便于友盟后台后期合并统计数据。
}
*/
type Data struct {
  Platform       Platform    `json:"-"`
  AppKey         string      `json:"appkey,omitempty"`
  TimeStamp      int64       `json:"timestamp,omitempty"`
  TaskId         string      `json:"task_id,omitempty"`
  FileContent    string      `json:"content,omitempty"`
  Type           string      `json:"type,omitempty"`
  DeviceTokens   string      `json:"device_tokens,omitempty"`
  AliasType      string      `json:"alias_type,omitempty"`
  Alias          string      `json:"alias,omitempty"`
  FileId         string      `json:"file_id,omitempty"`
  Filter         string      `json:"filter,omitempty"`
  Payload        interface{} `json:"payload,omitempty"`
  Policy         Policy      `json:"policy,omitempty"`
  ProductionMode bool        `json:"production_mode"`
  Description    string      `json:"description,omitempty"`
  ThirdPartyId   string      `json:"thirdparty_id,omitempty"`
}

type Result struct {
  Code string `json:"ret,omitempty"`
  Data map[string]string
}

func NewData(pf Platform) (data *Data) {
  data = new(Data)
  data.Platform = pf
  if data.Platform == APP_ANDROID {
    data.AppKey = AndroidAppKey
  }
  if data.Platform == APP_IOS {
    data.AppKey = IOSAppKey
  }

  return
}

func (data *Data) SetPolicy(policy Policy) {
  data.Policy = policy
}

func (data *Data) Push(body, aps, policy interface{}, extras map[string]string) (result Result) {
  if data.Platform == APP_ANDROID {
    payload := &AndroidPayload{}
    if v, ok := body.(AndroidBody); ok {
      if v.DisplayType == "message" && len(v.Custom) == 0 {
        return
      }
      payload.DisplayType = v.DisplayType
      payload.Body = v
    }

    if len(extras) > 0 {
      payload.Extra = extras
    }

    data.Payload = payload

  } else if data.Platform == APP_IOS {
    payload := make(IOSPayload, 0)
    payload["aps"] = aps
    if len(extras) > 0 {
      for key, val := range extras {
        payload[key] = val
      }

    }

    data.Payload = payload
  }

  if v, ok := policy.(Policy); ok {
    data.SetPolicy(v)
  }

  result = data.Send(data.Sign(PostPath))
  return
}

func (data *Data) Status() (result Result) {
  result = data.Send(data.Sign(StatusPath))
  return
}

func (data *Data) Cancel() (result Result) {
  result = data.Send(data.Sign(CancelPath))
  return
}

func (data *Data) Upload() (result Result) {
  result = data.Send(data.Sign(UploadPath))
  return
}

func (data *Data) Sign(reqPath string) (api string) {
  DataByte, _ = jsonFast.Marshal(data)
  jsonStr := string(DataByte)
  sign := ""
  if data.Platform == APP_ANDROID {
    sign = Md5(fmt.Sprintf("POST%s%s%s%s", Host, reqPath, jsonStr, AndroidAppMasterSecret))
  } else if data.Platform == APP_IOS {
    sign = Md5(fmt.Sprintf("POST%s%s%s%s", Host, reqPath, jsonStr, IOSAppMasterSecret))
  }

  api = fmt.Sprintf("%s%s?sign=%s", Host, reqPath, sign)
  return
}

func (data *Data) Send(url string) (result Result) {
  fmt.Printf("\n data:%s", string(DataByte))
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(DataByte))
  if err != nil {
    return
  }

  req.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  jsonFast.Unmarshal(body, &result)
  return
}

func Md5(s string) string {
  h := md5.New()
  h.Write([]byte(s))
  return hex.EncodeToString(h.Sum(nil))
}
