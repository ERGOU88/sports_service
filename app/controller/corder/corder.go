package corder

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/app/config"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/mcoach"
	"sports_service/server/models/morder"
	"errors"
	"sports_service/server/global/app/log"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
	"sports_service/server/tools/alipay"
	alipayCli "github.com/go-pay/gopay/alipay"
	wxCli "github.com/go-pay/gopay/wechat"
	"sports_service/server/tools/wechat"
	"sports_service/server/util"
	"time"
	"fmt"
)

type OrderModule struct {
	context     *gin.Context
	engine      *xorm.Session
	order       *morder.OrderModel
	appointment *mappointment.AppointmentModel
	user        *muser.UserModel
	venue       *mvenue.VenueModel
	coach       *mcoach.CoachModel
}

func New(c *gin.Context) OrderModule {
	socket := dao.VenueEngine.NewSession()
	defer socket.Close()

	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()
	return OrderModule{
		context: c,
		order: morder.NewOrderModel(socket),
		appointment: mappointment.NewAppointmentModel(socket),
		user: muser.NewUserModel(appSocket),
		venue: mvenue.NewVenueModel(socket),
		coach: mcoach.NewCoachModel(socket),
		engine: socket,
	}
}

// 获取订单
func (svc *OrderModule) GetOrder(orderId string) (*models.VenuePayOrders, error)  {
	ok, err := svc.order.GetOrder(orderId)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("order not found")
	}

	return svc.order.Order, nil
}

// 订单处理流程
func (svc *OrderModule) AliPayNotify(orderId, body, status, tradeNo string, payTm int64, notifyType int) error {
	switch status {
	case consts.TradeSuccess:
		if err := svc.OrderProcess(orderId, body, tradeNo, payTm, notifyType); err != nil {
			return err
		}

	case consts.TradeClosed:
		log.Log.Debug("trade closed, order_id:%v", orderId)

	case consts.WaitBuyerPay:
		log.Log.Debug("wait buyer pay, order_id:%v", orderId)

	case consts.TradeFinished:
		log.Log.Debug("trade finished, order_id:%v", orderId)

	}

	return nil
}

// 订单处理流程 1 支付成功 2 退款流程 [用户申请退款成功时执行] 3 退款申请
func (svc *OrderModule) OrderProcess(orderId, body, tradeNo string, payTm int64, notifyType int) error {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("payNotify_trace: session begin fail, err:%s", err)
		return err
	}

	// tips: 不可直接更新状态 并发情况下会有问题
	// 订单当前状态 及 需更新的状态
	var curStatus, status int
	// 如果是支付成功回调 则订单当前状态应是 待支付 需更新状态为 已支付
	if notifyType == consts.PAY_NOTIFY {
		curStatus = consts.PAY_TYPE_WAIT
		status = consts.PAY_TYPE_PAID
	}

	// 如果是申请退款 则订单当前状态 应是已付款 需更新状态为 退款中
	if notifyType == consts.APPLY_REFUND {
		curStatus = consts.PAY_TYPE_PAID
		status = consts.PAY_TYPE_REFUND_WAIT
	}

	now := int(time.Now().Unix())
	svc.order.Order.Status = status
	svc.order.Order.IsCallback = 1
	svc.order.Order.PayTime = int(payTm)
	svc.order.Order.UpdateAt = now
	svc.order.Order.Transaction = tradeNo
	// 更新订单状态
	affected, err := svc.order.UpdateOrderStatus(orderId, curStatus)
	if affected != 1 || err != nil {
		log.Log.Errorf("payNotify_trace: update order status fail, orderId:%s", orderId)
		svc.engine.Rollback()
		return errors.New("update order status fail")
	}

	svc.order.OrderProduct.Status = status
	svc.order.OrderProduct.UpdateAt = now
	// 更新订单商品流水状态
	if _, err = svc.order.UpdateOrderProductStatus(orderId, curStatus); err != nil {
		log.Log.Errorf("payNotify_trace: update order product status fail, err:%s, orderId:%s", err, orderId)
		svc.engine.Rollback()
		return errors.New("update order product status fail")
	}

	switch svc.order.Order.ProductType {
	case consts.ORDER_TYPE_APPOINTMENT_VENUE, consts.ORDER_TYPE_APPOINTMENT_COACH, consts.ORDER_TYPE_APPOINTMENT_COURSE:
		// 更新订单对应的预约流水状态
		if err := svc.appointment.UpdateAppointmentRecordStatus(orderId, now, status, curStatus); err != nil {
			log.Log.Errorf("payNotify_trace: update order product status fail, err:%s, orderId:%s", err, orderId)
			svc.engine.Rollback()
			return err
		}

	case consts.ORDER_TYPE_MONTH_CARD, consts.ORDER_TYPE_SEANSON_CARD, consts.ORDER_TYPE_YEAR_CARD:
		ok, err := svc.order.GetOrderProductsById(orderId)
		if !ok || err != nil {
			log.Log.Errorf("payNotify_trace: get order products by id fail, orderId:%s, err:%s", orderId, err)
			svc.engine.Rollback()
			return errors.New("get order product fail")
		}

		// 更新会员可用时长 及 过期时长
		if err := svc.UpdateVipInfo(svc.user.User.UserId, svc.order.OrderProduct.RelatedId, now, svc.order.OrderProduct.Count,
			svc.order.OrderProduct.ExpireDuration, svc.order.OrderProduct.Duration, notifyType); err != nil {
			log.Log.Errorf("payNotify_trace: update vip info fail, orderId:%s, err:%s", orderId, err)
			svc.engine.Rollback()
			return err
		}
	}

	if notifyType != consts.APPLY_REFUND {
		// 记录回调信息
		if err := svc.RecordNotifyInfo(now, notifyType, orderId, body, tradeNo); err != nil {
			log.Log.Errorf("payNotify_trace: record notify info fail, orderId:%s, err:%s", orderId, err)
			svc.engine.Rollback()
			return err
		}
	}

	svc.engine.Commit()
	log.Log.Debug("payNotify_trace: 订单成功， orderId: %s", orderId)
	return nil
}

func (svc *OrderModule) RecordNotifyInfo(now, notifyType int, orderId, body, tradeNo string) error {
	svc.order.Notify.CreateAt = now
	svc.order.Notify.UpdateAt = now
	svc.order.Notify.NotifyType = notifyType
	svc.order.Notify.PayType = svc.order.Order.PayType
	svc.order.Notify.PayOrderId = orderId
	svc.order.Notify.NotifyInfo = body
	svc.order.Notify.Transaction = tradeNo
	// 记录回调信息
	affected, err := svc.order.AddOrderPayNotify()
	if affected != 1 || err != nil {
		log.Log.Errorf("payNotify_trace: record pay notify fail, orderId:%s, err:%s", orderId, err)
		return errors.New("record pay notify fail")
	}

	return nil
}

// 获取订单列表
func (svc *OrderModule) GetOrderList(userId, status string, page, size int) (int, []*morder.OrderInfo) {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS, []*morder.OrderInfo{}
	}

	var condition string
	offset := (page - 1) * size
	switch status {
	case "-1":
		// 查看全部 0 待支付 1 订单超时/未支付/已取消 2 已支付 3 已完成  4 退款中 5 已退款 6 已过期
		condition = fmt.Sprintf("order_type=1001 AND status >= 0 AND user_id=%s", userId)
	case "0":
		// 0 待支付
		condition = fmt.Sprintf("order_type=1001 AND status = 0 AND user_id=%s", userId)
	case "1":
		// 1 可使用
		condition = fmt.Sprintf("order_type=1001 AND status = 2 AND user_id=%s", userId)
	case "2":
		// 2 退款/售后 包含[3 已完成 4 退款中 5 已退款 6 已过期]
		condition = fmt.Sprintf("order_type=1001 AND status >= 3 AND user_id=%s", userId)
	default:
		log.Log.Errorf("order_trace: unsupported query status, status:%d", status)
		return errdef.ERROR, nil
	}

	list, err := svc.order.GetOrderListByStatus(condition, offset, size)
	if err != nil {
		log.Log.Errorf("order_trace: get order list by status fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*morder.OrderInfo{}
	}

	orderList := svc.OrderInfo(list)
	return errdef.SUCCESS, orderList
}

// todo:
func (svc *OrderModule) OrderInfo(list []*models.VenuePayOrders) []*morder.OrderInfo {
	res := make([]*morder.OrderInfo, len(list))
	for index, order := range list {
		info := new(morder.OrderInfo)
		info.OrderType = int32(order.ProductType)
		info.CreatAt = time.Unix(int64(order.CreateAt), 0).Format(consts.FORMAT_TM)
		info.OrderStatus = int32(order.Status)
		info.OrderId = order.PayOrderId
		info.UserId = order.UserId
		info.Amount = fmt.Sprintf("%.2f", float64(order.Amount)/100)
		info.Title = order.Subject
		extra := &mappointment.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(order.Extra, extra); err != nil {
			log.Log.Errorf("order_trace: unmarshal extra fail, err:%s, orderId:%s", err, order.PayOrderId)
			continue
		}

		switch order.OrderType {
		// 预约场馆、私教、大课
		case consts.ORDER_TYPE_APPOINTMENT_VENUE, consts.ORDER_TYPE_APPOINTMENT_COACH, consts.ORDER_TYPE_APPOINTMENT_COURSE:
			info.Count = len(extra.TimeNodeInfo)
			if order.ProductType == consts.ORDER_TYPE_APPOINTMENT_COACH && order.Status == consts.PAY_TYPE_PAID {
				// 查询是否评价
				ok, err := svc.coach.HasEvaluateByUserId(svc.order.Order.UserId, svc.order.Order.PayOrderId)
				if !ok || err != nil {
					log.Log.Errorf("order_trace: already evaluate, userId:%s, orderId:%s", svc.order.Order.UserId, svc.order.Order.PayOrderId)
				}

				info.HasEvaluate = ok
			}
		case consts.ORDER_TYPE_MONTH_CARD, consts.ORDER_TYPE_SEANSON_CARD, consts.ORDER_TYPE_YEAR_CARD:
			info.Count = extra.Count
		}

		res[index] = info
	}

	return res
}

// 更新会员信息
func (svc *OrderModule) UpdateVipInfo(userId string, venueId int64, now, count, expireDuration, duration, notifyType int) error {
	ok, err := svc.venue.GetVenueVipInfo(userId, venueId)
	if !ok || err != nil {
		log.Log.Errorf("order_trace: get venue vip info fail, userId:%s, err:%s", userId, err)
		return errors.New("vip not exists")
	}

	var cols string
	svc.venue.Vip.UpdateAt = now
	switch notifyType {
	// todo:月卡/季卡/年卡会员不可退款
	// 如果是申请退款 走退款流程
	//case consts.APPLY_REFUND:
	//	// 会员 需扣减可用时长  同时 扣减过期时长
	//	svc.venue.Vip.Duration += int64(duration * -1)
	//	svc.venue.Vip.EndTm = svc.venue.Vip.EndTm - int64(expireDuration * count)
	//	cols = "end_tm, duration, update_at"
	// 支付成功回调通知
	case consts.PAY_NOTIFY:
		// 如果vip结束时间 >= 当前时间戳 则为续费
		if int(svc.venue.Vip.EndTm) >= now {
			svc.venue.Vip.Duration += int64(duration)
			svc.venue.Vip.EndTm = int64(now + expireDuration * count)
			cols = "end_tm, duration, update_at"
		} else {
			// 否则 为 重新购买
			svc.venue.Vip.StartTm = int64(now)
			// 过期时间 叠加
			svc.venue.Vip.EndTm = int64(now + expireDuration * count)
			// 可用时长
			svc.venue.Vip.Duration = int64(duration)
			cols = "start_tm, end_tm, duration, update_at"
		}

	default:
		log.Log.Errorf("order_trace: unsupported notify type:%d", notifyType)
		return errors.New("unsupported notify type")
	}


	if _, err := svc.venue.UpdateVenueVipInfo(cols); err != nil {
		log.Log.Errorf("order_trace: update venue vip info err:%s", err)
		return err
	}

	return nil
}

// 订单详情
func (svc *OrderModule) OrderDetail(orderId, userId string) (int, *mappointment.OrderResp) {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("order_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS, nil
	}

	ok, err := svc.order.GetOrder(orderId)
	if !ok || err != nil {
		log.Log.Errorf("order_trace: order not exists, orderId:%s, err:%s", orderId, err)
		return errdef.ORDER_NOT_EXISTS, nil
	}

	// already delete
	if svc.order.Order.Status == -1 {
		log.Log.Errorf("order_trace: order already delete, orderId:%s", orderId)
		return errdef.ORDER_ALREADY_DEL, nil
	}

	rsp := &mappointment.OrderResp{}
	if err := util.JsonFast.UnmarshalFromString(svc.order.Order.Extra, rsp); err != nil {
		log.Log.Errorf("order_trace: unmarshal extra info fail, err:%s", err)
		return errdef.ERROR, nil
	}

	if svc.order.Order.Status == consts.PAY_TYPE_PAID && svc.order.Order.ProductType == consts.ORDER_TYPE_APPOINTMENT_COACH  {
		// 查询是否评价
		ok, err := svc.coach.HasEvaluateByUserId(svc.order.Order.UserId, svc.order.Order.PayOrderId)
		if !ok || err != nil {
			log.Log.Errorf("order_trace: already evaluate, userId:%s, orderId:%s", svc.order.Order.UserId, svc.order.Order.PayOrderId)
		}

		rsp.HasEvaluate = ok
	}

	rsp.PayDuration = 0
	// 待支付订单 剩余支付时长
	if svc.order.Order.Status == consts.PAY_TYPE_WAIT {
		// 已过时长 =  当前时间戳 - 订单创建时间戳
		duration := time.Now().Unix() - int64(svc.order.Order.CreateAt)
		// 订单状态是待支付 且 已过时长 <= 总时差
		if duration < consts.PAYMENT_DURATION {
			log.Log.Debugf("order_trace: duration:%v", duration)
			// 剩余支付时长 = 总时长[15分钟] - 已过时长
			rsp.PayDuration = consts.PAYMENT_DURATION - duration
		}
	}

	rsp.OrderId = orderId
	rsp.OrderStatus = int32(svc.order.Order.Status)
	rsp.OrderDescription = "订单需知"
	// 是否可退款
	rsp.CanRefund = svc.CanRefund(svc.order.Order.Status, svc.order.Order.OrderType, svc.order.Order.PayTime,
		svc.order.Order.PayOrderId, svc.order.Order.Extra)

	return errdef.SUCCESS, rsp
}

// 订单退款流程
func (svc *OrderModule) OrderRefund(param *morder.OrderRefund) int {
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR
	}

	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("order_trace: user not exists, userId:%s", param.UserId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}

	ok, err := svc.order.GetOrder(param.OrderId)
	if !ok || err != nil {
		log.Log.Errorf("order_trace: order not exists, orderId:%s", param.OrderId)
		svc.engine.Rollback()
		return errdef.ORDER_NOT_EXISTS
	}

	if svc.order.Order.UserId != user.UserId {
		log.Log.Errorf("order_trace: user not match, userId:%s, curUser:%s", svc.order.Order.UserId, user.UserId)
		svc.engine.Rollback()
		return errdef.ORDER_REFUND_FAIL
	}

	// 是否可退款
	if can := svc.CanRefund(svc.order.Order.Status, svc.order.Order.OrderType, svc.order.Order.PayTime,
		svc.order.Order.PayOrderId, svc.order.Order.Extra); !can {
		log.Log.Errorf("order_trace: user can't refund, orderId:%s", svc.order.Order.PayOrderId)
		svc.engine.Rollback()
		return errdef.ORDER_REFUND_FAIL
	}

	// 可退款
	if _, err := svc.TradeRefund(); err != nil {
		log.Log.Errorf("order_trace: trade refund err:%s", err)
		svc.engine.Rollback()
		return errdef.ORDER_REFUND_FAIL
	}

	if err := svc.OrderProcess(svc.order.Order.PayOrderId, "", svc.order.Order.Transaction,
		int64(svc.order.Order.PayTime), consts.APPLY_REFUND); err != nil {
		svc.engine.Rollback()
		return errdef.ORDER_REFUND_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 交易退款 todo:计算手续费
func (svc *OrderModule) TradeRefund() (string, error) {
	var body string
	switch svc.order.Order.PayType {
	case consts.ALIPAY:
		// 支付宝
		resp, err := svc.AliRefund()
		if err != nil {
			log.Log.Errorf("pay_trace: alipay refund fail, orderId:%s, payType:%d", svc.order.Order.PayOrderId,
				svc.order.Order.PayType)
			return "", err
		}

		log.Log.Infof("order_trace: orderId:%s, alipay refund resp:%+v", svc.order.Order.PayOrderId, resp)
		body, _ = util.JsonFast.MarshalToString(resp)

	case consts.WEICHAT:
		// 微信
		resp, err := svc.WechatRefund()
		if err != nil {
			log.Log.Errorf("pay_trace: get wechatPay param fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			return "", err
		}

		log.Log.Infof("order_trace: orderId:%s, wechat refund resp:%+v", svc.order.Order.PayOrderId, resp)
		body, _ = util.JsonFast.MarshalToString(resp)

	default:
		log.Log.Errorf("order_trace: unsupported payType:%d", svc.order.Order.PayType)
		return "", errors.New("unsupported payType")
	}


	return body, nil
}

// 支付宝退款
func (svc *OrderModule) AliRefund() (*alipayCli.TradeRefundResponse, error) {
	client := alipay.NewAliPay(true)
	client.OutTradeNo = svc.order.Order.PayOrderId
	client.RefundAmount = fmt.Sprintf("%.2f", float64(svc.order.Order.Amount)/100)
	client.RefundReson = fmt.Sprintf("%s%s", svc.order.Order.Subject, "退款")
	resp, err := client.TradeRefund()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 微信退款
func (svc *OrderModule) WechatRefund() (*wxCli.RefundResponse, error) {
	client := wechat.NewWechatPay(true)
	client.OutTradeNo = svc.order.Order.PayOrderId
	client.TotalAmount = svc.order.Order.Amount
	client.RefundAmount = svc.order.Order.Amount
	client.RefundNotify = config.Global.WechatRefundNotify
	resp, err := client.TradeRefund()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 是否可退款
func (svc *OrderModule) CanRefund(status, orderType, payTime int, orderId, extra string) bool {
	// 如果订单状态不等于已支付
	if status != consts.PAY_TYPE_PAID {
		return false
	}

	// 只有预约场馆/次卡/私教/课程可申请退款 同时需要判断 订单是否过期
	if orderType == consts.ORDER_TYPE_YEAR_CARD || orderType == consts.ORDER_TYPE_MONTH_CARD ||
		orderType == consts.ORDER_TYPE_SEANSON_CARD {
		return false
	}

	now := time.Now().Unix()
	switch orderType {
	case consts.ORDER_TYPE_APPOINTMENT_VENUE, consts.ORDER_TYPE_APPOINTMENT_COACH, consts.ORDER_TYPE_APPOINTMENT_COURSE:
		// 获取预约流水
		infos, err := svc.appointment.GetAppointmentRecordByOrderId(orderId, consts.PAY_TYPE_PAID)
		if len(infos) == 0 || err != nil {
			log.Log.Errorf("order_trace: get appointment record by orderId fail, orderId:%s, err:%s", orderId, err)
			return false
		}

		rsp := &mappointment.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(extra, rsp); err != nil {
			log.Log.Errorf("order_trace: unmarshal extra fail, orderId:%s, err:%s", orderId, err)
			return false
		}

		if len(rsp.TimeNodeInfo) == 0 {
			log.Log.Errorf("order_trace: time node is empty, orderId:%s", orderId)
			return false
		}

		for _, node := range rsp.TimeNodeInfo {
			// 如果预约中 某个节点的开始时间 <= 当前时间 表示已过期 不能退款
			if node.StartTm <= now {
				return false
			}
		}

	case consts.ORDER_TYPE_EXPERIENCE_CARD:
		// 获取订单商品
		ok, err := svc.order.GetOrderProductsById(orderId)
		if !ok || err != nil {
			log.Log.Errorf("order_trace: get order product by id fail, orderId:%s, err:%s", orderId, err)
			return false
		}

		// 次卡  订单完成时间 + 过期时长 <= 当前时间戳 表示已过期
		if payTime + svc.order.OrderProduct.ExpireDuration <= int(now) {
			return false
		}
	}


	return true
}
