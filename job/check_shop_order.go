package job

import (
	"errors"
	"fmt"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/global/rdskey"
	"sports_service/server/models/mshop"
	"time"
)

// 检测商城订单 支付是否超时（5min）
func CheckShopOrder() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	
	for {
		select {
		case <- ticker.C:
			checkShopOrderTimeOut()
		}
	}
}

func checkShopOrderTimeOut() {
	orderIds, err := GetShopOrderIds()
	if err != nil {
		log.Log.Errorf("shopJob_trace: get orderIds fail, err:%s", err)
		return
	}
	
	if len(orderIds) == 0 {
		return
	}
	
	for _, orderId := range orderIds {
		if err := shopOrderTimeOut(orderId); err != nil {
			log.Log.Errorf("shopJob_trace: orderTimeOut fail, err:%s", err)
			continue
		}
	}
}

// 获取需处理超时的商城订单号
func GetShopOrderIds() ([]string, error) {
	rds := dao.NewRedisDao()
	return rds.SMEMBERS(rdskey.SHOP_ORDER_EXPIRE)
}

// 超时处理完毕 / 订单已成功 删除缓存中的商城订单号
func DelShopOrderId(orderId string) (int, error) {
	rds := dao.NewRedisDao()
	return rds.SREM(rdskey.SHOP_ORDER_EXPIRE, orderId)
}

// 商城订单超时
func shopOrderTimeOut(orderId string) error {
	session := dao.AppEngine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Log.Errorf("shopJob_trace: session begin err:%s, orderId:%s", err, orderId)
		return err
	}
	
	shopModel := mshop.NewShop(session)
	order, err := shopModel.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Errorf("shopJob_trace: get order info fail, err:%s, orderId:%s", err, orderId)
		session.Rollback()
		return errors.New("fail")
	}
	
	// 订单支付状态 != 0 (待支付) 表示 订单 已设为超时/已支付/已完成 等...
	if order.PayStatus != consts.SHOP_ORDER_TYPE_WAIT {
		log.Log.Errorf("shopJob_trace: don't need to change，orderId:%s, status:%d", orderId, order.PayStatus)
		DelShopOrderId(orderId)
		session.Rollback()
		return errors.New("fail")
	}
	
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := int(time.Now().In(loc).Unix())
	// 如果当前时间 < 超时处理时间 不处理
	if now < order.CreateAt + consts.SHOP_PAYMENT_DURATION {
		log.Log.Errorf("shopJob_trace: now < processTm, orderId:%s, now:%d, createAt:%d", orderId,
			now, order.CreateAt)
		session.Rollback()
		return errors.New("fail")
	}
	
	condition := fmt.Sprintf("`pay_status`=%d AND `order_id`='%s'", consts.SHOP_ORDER_TYPE_WAIT, order.OrderId)
	cols := "pay_status, close_time, update_at"
	order.PayStatus = consts.SHOP_ORDER_TYPE_UNPAID
	order.UpdateAt = now
	order.CloseTime = now
	
	// 更新订单状态为 超时未支付
	affected, err := shopModel.UpdateOrderInfo(condition, cols, order)
	if affected != 1 || err != nil {
		log.Log.Errorf("shopJob_trace: update order status fail, orderId:%s, err:%s", orderId, err)
		session.Rollback()
		return errors.New("update order status fail")
	}
	
	list, err := shopModel.GetOrderProductList(order.OrderId)
	if err != nil {
		log.Log.Errorf("shopJob_trace: get order product list fail, orderId:%s, err:%s", order.OrderId, err)
		session.Rollback()
		return err
	}
	
	for _, item := range list {
		// 归还冻结库存
		if _, err := shopModel.UpdateProductSkuStock(fmt.Sprint(item.SkuId), item.Count * -1) ; err != nil {
			log.Log.Errorf("shop_trace: update product sku stock info fail, skuId:%d, err:%s", item.SkuId, err)
			session.Rollback()
			return err
		}
	}
	
	
	if _, err := DelShopOrderId(orderId); err != nil {
		log.Log.Errorf("shopJob_trace: del orderId fail, err:%s", err)
		session.Rollback()
		return err
	}
	
	session.Commit()
	
	return nil
}
