package order

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/corder"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/morder"
	"sports_service/server/util"
)

// 订单列表
func OrderList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	userId, _ := c.Get(consts.USER_ID)
	status := c.Query("status")
	svc := corder.New(c)
	code, list := svc.GetOrderList(userId.(string), status, page, size)
	if code == errdef.SUCCESS {
		reply.Data["list"] = list
	}

	reply.Response(http.StatusOK, code)
}

// 订单详情
func OrderDetail(c *gin.Context) {
	reply := errdef.New(c)
	orderId := c.Query("order_id")
	userId, _ := c.Get(consts.USER_ID)
	svc := corder.New(c)
	code, detail := svc.OrderDetail(orderId, userId.(string))
	if code == errdef.SUCCESS {
		reply.Data["detail"] = detail
	}

	reply.Response(http.StatusOK, code)
}

// 订单退款
func OrderRefund(c *gin.Context) {
	reply := errdef.New(c)
	param := &morder.ChangeOrder{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("order_trace: invalid param, param:%+v, err:%s", param, err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	userId, _ := c.Get(consts.USER_ID)
	svc := corder.New(c)
	param.UserId = userId.(string)
	code, _, _, _ := svc.OrderRefund(param, consts.EXECUTE_TYPE_REFUND)
	reply.Response(http.StatusOK, code)
}

// 删除订单
func OrderDelete(c *gin.Context) {
	reply := errdef.New(c)
	param := &morder.ChangeOrder{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("order_trace: invalid param, param:%+v, err:%s", param, err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	userId, _ := c.Get(consts.USER_ID)
	svc := corder.New(c)
	param.UserId = userId.(string)
	code := svc.DeleteOrder(param)
	reply.Response(http.StatusOK, code)
}

// 查看券码[次卡/预约场馆]
func OrderCouponCode(c *gin.Context) {
	reply := errdef.New(c)
	orderId := c.Query("order_id")
	userId, _ := c.Get(consts.USER_ID)
	svc := corder.New(c)
	code, resp := svc.GetCouponCodeInfo(userId.(string), orderId)
	if code == errdef.SUCCESS {
		reply.Data["detail"] = resp
	}

	reply.Response(http.StatusOK, code)
}

// 取消订单
func OrderCancel(c *gin.Context) {
	reply := errdef.New(c)
	param := &morder.ChangeOrder{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("order_trace: invalid param, err:%s, param:%+v", err, param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	userId, _ := c.Get(consts.USER_ID)
	param.UserId = userId.(string)
	svc := corder.New(c)
	code := svc.OrderCancel(param)
	reply.Response(http.StatusOK, code)
}

// 退款规则
func RefundRules(c *gin.Context) {
	reply := errdef.New(c)
	orderId := c.Query("order_id")
	userId, _ := c.Get(consts.USER_ID)
	svc := corder.New(c)
	code, refundAmount, refundFee, ruleId := svc.OrderRefund(&morder.ChangeOrder{
		UserId: userId.(string),
		OrderId: orderId,
	}, consts.EXECUTE_TYPE_QUERY)

	if code != errdef.SUCCESS {
		reply.Response(http.StatusOK, code)
		return
	}

	code, rules := svc.OrderRefundRules()
	if code != errdef.SUCCESS {
		reply.Response(http.StatusOK, code)
		return
	}

	reply.Data["rule_id"] = ruleId
	reply.Data["list"] = rules
	reply.Data["refund_amount"] = refundAmount
	reply.Data["refund_fee"] = refundFee
	reply.Data["order_amount"] = refundAmount + refundFee

	reply.Response(http.StatusOK, code)
}
