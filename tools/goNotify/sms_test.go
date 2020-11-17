//@Description todo
//@Author 凌云（何玉福）heyf@youzu.com
//@DateTime 2020/9/3 7:38 下午

package notify

import (
  "encoding/json"
  "testing"
  "time"
)

type TemplateParams struct {
  Code       string    `json:"code"`
  Username   string    `json:"username"`
}

func TestSms_Send(t *testing.T) {
  params := &TemplateParams{
    Code:   "4321",
  }

  bts, _ := json.Marshal(params)

  s := &Sms{}
  s.Mobile = "13177656222"
  s.TemplateCode = "SMS_000042"
  s.TemplateParams = string(bts)
  s.Time = time.Now().Unix()
  s.ServiceName = "test"
  t.Log(s.Send())
}

