package pay

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cpay"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/morder"
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

}

// 微信回调通知
func WechatNotify(c *gin.Context) {

}
