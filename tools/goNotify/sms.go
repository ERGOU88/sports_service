//@Description todo
//@Author 凌云（何玉福）heyf@youzu.com
//@DateTime 2020/9/3 5:54 下午

package notify

import (
  "crypto/md5"
  "encoding/hex"
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "sort"
  "strings"
  "sports_service/server/util"
  "time"
)

var (
  SMS_URI = "http://sms.uuzuonline.com/api/sp/sendSMS"
)

//type Sms struct {
//	content
//	AppId  string
//	AppKey string
//}

type Sms struct {
  AppId           string
  AppKey          string
  Mobile          string      // 手机号码
  TemplateCode    string      // 短信模板编号
  TemplateParams  string      // 模板变量json形式，例如：{"code":1234,"username":"testUser"}
  Time            int64       // 请求时间戳：推送的时间戳
  ServiceName     string
}

const (
  APP_ID        = "dhlQN9jYck"
  APP_KEY       = "WKbPKSr8TuRfrfmzch"
)

func (s Sms) Send() error {
  if s.AppId == "" || s.AppKey == "" {
    s.AppId = APP_ID
    s.AppKey = APP_KEY
  }

  if err := s.doRequest(); err != nil {
    fmt.Println("\n err:", err)
    return err
  }

  return nil


  //t := time.Now().Unix()
  //query := fmt.Sprintf("app_id=%s&mobile=%s&template_code=%s&template_params=%s&time=%d", s.AppId, s.Mobile, s.TemplateCode, s.TemplateParams,  t)
  //query = fmt.Sprintf("%s&verify=%s", query, s.sign(query))
  //header := http.Header{}
  //header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
  //resp, status, err := util.HttpDo(SMS_URI, http.MethodGet, s.ServiceName, []byte(query), 5, header)
  //if err != nil || status != 200 {
  //  return err
  //}
  //code := gjson.GetBytes(resp, "status").String()
  //if code != "0" {
  //  return errors.New("短信发送失败,err:" + gjson.GetBytes(resp, "desc").String())
  //}
  //return nil
}

func (s Sms) sign(query string) string {
  query += s.AppKey
  m := md5.Sum([]byte(query))
  return fmt.Sprintf("%x", m)
}

func (s Sms) Params() map[string]string {
  var m = make(map[string]string)
  m["app_id"] = s.AppId
  m["mobile"] = s.Mobile
  m["template_code"] = s.TemplateCode
  m["template_params"] = s.TemplateParams
  m["time"] = fmt.Sprint(s.Time)
  return m
}

func (s Sms) URLValues() (data string, value url.Values, err error) {
  var p = url.Values{}
  var ps = s.Params()
  if ps != nil {
    for key, value := range ps {
      if value != "" {
        p.Add(key, value)
      }
    }
  }

  sv := s.GetKeysAndValuesBySortKeys(p)
  md5sum := md5.New()
  params := strings.Join(sv, "&")
  md5sum.Write([]byte(params))
  md5sum.Write([]byte(s.AppKey))

  sign := hex.EncodeToString(md5sum.Sum([]byte(nil)))

  p.Add("verify", sign)
  params = fmt.Sprintf("%s&verify=%s", params, sign)
  return params, p, nil
}

// 排序
func (s Sms) GetKeysAndValuesBySortKeys(urlValues url.Values) (values []string) {
  vLen := len(urlValues)
  // get keys
  keys := make([]string, vLen)
  i := 0
  for k := range urlValues {
    keys[i] = k
    i++
  }
  // sort keys
  sort.Sort(sort.StringSlice(keys))
  values = make([]string, vLen)
  for i, k := range keys {
    values[i] = fmt.Sprintf(`%s=%s`, k, urlValues.Get(k))
  }

  return
}

// 返回值
type Resp struct {
  Status   int     `json:"status"`
  Desc     string  `json:"desc"`
  Data     struct {
   SpId    int  `json:"sp_id"`
  } `json:"data"`
}

func (s Sms) doRequest() error {
  params, _, err := s.URLValues()
  if err != nil {
    return err
  }

  fmt.Printf("\n params:%s", params)

  req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", SMS_URI, params), nil)
  if err != nil {
    return err
  }
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  //req.Header.Set("User-Agent", "youzu-go-notify")
  //req.Header.Set("Service-Name", s.ServiceName)
  fmt.Printf("\nreq:%+v, \nurl:%s", req, req.URL)

  client := http.Client{
    Timeout: time.Second * time.Duration(3),
  }
  resp, err := client.Do(req)
  if resp != nil {
    defer resp.Body.Close()
  }

  if err != nil {
    return err
  }

  fmt.Printf("\nresp:%v", resp)
  if resp.StatusCode != 200 {
    return errors.New("request fail")
  }

  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }

  res := &Resp{}
  if err = util.JsonFast.Unmarshal(data, res); err != nil {
    return err
  }

  fmt.Printf("res:%+v", res)

  if res.Status != 0 {
    return errors.New("短信发送失败, err:" + res.Desc)
  }

  return nil
}
