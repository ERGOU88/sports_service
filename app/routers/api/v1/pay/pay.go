package pay

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
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
	"sports_service/server/tools/wechat"
	"sports_service/server/util"
	"strconv"
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
	code, payParam := svc.InitiatePayment(param)
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

	svc := corder.New(c)
	if b := svc.VerifySign(consts.ALIPAY, msg, sign, nil); !b {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if code := svc.AliPayNotify(params, string(body)); code != errdef.SUCCESS {
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}


type WXPayNotify struct {
	ReturnCode    string `json:"return_code"`
	ReturnMsg     string `json:"return_msg"`
	Appid         string `json:"appid"`
	MchID         string `json:"mch_id"`
	DeviceInfo    string `json:"device_info"`
	NonceStr      string `json:"nonce_str"`
	Sign          string `json:"sign"`
	ResultCode    string `json:"result_code"`
	ErrCode       string `json:"err_code"`
	ErrCodeDes    string `json:"err_code_des"`
	Openid        string `json:"openid"`
	IsSubscribe   string `json:"is_subscribe"`
	TradeType     string `json:"trade_type"`
	BankType      string `json:"bank_type"`
	TotalFee      int64  `json:"total_fee"`
	FeeType       string `json:"fee_type"`
	CashFee       int64  `json:"cash_fee"`
	CashFeeType   string `json:"cash_fee_type"`
	CouponFee     int64  `json:"coupon_fee"`
	CouponCount   int64  `json:"coupon_count"`
	CouponID0     string `json:"coupon_id_0"`
	CouponFee0    int64  `json:"coupon_fee_0"`
	TransactionID string `json:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no"`
	Attach        string `json:"attach"`
	TimeEnd       string `json:"time_end"`
}

// 微信回调通知
func WechatNotify(c *gin.Context) {
	bm, err := wxCli.ParseNotifyToBodyMap(c.Request)
	if err != nil {
		log.Log.Errorf("wxNotify_trace: parse notify to bodyMap fail, err:%s", err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	body, _ := util.JsonFast.Marshal(bm)
	log.Log.Debug("wxNotify_trace: body:%s, bm:%+v", string(body), bm)
	svc := corder.New(c)
	if b := svc.VerifySign(consts.WEICHAT, "", "", bm); !b {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if hasExist := util.MapExistBySlice(bm, []string{"return_code", "result_code", "out_trade_no", "total_fee",
		"time_end", "transaction_id"}); !hasExist {
		log.Log.Error("wxNotify_trace: map key not exists")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if bm["return_code"].(string) != "SUCCESS" || bm["result_code"].(string) != "SUCCESS" {
		log.Log.Errorf("wxNotify_trace: trade not success")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	orderId := bm["out_trade_no"].(string)

	order, err := svc.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Error("wxNotify_trace: order not found, orderId:%s, err:%s", orderId, err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	rsp := new(wxCli.NotifyResponse)
	rsp.ReturnCode = gopay.SUCCESS
	rsp.ReturnMsg = "OK"
	if order.Status != consts.ORDER_TYPE_WAIT {
		log.Log.Error("wxNotify_trace: order already pay, orderId:%s, err:%s", orderId, err)
		c.String(http.StatusOK, rsp.ToXmlString())
		return
	}

	totalFee := bm["total_fee"].(string)
	fee, err := strconv.Atoi(totalFee)
	if err != nil {
		log.Log.Error("wxNotify_trace: amount fail, orderId:%s, err:%s", orderId, err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if fee != order.Amount {
		log.Log.Error("wxNotify_trace: amount not match, orderAmount:%d, amount:%d", order.Amount, fee)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	payTime, _ := time.ParseInLocation("20060102150405", bm["time_end"].(string), time.Local)
	if err := svc.WechatPayNotify(orderId, string(body),  bm["transaction_id"].(string), totalFee,"", payTime.Unix(), 0, consts.PAY_NOTIFY); err != nil {
		log.Log.Errorf("wxNotify_trace: order process fail, err:%s", err)
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, rsp.ToXmlString())
}

// 微信退款回调
// 退款通知无sign，不用验签 只需解密数据
func WechatRefundNotify(c *gin.Context) {
	// 解析参数
	notifyReq, err := wxCli.ParseRefundNotify(c.Request)
	if err != nil {
		log.Log.Errorf("wxNotify_trace: parse refund notify fail, err:%s", err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	// 解密退款异步通知的加密数据
	refundNotify, err := wxCli.DecryptRefundNotifyReqInfo(notifyReq.ReqInfo, wechat.WECHAT_SECRET)
	if err != nil {
		log.Log.Errorf("wxNotify_trace: decrypt refund notify fail, err:%s", err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	body, _ := util.JsonFast.Marshal(refundNotify)
	log.Log.Debug("wxNotify_trace: body:%s, notify:%+v", string(body), refundNotify)

	if refundNotify.RefundStatus != "SUCCESS" || notifyReq.ReturnCode != "SUCCESS" {
		log.Log.Errorf("wxNotify_trace: trade not success")
		c.String(http.StatusBadRequest, "fail")
		return
	}

	orderId := refundNotify.OutTradeNo
	svc := corder.New(c)
	order, err := svc.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Error("wxNotify_trace: order not found, orderId:%s, err:%s", orderId, err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	rsp := new(wxCli.NotifyResponse)
	rsp.ReturnCode = gopay.SUCCESS
	rsp.ReturnMsg = "OK"
	if order.Status != consts.ORDER_TYPE_REFUND_WAIT {
		log.Log.Error("wxNotify_trace: order status fail, orderId:%s, status:%d, err:%s", orderId, order.Status, err)
		c.String(http.StatusOK, rsp.ToXmlString())
		return
	}

	refundFee := refundNotify.RefundFee
	refund, err := strconv.Atoi(refundFee)
	if err != nil {
		log.Log.Error("wxNotify_trace: refund amount fail, orderId:%s, err:%s", orderId, err)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if refund != order.RefundAmount {
		log.Log.Errorf("wxNotify_trace: refund amount not match, orderId:%s", orderId)
		c.String(http.StatusBadRequest, "fail")
		return
	}

	refundTm, _ := time.ParseInLocation("20060102150405", refundNotify.SuccessTime, time.Local)
	if err := svc.WechatPayNotify(orderId, string(body), refundNotify.TransactionId, refundNotify.TotalFee, refundNotify.OutRefundNo,
		int64(order.PayTime), refundTm.Unix(), consts.REFUND_NOTIFY); err != nil {
		log.Log.Errorf("wxNotify_trace: order process fail, err:%s", err)
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, rsp.ToXmlString())
}

func AppletPay(c *gin.Context) {
	reply := errdef.New(c)
	param := &morder.PayReqParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("pay_trace: invalid param, params:%+v", param)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	
	userId, _ := c.Get(consts.USER_ID)
	param.UserId = userId.(string)
	param.Platform = 1
	svc := cpay.New(c)
	code, payParam := svc.InitiatePayment(param)
	if code == errdef.SUCCESS {
		reply.Data["pay_param"] = payParam
	}
	
	reply.Response(http.StatusOK, code)
}
