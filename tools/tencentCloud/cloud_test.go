package tencentCloud

import (
	"sports_service/global/consts"
	"testing"
)

// 文本检测
func TestTextModeration(t *testing.T) {
	content := "你好 草尼玛"
	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	b, content, err := client.TextModeration(content)
	if err != nil {
		t.Errorf("text moderation err:%s", err)
		return
	}

	if !b {
		t.Logf("content:%s", content)

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
//func TestPullEvents(t *testing.T) {
//	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
//	events, err := client.PullEvents()
//	if err != nil {
//		t.Errorf("pull events err:%v", err)
//		return
//	}
//
//	t.Logf("pull events: %+v", events)
//}

// 确认事件回调
//func TestConfirmEvent(t *testing.T) {
//	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
//	if err := client.ConfirmEvents([]string{"test"}); err != nil {
//		t.Errorf("confirm events err:%v", err)
//		return
//	}
//}

// 上传
//func TestUpload(t *testing.T) {
//	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
//	err, _ := client.Upload(123, "321", "", "./test.mp4", "ap-shanghai", consts.VOD_PROCEDURE_NAME)
//	if err != nil {
//		t.Errorf("upload err:%v", err)
//		return
//	}
//}

//func TestFreeLogin(t *testing.T) {
//	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
//	mobileNum, err := client.FreeLogin("testtest", "mobile", "86")
//	t.Logf("mobileNum:%s, err:%s", mobileNum, err)
//}

//func TestRecognitionImage(t *testing.T) {
//	client := New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, "")
//	res, err := client.RecognitionImage("https://fpv-1253904687.cos.ap-shanghai.myqcloud.com/fpv/", "qweqweqwe.png")
//	t.Logf("res:%+v, err:%v", res, err)
//}
