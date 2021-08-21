package pay

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
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
	"sports_service/server/app/controller/corder"
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
	reply := errdef.New(c)
	req := c.Request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Log.Errorf("aliNotify_trace: err:%s", err.Error())
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	log.Log.Debug("aliNotify_trace: info %s", string(body))
	params, err := url.ParseQuery(string(body))
	if err != nil {
		log.Log.Errorf("aliNotify_trace: err:%s, params:%v", err.Error(), params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	sign := params.Get("sign")
	params.Del("sign")
	params.Del("sign_type")
	query := params.Encode()
	msg, err := url.QueryUnescape(query)
	if err != nil {
		log.Log.Error("aliNotify_trace: QueryUnescape failed: %s", err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	sign, _ = url.PathUnescape(sign)
	log.Log.Debug("aliNotify_trace: msg:%s, sign:%v", msg, sign)

	orderId := params.Get("out_trade_no")
	svc := corder.New(c)
	order, err := svc.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Error("aliNotify_trace: order not found, orderId:%s, err:%s", orderId, err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	if order.Status != consts.PAY_TYPE_WAIT {
		log.Log.Error("aliNotify_trace: order already pay, orderId:%s, err:%s", orderId, err)
		reply.Response(http.StatusOK, errdef.SUCCESS)
		return
	}

	cli := alipay.NewAliPay(true)
	ok, err := cli.VerifyData(msg, "RSA2", sign)
	if !ok || err != nil {
		log.Log.Errorf("aliNotify_trace: verify data fail, err:%s", err)
		reply.Response(http.StatusBadRequest, errdef.ERROR)
		return
	}

	amount, err := strconv.ParseFloat(strings.Trim(params.Get("total_amount"), " "), 64)
	if err != nil {
		log.Log.Errorf("aliNotify_trace: parse float fail, err:%s", err)
		reply.Response(http.StatusBadRequest, errdef.ERROR)
		return
	}

	if int(amount * 100) != order.Amount {
		log.Log.Error("aliNotify_trace: amount not match, orderAmount:%d, amount:%d", order.Amount, amount * 100)
		reply.Response(http.StatusBadRequest, errdef.ERROR)
		return
	}

	status := strings.Trim(params.Get("trade_status"), " ")
	payTime, _ := time.Parse("2006-01-02 15:04:05", params.Get("gmt_payment"))
	if err := svc.OrderProcess(string(body), status, payTime.Unix()); err != nil {
		reply.Response(http.StatusInternalServerError, errdef.ERROR)
		return
	}

	reply.Response(http.StatusOK, errdef.PAY_SUCCESS)
}

// 微信回调通知
func WechatNotify(c *gin.Context) {
	reply := errdef.New(c)
	req := c.Request
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Log.Errorf("wxNotify_trace: err:%s", err.Error())
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	log.Log.Errorf("wxNotify_trace: body:%s", string(body))

	cli := wechat.NewWechatPay(true)
	ok, err := cli.VerifySign(string(body))
	if !ok || err != nil {
		log.Log.Error("wxNotify_trace: sign not match")
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	log.Log.Debug("wxNotify_trace: info %s", string(body))
	params, err := url.ParseQuery(string(body))
	if err != nil {
		log.Log.Errorf("wxNotify_trace: err:%s, params:%v", err.Error(), params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	if params.Get("result_code") != "SUCCESS" {
		log.Log.Errorf("wxNotify_trace: trade not success")
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	orderId := params.Get("out_trade_no")
	svc := corder.New(c)
	order, err := svc.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Error("wxNotify_trace: order not found, orderId:%s, err:%s", orderId, err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	if order.Status != consts.PAY_TYPE_WAIT {
		log.Log.Error("wxNotify_trace: order already pay, orderId:%s, err:%s", orderId, err)
		reply.Response(http.StatusOK, errdef.SUCCESS)
		return
	}

	amount, err := strconv.Atoi(strings.Trim(params.Get("total_fee"), " "))
	if err != nil {
		log.Log.Errorf("wxNotify_trace: parse float fail, err:%s", err)
		reply.Response(http.StatusBadRequest, errdef.ERROR)
		return
	}

	if amount != order.Amount {
		log.Log.Error("wxNotify_trace: amount not match, orderAmount:%d, amount:%d", order.Amount, amount * 100)
		reply.Response(http.StatusBadRequest, errdef.ERROR)
		return
	}

	status := strings.Trim(params.Get("trade_status"), " ")
	payTime, _ := time.Parse("20060102150405", params.Get("time_end"))
	if err := svc.OrderProcess(string(body), status, payTime.Unix()); err != nil {
		reply.Response(http.StatusInternalServerError, errdef.ERROR)
		return
	}

	reply.Response(http.StatusOK, errdef.PAY_SUCCESS)
}
