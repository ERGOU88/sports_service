package job

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"net/http/httptest"
	"sports_service/app/controller/corder"
	"sports_service/models/mappointment"
	//"github.com/go-xorm/xorm"
	"sports_service/dao"
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/global/rdskey"
	"sports_service/models/morder"
	"time"
)

// 检测订单支付是否超时（30秒）
func CheckOrder() {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkOrderTimeOut()
			checkOrderExpire()
		}
	}
}

func checkOrderExpire() {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	svc := corder.New(c)
	if err := svc.CheckOrderExpire(); err != nil {
		log.Log.Errorf("orderJob_trace: check order expire fail, err:%s", err)
	}
}

func checkOrderTimeOut() {
	orderIds, err := GetOrderIds()
	if err != nil {
		log.Log.Errorf("orderJob_trace: get orderIds fail, err:%s", err)
		return
	}

	if len(orderIds) == 0 {
		return
	}

	for _, orderId := range orderIds {
		if err := orderTimeOut(orderId); err != nil {
			log.Log.Errorf("orderJob_trace: orderTimeOut fail, err:%s", err)
			continue
		}
	}
}

// 获取需处理超时的订单号
func GetOrderIds() ([]string, error) {
	rds := dao.NewRedisDao()
	return rds.SMEMBERS(rdskey.ORDER_EXPIRE_INFO)
}

// 超时处理完毕 / 订单已成功 删除缓存中的订单号
func DelOrderId(orderId string) (int, error) {
	rds := dao.NewRedisDao()
	return rds.SREM(rdskey.ORDER_EXPIRE_INFO, orderId)
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
		return errors.New("fail")
	}

	// 订单状态 != 0 (待支付) 表示 订单 已设为超时/已支付/已完成 等等...
	if orderModel.Order.Status != consts.ORDER_TYPE_WAIT {
		log.Log.Errorf("orderJob_trace: don't need to change，orderId:%s, status:%d", orderId,
			orderModel.Order.Status)
		DelOrderId(orderId)
		session.Rollback()
		return errors.New("fail")
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := int(time.Now().In(loc).Unix())
	// 如果当前时间 < 超时处理时间 不处理
	if now < orderModel.Order.CreateAt+consts.PAYMENT_DURATION {
		log.Log.Errorf("orderJob_trace: now < processTm, orderId:%s, now:%d, createAt:%d", orderId,
			now, orderModel.Order.CreateAt)
		session.Rollback()
		return errors.New("fail")
	}

	orderModel.Order.UpdateAt = now
	orderModel.Order.Status = consts.ORDER_TYPE_UNPAID
	// 更新订单状态为 超时未支付
	affected, err := orderModel.UpdateOrderStatus(orderId, consts.ORDER_TYPE_WAIT)
	if affected != 1 || err != nil {
		log.Log.Errorf("orderJob_trace: update order status fail, orderId:%s, err:%s", orderId, err)
		session.Rollback()
		return errors.New("update order status fail")
	}

	orderModel.OrderProduct.Status = consts.ORDER_TYPE_UNPAID
	orderModel.OrderProduct.UpdateAt = now
	// 更新订单商品流水状态
	if _, err = orderModel.UpdateOrderProductStatus(orderId, consts.ORDER_TYPE_WAIT); err != nil {
		log.Log.Errorf("orderJob_trace: update order product status fail, err:%s, affected:%d, orderId:%s", err, affected, orderId)
		session.Rollback()
		return errors.New("update order product status fail")
	}

	switch orderModel.Order.ProductType {
	// 预约类型的订单 需修改预约相关数据
	case consts.ORDER_TYPE_APPOINTMENT_VENUE, consts.ORDER_TYPE_APPOINTMENT_COACH, consts.ORDER_TYPE_APPOINTMENT_COURSE:
		if err := updateAppointmentInfo(session, orderId, now); err != nil {
			log.Log.Errorf("orderJob_trace: update appointment info fail, err:%s", err)
			session.Rollback()
			return err
		}
	}

	if _, err := DelOrderId(orderId); err != nil {
		log.Log.Errorf("orderJob_trace: del orderId fail, err:%s", err)
		session.Rollback()
		return err
	}

	log.Log.Errorf("orderJob_trace: del redis orderId success, orderId:%s", orderId)

	session.Commit()

	return nil
}

// 更新预约信息
func updateAppointmentInfo(session *xorm.Session, orderId string, now int) error {
	// 获取订单对应的预约流水
	amodel := mappointment.NewAppointmentModel(session)
	list, err := amodel.GetAppointmentRecordByOrderId(orderId)
	if err != nil {
		log.Log.Errorf("orderJob_trace: get appointment record by orderId fail, orderId:%s, err:%s", orderId, err)
		return err
	}

	for _, record := range list {
		switch record.AppointmentType {
		case consts.APPOINTMENT_VENUE:
			// 归还场馆预约对应节点的冻结库存
			affected, err := amodel.RevertStockNum(record.TimeNode, record.Date, record.PurchasedNum*-1, now,
				record.AppointmentType, int(record.VenueId))
			if affected != 1 || err != nil {
				log.Log.Errorf("orderJob_trace: update stock info fail, orderId:%s, err:%s, affected:%d, id:%d", orderId, err, affected, record.Id)
				return errors.New("update stock info fail")
			}

			// 归还抵扣的会员时长
			if record.DeductionTm > 0 {
				affected, err := amodel.UpdateVenueVipInfo(int(record.DeductionTm), record.VenueId, record.UserId)
				if affected != 1 || err != nil {
					log.Log.Errorf("order_trace: revert vip duration fail, orderId:%s, err:%s", record.PayOrderId, err)
					return err
				}

			}

			// 更新标签状态[废弃]
			amodel.Labels.Status = 1
			if _, err = amodel.UpdateLabelsStatus(orderId, 0); err != nil {
				log.Log.Errorf("orderJob_trace: update labels status fail, orderId:%s, err:%s", orderId, err)
				return errors.New("update label status fail")
			}

		case consts.APPOINTMENT_COACH, consts.APPOINTMENT_COURSE:
			// 归还课程对应节点的冻结库存
			affected, err := amodel.RevertCourseStockNum(record.TimeNode, record.Date, record.PurchasedNum*-1, now,
				record.AppointmentType, int(record.VenueId), int(record.CourseId), int(record.CoachId))
			if affected != 1 || err != nil {
				log.Log.Errorf("orderJob_trace: update stock info fail, orderId:%s, err:%s, affected:%d, id:%d", orderId, err, affected, record.Id)
				return errors.New("update stock info fail")
			}

		}

	}

	// 更新订单对应的预约流水状态
	//if err := amodel.UpdateAppointmentRecordStatus(orderId, now, 0); err != nil {
	//	log.Log.Errorf("payNotify_trace: update order product status fail, err:%s, orderId:%s", err, orderId)
	//	return err
	//}

	return nil
}
