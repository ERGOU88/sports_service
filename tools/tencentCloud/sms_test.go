package tencentCloud

import (
	"sports_service/server/global/consts"
	"testing"
)

// 发送短信
func TestSendSms(t *testing.T) {
	client := New(consts.TX_SMS_SECRET_ID, consts.TX_SMS_SECRET_KEY, consts.TMS_API_DOMAIN)
	b, err := client.SendSms("13177656222", "1001")
	if err != nil {
		t.Errorf("text moderation err:%s", err)
		return
	}

	t.Logf("text moderation pass:%v", b)
}
