package wechat

import (
	"testing"
	"time"
)

func TestTradeWapPay(t *testing.T) {
	client := NewWechatPay(true, APPLET_APPID, MERCHANT_ID, WECHAT_SECRET)
	client.TotalAmount = 1
	client.Subject = "测试商品"
	client.NotifyUrl = "https://www.baidu.com"
	client.CreateIp = "127.0.0.1"
	client.OutTradeNo = time.Now().Format("20060102150405")
	client.OpenId = "oWtAR40Z1Rrm-zd-1uonKTw_Z9Bg"
	mp, err := client.TradeJsAPIPay()
	if err != nil {
		t.Logf("err:%v", err)
		return
	}

	t.Logf("mp:%+v", mp)

	return
}
