//@Description todo
//@Author 凌云（何玉福）heyf@youzu.com
//@DateTime 2020/9/3 5:54 下午

package notify

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"sports_service/tools/goNotify/util"
	"strconv"
	"strings"
	"time"
)

const (
	FEISHU_URI = "https://dingcenter.youzu.com/api/v2/message/batch_send"
)

type Feishu struct {
	content
	AppId  string
	AppKey string
}

func (f Feishu) Send() error {
	if f.AppKey == "" || f.AppId == "" {
		f.AppId = "36"
		f.AppKey = "rK4YokBY8MMEejrDvctlMxo45I0Jksi2"
	}
	t := time.Now().Unix()
	header := http.Header{}
	header.Set("accept", "application/json")
	header.Set("appId", f.AppId)
	header.Set("time", strconv.FormatInt(t, 10))
	header.Set("Content-Type", "application/json")
	header.Set("token", f.token(t))
	query := make(map[string]interface{})
	query["msg_type"] = "text"
	query["accounts"] = strings.Split(f.To, ",")
	query["text"] = string(f.Content)
	body, _ := json.Marshal(query)
	resp, status, err := util.HttpDo(FEISHU_URI, http.MethodPost, f.ServiceName, body, 3, header)
	if err != nil || status != 200 {
		return errors.New("系统错误，飞书信息发送失败，err:" + err.Error())
	}
	code := gjson.GetBytes(resp, "code").Int()
	if code != 1 {
		return errors.New("飞书信息发送失败,err:" + gjson.GetBytes(resp, "msg").String())
	}
	return nil

}

func (f Feishu) token(time int64) string {
	str := fmt.Sprintf("%d%s", time, f.AppKey)
	sum := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", sum)
}
