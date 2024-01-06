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
	"sports_service/util"
	"strings"
	"time"
)

var (
	SMS_URI_NEW = "http://notify.youzu.com/api/sp/sendSMS"
	SMS_URI     = "http://sms.uuzuonline.com/api/sp/sendSMS"
)

type Sms struct {
	AppId          string
	AppKey         string
	Mobile         string // 手机号码
	TemplateCode   string // 短信模板编号
	TemplateParams string // 模板变量json形式，例如：{"code":1234,"username":"testUser"}
	Time           int64  // 请求时间戳：推送的时间戳
	ServiceName    string
}

const (
	APP_ID  = "Gb9cIfURgH"
	APP_KEY = "fO2txAmo5FOjgn4jAj"
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
	Status string `json:"status"`
	Desc   string `json:"desc"`
	//Data     struct {
	//  SpId    int  `json:"sp_id"`
	//} `json:"data"`
}

func (s Sms) doRequest() error {
	params, p, err := s.URLValues()
	if err != nil {
		return err
	}

	fmt.Printf("\n params:%s", params)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", SMS_URI_NEW, p.Encode()), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "youzu-go-notify")
	req.Header.Set("Service-VenueName", s.ServiceName)
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

	if res.Status != "0" {
		return errors.New("短信发送失败, err:" + res.Desc)
	}

	return nil
}
