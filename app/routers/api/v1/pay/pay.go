package pay

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	wxCli "github.com/go-pay/gopay/wechat"
	"io/ioutil"
	"net/http"
	"net/url"
	"sports_service/server/app/controller/corder"
	"sports_service/server/app/controller/cpay"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/morder"
	"sports_service/server/tools/alipay"
	"sports_service/server/tools/wechat"
	"strconv"
	"strings"
	"time"
)

// app发起支付
func AppPay(c *gin.Context) {
	reply := errdef.New(c)
	param := &morder.PayReqParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("pay_trace: invalid param, params:%+v", param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	userId, _ := c.Get(consts.USER_ID)
	param.UserId = userId.(string)
	svc := cpay.New(c)
	code, payParam := svc.AppPay(param)
	if code == errdef.SUCCESS {
		reply.Data["pay_param"] = payParam
	}

	reply.Response(http.StatusOK, code)
}

// 支付宝回调通知
func AliPayNotify(c *gin.Context) {
	req := c.Request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Log.Errorf("aliNotify_trace: err:%s", err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	log.Log.Debug("aliNotify_trace: info %s", string(body))
	params, err := url.ParseQuery(string(body))
	if err != nil {
		log.Log.Errorf("aliNotify_trace: err:%s, params:%v", err.Error(), params)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	sign := params.Get("sign")
	params.Del("sign")
	params.Del("sign_type")
	query := params.Encode()
	msg, err := url.QueryUnescape(query)
	if err != nil {
		log.Log.Error("aliNotify_trace: QueryUnescape failed: %s", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	sign, _ = url.PathUnescape(sign)
	log.Log.Debug("aliNotify_trace: msg:%s, sign:%v", msg, sign)

	orderId := params.Get("out_trade_no")
	svc := corder.New(c)
	order, err := svc.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Error("aliNotify_trace: order not found, orderId:%s, err:%s", orderId, err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if order.Status != consts.PAY_TYPE_WAIT {
		log.Log.Error("aliNotify_trace: order already pay, orderId:%s, err:%s", orderId, err)
		c.String(http.StatusOK, "success")
		return
	}

	cli := alipay.NewAliPay(true)
	ok, err := cli.VerifyData(msg, "RSA2", sign)
	if !ok || err != nil {
		log.Log.Errorf("aliNotify_trace: verify data fail, err:%s", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	amount, err := strconv.ParseFloat(strings.Trim(params.Get("total_amount"), " "), 64)
	if err != nil {
		log.Log.Errorf("aliNotify_trace: parse float fail, err:%s", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if int(amount * 100) != order.Amount {
		log.Log.Error("aliNotify_trace: amount not match, orderAmount:%d, amount:%d", order.Amount, amount * 100)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	status := strings.Trim(params.Get("trade_status"), " ")
	payTime, _ := time.Parse("2006-01-02 15:04:05", params.Get("gmt_payment"))
	tradeNo := params.Get("trade_no")
	if err := svc.AliPayNotify(orderId, string(body), status, tradeNo, payTime.Unix()); err != nil {
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}


type WXPayNotify struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	Appid         string `xml:"appid"`
	MchID         string `xml:"mch_id"`
	DeviceInfo    string `xml:"device_info"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	ResultCode    string `xml:"result_code"`
	ErrCode       string `xml:"err_code"`
	ErrCodeDes    string `xml:"err_code_des"`
	Openid        string `xml:"openid"`
	IsSubscribe   string `xml:"is_subscribe"`
	TradeType     string `xml:"trade_type"`
	BankType      string `xml:"bank_type"`
	TotalFee      int64  `xml:"total_fee"`
	FeeType       string `xml:"fee_type"`
	CashFee       int64  `xml:"cash_fee"`
	CashFeeType   string `xml:"cash_fee_type"`
	CouponFee     int64  `xml:"coupon_fee"`
	CouponCount   int64  `xml:"coupon_count"`
	CouponID0     string `xml:"coupon_id_0"`
	CouponFee0    int64  `xml:"coupon_fee_0"`
	TransactionID string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	Attach        string `xml:"attach"`
	TimeEnd       string `xml:"time_end"`
}

// 微信回调通知
func WechatNotify(c *gin.Context) {
	req := c.Request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Log.Errorf("wxNotify_trace: err:%s", err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	log.Log.Errorf("wxNotify_trace: body:%s", string(body))
	var wx WXPayNotify
	err = xml.Unmarshal(body, &wx)
	if err != nil {
		log.Log.Errorf("wxNotify_trace: xml unmarshal err:%s", err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	bm, err := wxCli.ParseNotifyToBodyMap(req)
	if err != nil {
		log.Log.Errorf("wxNotify_trace: parse notify to bodyMap fail, err:%s", err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	cli := wechat.NewWechatPay(true)
	ok, err := cli.VerifySign(bm)
	if !ok || err != nil {
		log.Log.Error("wxNotify_trace: sign not match, err:%s", err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	log.Log.Debug("wxNotify_trace: info %s", string(body))
	params, err := url.ParseQuery(string(body))
	if err != nil {
		log.Log.Errorf("wxNotify_trace: err:%s, params:%v", err.Error(), params)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if !(wx.ReturnCode == "SUCCESS" && wx.ResultCode == "SUCCESS") {
		log.Log.Errorf("wxNotify_trace: trade not success")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	orderId := wx.OutTradeNo
	svc := corder.New(c)
	order, err := svc.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Error("wxNotify_trace: order not found, orderId:%s, err:%s", orderId, err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if order.Status != consts.PAY_TYPE_WAIT {
		log.Log.Error("wxNotify_trace: order already pay, orderId:%s, err:%s", orderId, err)
		c.String(http.StatusOK, "SUCCESS")
		return
	}

	if int(wx.TotalFee) != order.Amount {
		log.Log.Error("wxNotify_trace: amount not match, orderAmount:%d, amount:%d", order.Amount, wx.TotalFee)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	payTime, _ := time.Parse("20060102150405", wx.TimeEnd)
	if err := svc.OrderProcess(orderId, string(body), wx.TransactionID, payTime.Unix()); err != nil {
		log.Log.Errorf("wxNotify_trace: order process fail, err:%s", err)
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "SUCCESS")
}
