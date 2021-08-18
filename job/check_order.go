package job

import (
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/global/rdskey"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/morder"
	"time"
	"errors"
)

// 检测订单支付是否超时（每分钟）
func CheckOrder() {
	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			log.Log.Debugf("开始检测订单支付是否超时")
			checkOrderTimeOut()
			log.Log.Debugf("检测完毕")
		}
	}
}

func checkOrderTimeOut() {
	orderIds, err := GetOrderIds()
	if err != nil {
		log.Log.Errorf("orderJob_trace: get orderIds fail, err:%s", err)
		return
	}

	for _, orderId := range orderIds {
		if err := orderTimeOut(orderId); err != nil {
			log.Log.Errorf("orderJob_trace: orderTimeOut fail, err:%s", err)
		}
	}
}

// 获取需处理超时的订单号
func GetOrderIds() ([]string, error) {
	rds := dao.NewRedisDao()
	return rds.SMEMBERS(rdskey.ORDER_EXPIRE_INFO)
}

// 超时处理完毕 删除缓存中的订单号
func DelOrderId(orderId string) (int, error) {
	rds := dao.NewRedisDao()
	return rds.Del(rdskey.ORDER_EXPIRE_INFO, orderId)
}

// 订单超时
func orderTimeOut(orderId string) error {
	session := dao.VenueEngine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		log.Log.Errorf("orderJob_trace: session begin err:%s, orderId:%s", err, orderId)
		return err
	}

	orderModel := morder.NewOrderModel(session)
	ok, err := orderModel.GetOrder(orderId)
	if !ok || err != nil {
		log.Log.Errorf("orderJob_trace: get order info fail, err:%s, ok:%v, orderId:%s", err, ok, orderId)
		session.Rollback()
		return nil
	}

	// 订单状态 != 0 (待支付) 表示 订单 已设为超时/已支付/已完成 等等...
	if orderModel.Order.Status != consts.PAY_TYPE_WAIT {
		log.Log.Errorf("orderJob_trace: don't need to change，orderId:%s, status:%d", orderId,
			orderModel.Order.Status)
		session.Rollback()
		return nil
	}

	now := int(time.Now().Unix())
	orderModel.Order.UpdateAt = now
	orderModel.Order.Status = consts.PAY_TYPE_UNPAID
	// 更新订单状态为 超时未支付
	affected, err := orderModel.UpdateOrderStatus(orderId, consts.PAY_TYPE_WAIT)
	if affected != 1 || err != nil {
		log.Log.Errorf("orderJob_trace: update order status fail, orderId:%s, err:%s", orderId, err)
		session.Rollback()
		return errors.New("update order status fail")
	}

	orderModel.OrderProduct.Status = consts.PAY_TYPE_UNPAID
	orderModel.OrderProduct.UpdateAt = now
	// 更新订单商品流水状态
	if _, err = orderModel.UpdateOrderProductStatus(orderId, consts.PAY_TYPE_WAIT); err != nil {
		log.Log.Errorf("orderJob_trace: update order product status fail, err:%s, affected:%d, orderId:%s", err, affected, orderId)
		session.Rollback()
		return errors.New("update order product status fail")
	}

	// 获取订单对应的预约流水
	amodel := mappointment.NewAppointmentModel(session)
	list, err := amodel.GetAppointmentRecordByOrderId(orderId, consts.PAY_TYPE_WAIT)
	if err != nil {
		log.Log.Errorf("orderJob_trace: get appointment record by orderId fail, orderId:%s, err:%s", orderId, err)
		session.Rollback()
		return err
	}

	for _, record := range list {
		// 归还对应节点的冻结库存
		affected, err = amodel.RevertStockNum(record.TimeNode, record.Date,  record.PurchasedNum * -1, now,
			record.AppointmentType, int(record.RelatedId))
		if affected != 1 || err != nil {
			log.Log.Errorf("orderJob_trace: update stock info fail, orderId:%s, err:%s, affected:%d, id:%d", orderId, err, affected, record.Id)
			session.Rollback()
			return errors.New("update stock info fail")
		}
	}

	// 更新标签状态[废弃]
	amodel.Labels.Status = 1
	if _, err = amodel.UpdateLabelsStatus(orderId, 0); err != nil {
		log.Log.Errorf("orderJob_trace: update labels status fail, orderId:%s, err:%s", orderId, err)
		session.Rollback()
		return errors.New("update label status fail")
	}

	if _, err := DelOrderId(orderId); err != nil {
		log.Log.Errorf("orderJob_trace: del orderId fail, err:%s", err)
		session.Rollback()
		return err
	}

	session.Commit()

	return nil
}
