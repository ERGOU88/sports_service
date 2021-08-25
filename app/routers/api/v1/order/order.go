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
	param := &morder.OrderRefund{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("order_trace: invalid param, param:%+v, err:%s", param, err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	userId, _ := c.Get(consts.USER_ID)
	svc := corder.New(c)
	param.UserId = userId.(string)
	code := svc.OrderRefund(param)
	reply.Response(http.StatusOK, code)
}
