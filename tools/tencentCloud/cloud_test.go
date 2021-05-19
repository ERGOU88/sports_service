package tencentCloud

import (
	"sports_service/server/global/consts"
	"testing"
)

// 文本检测
func TestTextModeration(t *testing.T) {
	content := "尼玛"
	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	b, err := client.TextModeration(content)
	if err != nil {
		t.Errorf("text moderation err:%s", err)
		return
	}

	t.Logf("text moderation pass:%v", b)
}

// 生成签名
//func TestGenerateSign(t *testing.T) {
//	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
//	sign := client.GenerateSign("1234567", consts.VOD_PROCEDURE_NAME, 7654321)
//	t.Logf("upload sign: %s", sign)
//}

// 主动拉取事件回调
func TestPullEvents(t *testing.T) {
	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
	events, err := client.PullEvents()
	if err != nil {
		t.Errorf("pull events err:%v", err)
		return
	}

	t.Logf("pull events: %+v", events)
}

// 确认事件回调
//func TestConfirmEvent(t *testing.T) {
//	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
//	if err := client.ConfirmEvents([]string{"test"}); err != nil {
//		t.Errorf("confirm events err:%v", err)
//		return
//	}
//}

// 上传
func TestUpload(t *testing.T) {
	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
	err, _ := client.Upload(123, "321", "", "./test.mp4", "ap-shanghai", consts.VOD_PROCEDURE_NAME)
	if err != nil {
		t.Errorf("upload err:%v", err)
		return
	}
}

