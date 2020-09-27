//@Description todo
//@Author 凌云（何玉福）heyf@youzu.com
//@DateTime 2020/9/3 5:54 下午

package notify

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"sports_service/server/tools/goNotify/util"
	"time"
)

var (
	SMS_URI = "http://sms.uuzuonline.com/api/sp/sendSMS"
)

type Sms struct {
	content
	AppId  string
	AppKey string
}

func (s Sms) Send() error {
	if s.AppId == "" || s.AppKey == "" {
		s.AppId = "dhlQN9jYck"
		s.AppKey = "WKbPKSr8TuRfrfmzch"
	}
	t := time.Now().Unix()
	query := fmt.Sprintf("app_id=%s&content=%s&mobile=%s&time=%d", s.AppId, s.Content, s.To, t)
	query = fmt.Sprintf("%s&verify=%s", query, s.sign(query))
	header := http.Header{}
	header.Set("Content-Type", "application/www-form-urlencoude")
	resp, status, err := util.HttpDo(SMS_URI, http.MethodGet, s.ServiceName, []byte(query), 5, header)
	if err != nil || status != 200 {
		return err
	}
	code := gjson.GetBytes(resp, "status").String()
	if code != "0" {
		return errors.New("短信发送失败,err:" + gjson.GetBytes(resp, "desc").String())
	}
	return nil
}

func (s Sms) sign(query string) string {
	query += s.AppKey
	m := md5.Sum([]byte(query))
	return fmt.Sprintf("%x", m)
}
