package corder

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/morder"
	"errors"
	"sports_service/server/global/app/log"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
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
func (svc *OrderModule) AliPayNotify(orderId, body, status, tradeNo string, payTm int64) error {
	switch status {
	case consts.TradeSuccess:
		if err := svc.OrderProcess(orderId, body, tradeNo, payTm); err != nil {
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

// 订单处理流程
func (svc *OrderModule) OrderProcess(orderId, body, tradeNo string, payTm int64) error {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("payNotify_trace: session begin fail, err:%s", err)
		return err
	}

	now := int(time.Now().Unix())
	svc.order.Order.Status = consts.PAY_TYPE_PAID
	svc.order.Order.IsCallback = 1
	svc.order.Order.PayTime = int(payTm)
	svc.order.Order.UpdateAt = now
	svc.order.Order.Transaction = tradeNo
	// 更新订单状态
	affected, err := svc.order.UpdateOrderStatus(orderId, consts.PAY_TYPE_WAIT)
	if affected != 1 || err != nil {
		log.Log.Errorf("payNotify_trace: update order status fail, orderId:%s", orderId)
		svc.engine.Rollback()
		return errors.New("update order status fail")
	}

	svc.order.OrderProduct.Status = consts.PAY_TYPE_PAID
	svc.order.OrderProduct.UpdateAt = now
	// 更新订单商品流水状态为已支付
	if _, err = svc.order.UpdateOrderProductStatus(orderId, consts.PAY_TYPE_WAIT); err != nil {
		log.Log.Errorf("payNotify_trace: update order product status fail, err:%s, orderId:%s", err, orderId)
		svc.engine.Rollback()
		return errors.New("update order product status fail")
	}

	switch svc.order.Order.ProductType {
	case consts.ORDER_TYPE_APPOINTMENT_VENUE,consts.ORDER_TYPE_APPOINTMENT_COACH,consts.ORDER_TYPE_APPOINTMENT_COURSE:
		// 更新订单对应的预约流水状态
		if err := svc.appointment.UpdateAppointmentRecordStatus(orderId, now, consts.PAY_TYPE_PAID, consts.PAY_TYPE_WAIT); err != nil {
			log.Log.Errorf("payNotify_trace: update order product status fail, err:%s, orderId:%s", err, orderId)
			svc.engine.Rollback()
			return err
		}
	case consts.ORDER_TYPE_MONTH_CARD, consts.ORDER_TYPE_SEANSON_CARD, consts.ORDER_TYPE_YEAR_CARD:

	}

	svc.order.Notify.CreateAt = now
	svc.order.Notify.UpdateAt = now
	svc.order.Notify.PayType = svc.order.Order.PayType
	svc.order.Notify.PayOrderId = orderId
	svc.order.Notify.NotifyInfo = body
	// 记录回调信息
	affected, err = svc.order.AddOrderPayNotify()
	if affected != 1 || err != nil {
		log.Log.Errorf("payNotify_trace: record pay notify fail, err:%s", err)
		svc.engine.Rollback()
		return errors.New("record pay notify fail")
	}

	svc.engine.Commit()
	log.Log.Debug("payNotify_trace: 订单成功， orderId: %s", orderId)
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
		switch order.OrderType {
		// 预约场馆、私教、大课
		case consts.ORDER_TYPE_APPOINTMENT_VENUE, consts.ORDER_TYPE_APPOINTMENT_COACH, consts.ORDER_TYPE_APPOINTMENT_COURSE:
			extra := &mappointment.OrderResp{}
			if err := util.JsonFast.UnmarshalFromString(order.Extra, extra); err != nil {
				log.Log.Errorf("order_trace: unmarshal extra fail, err:%s, orderId:%s", err, order.PayOrderId)
				continue
			}

			info.Title = order.Subject
			info.Count = len(extra.TimeNodeInfo)
		}

		res[index] = info
	}

	return res
}

// 更新会员信息
func (svc *OrderModule) UpdateVipInfo(userId string, venueId int64, now, count, expireDuration, duration int) error {
	ok, err := svc.venue.GetVenueVipInfo(userId, venueId)
	if !ok || err != nil {
		log.Log.Errorf("venue_trace: get venue vip info fail, userId:%s, err:%s", userId, err)
		return errors.New("vip not exists")
	}

	var cols string
	svc.venue.Vip.UpdateAt = now
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

	// 待支付订单 剩余支付时长
	if svc.order.Order.Status == consts.PAY_TYPE_WAIT {
		// 已过时长 =  当前时间戳 - 订单创建时间戳
		duration := time.Now().Unix() - int64(svc.order.Order.CreateAt)
		// 订单状态是待支付 且 已过时长 <= 总时差
		if svc.order.Order.Status == consts.PAY_TYPE_WAIT && duration < consts.PAYMENT_DURATION {
			log.Log.Debugf("order_trace: duration:%v", duration)
			// 剩余支付时长 = 总时长[15分钟] - 已过时长
			rsp.PayDuration = consts.PAYMENT_DURATION - duration
		}
	}

	rsp.OrderId = orderId
	//rsp.OrderStatus = svc.order.Order.Status

	return errdef.SUCCESS, rsp
}
