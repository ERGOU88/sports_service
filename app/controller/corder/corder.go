package corder

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/morder"
	"errors"
	"sports_service/server/global/app/log"
	"time"
)

type OrderModule struct {
	context     *gin.Context
	engine      *xorm.Session
	order       *morder.OrderModel
	appointment *mappointment.AppointmentModel
}

func New(c *gin.Context) OrderModule {
	socket := dao.VenueEngine.NewSession()
	defer socket.Close()
	return OrderModule{
		context: c,
		order: morder.NewOrderModel(socket),
		appointment: mappointment.NewAppointmentModel(socket),
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
func (svc *OrderModule) AliPayNotify(orderId, body, status string, payTm int64) error {
	switch status {
	case consts.TradeSuccess:
		if err := svc.OrderProcess(orderId, body, payTm); err != nil {
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
func (svc *OrderModule) OrderProcess(orderId, body string, payTm int64) error {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("payNotify_trace: session begin fail, err:%s", err)
		return err
	}

	now := int(time.Now().Unix())
	svc.order.Order.Status = consts.PAY_TYPE_PAID
	svc.order.Order.IsCallback = 1
	svc.order.Order.PayTime = int(payTm)
	svc.order.Order.UpdateAt = now
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

	// 更新订单对应的预约流水状态
	if err := svc.appointment.UpdateAppointmentRecordStatus(orderId, now, consts.PAY_TYPE_PAID, consts.PAY_TYPE_WAIT); err != nil {
		log.Log.Errorf("payNotify_trace: update order product status fail, err:%s, orderId:%s", err, orderId)
		svc.engine.Rollback()
		return err
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
