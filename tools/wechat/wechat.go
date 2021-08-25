package wechat

import (
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/wechat"
	"sports_service/server/global/app/log"
	"strconv"
	"time"
	"errors"
)

const (
	MERCHANT_ID    = "1610121931"
	WECHAT_APP_ID  = "wxd693805bd4a2a39e"
	WECHAT_SECRET  = "UyPTFdVsJPPxMYzfXBztTKwsUusxSIFw"
)

type WechatPayClient struct {
	Client       *wechat.Client
	OutTradeNo   string             // 订单号
	TotalAmount  int                // 金额（分）
	Subject      string             // 商品名称
	CreateIp     string             // 请求ip
	NotifyUrl    string             // 付款回调地址
	TimeStart    string             // 交易生成时间
	TimeExpire   string             // 交易结束时间
	RefundAmount int                // 退款金额
	RefundNotify string             // 退款回调地址
}

// 初始化微信客户端
// appId：应用ID
// mchId：商户ID
// apiKey：API秘钥值
// isProd：是否是正式环境
func NewWechatPay(isProd bool) *WechatPayClient {
	wx := &WechatPayClient{}
	wx.Client = wechat.NewClient(WECHAT_APP_ID, MERCHANT_ID, WECHAT_SECRET, isProd)
	// 设置国家
	wx.Client.SetCountry(wechat.China)
	return wx
}

// 微信app支付
func (c *WechatPayClient) TradeAppPay() (map[string]interface{}, error){
	mp := make(map[string]interface{}, 0)
	// 初始化参数Map
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.GetRandomString(32)).
		Set("body", c.Subject).
		Set("out_trade_no", c.OutTradeNo).
		Set("total_fee", c.TotalAmount).
		Set("spbill_create_ip", c.CreateIp).
		Set("notify_url", c.NotifyUrl).
		Set("trade_type", wechat.TradeType_App).
		Set("sign_type", wechat.SignType_MD5)

	if c.TimeStart != "" {
		bm.Set("time_start", c.TimeStart)
	}

	if c.TimeExpire != "" {
		bm.Set("time_expire", c.TimeExpire)
	}

	// 请求支付下单，成功后得到结果
	wxRsp, err := c.Client.UnifiedOrder(bm)
	if err != nil {
		return nil, err
	}

	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	// 获取App支付需要的paySign
	paySign := wechat.GetAppPaySign(WECHAT_APP_ID, MERCHANT_ID, wxRsp.NonceStr, wxRsp.PrepayId, wechat.SignType_MD5,
		timeStamp, WECHAT_SECRET)

	mp["partner_id"] = MERCHANT_ID
	mp["sign"] = paySign
	mp["time_stamp"] = timeStamp
	mp["nonce_str"] = wxRsp.NonceStr
	mp["prepay_id"] = wxRsp.PrepayId
	mp["pkg_name"] = "Sign=WXPay"

	return mp, nil
}

// 校验签名
func (c *WechatPayClient) VerifySign(body interface{}) (bool, error) {
	ok, err := wechat.VerifySign(WECHAT_SECRET, wechat.SignType_MD5, body)
	if !ok || err != nil {
		return false, err
	}

	return true, nil
}

// 微信退款
func (c *WechatPayClient) TradeRefund() (*wechat.RefundResponse, error) {
	// 初始化参数结构体
	bm := make(gopay.BodyMap)
	bm.Set("out_trade_no", c.OutTradeNo).
		Set("nonce_str", util.GetRandomString(32)).
		Set("sign_type", wechat.SignType_MD5).
		Set("out_refund_no", util.GetRandomString(64)).
		Set("total_fee", c.TotalAmount).
		Set("refund_fee", c.RefundAmount).
		Set("notify_url", c.RefundNotify)

	wxRsp, _, err := c.Client.Refund(bm)
	if err != nil {
		return nil, err
	}

	log.Log.Info("wxRsp:", *wxRsp)

	if wxRsp.ReturnCode != "SUCCESS" || wxRsp.ResultCode != "SUCCESS" {
		return nil, errors.New("wx refund fail")
	}

	return wxRsp, nil
}
