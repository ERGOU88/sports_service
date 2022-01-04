package corder

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	alipayCli "github.com/go-pay/gopay/alipay"
	wxCli "github.com/go-pay/gopay/wechat"
	"github.com/go-xorm/xorm"
	"net/url"
	"sports_service/server/app/config"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/global/rdskey"
	"sports_service/server/models"
	"sports_service/server/models/mappointment"
	"sports_service/server/models/mcoach"
	"sports_service/server/models/morder"
	"sports_service/server/models/mpay"
	"sports_service/server/models/mshop"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
	"sports_service/server/models/sms"
	"sports_service/server/tools/alipay"
	"sports_service/server/tools/tencentCloud"
	"sports_service/server/tools/wechat"
	"sports_service/server/util"
	"strconv"
	"strings"
	"time"
)

type OrderModule struct {
	context     *gin.Context
	engine      *xorm.Session
	order       *morder.OrderModel
	appointment *mappointment.AppointmentModel
	user        *muser.UserModel
	venue       *mvenue.VenueModel
	coach       *mcoach.CoachModel
	pay         *mpay.PayModel
	shop        *mshop.ShopModel
	sms         *sms.SmsModel
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
		pay: mpay.NewPayModel(socket),
		shop: mshop.NewShop(appSocket),
		sms: sms.NewSmsModel(),
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


// 支付宝通知 包含[支付成功、部分退款成功、全额退款成功]
func (svc *OrderModule) AliPayNotify(params url.Values, body string) int {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("payNotify_trace: session begin fail, err:%s", err)
		return errdef.ERROR
	}
	
	orderId := params.Get("out_trade_no")
	length := len(orderId)
	switch length {
	// 场馆订单
	case 16:
		return svc.AliPayNotifyByVenue(params, body)
	// 商城订单
	case 18:
		return svc.AliPayNotifyByShop(params, body)
	}
	
	log.Log.Errorf("order_trace: invalid orderId, orderId:%s", orderId)
	return errdef.ERROR
}

// 支付宝通知 商城订单 [暂只有支付成功]
func (svc *OrderModule) AliPayNotifyByShop(params url.Values, body string) int {
	orderId := params.Get("out_trade_no")
	order, err := svc.shop.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Errorf("aliNotify_trace: get order fail, orderId:%s, err:%s", orderId, err)
		return errdef.ERROR
	}
	
	status := strings.Trim(params.Get("trade_status"), " ")
	payTime, _ := time.ParseInLocation("2006-01-02 15:04:05", params.Get("gmt_payment"), time.Local)
	tradeNo := params.Get("trade_no")
	
	switch status {
	// 成功需区分 部分退款和支付成功
	case consts.TradeSuccess:
		switch order.PayStatus {
		// 待支付状态
		case consts.SHOP_ORDER_TYPE_WAIT:
			amount, err := strconv.ParseFloat(strings.Trim(params.Get("total_amount"), " "), 64)
			if err != nil {
				log.Log.Errorf("aliNotify_trace: parse float fail, err:%s", err)
				return errdef.ERROR
			}
			
			if int(amount * 100) != order.PayAmount {
				log.Log.Error("aliNotify_trace: amount not match, orderAmount:%d, amount:%d",
					order.PayAmount, amount * 100)
				return errdef.ERROR
			}
			
			
			condition := fmt.Sprintf("`pay_status`=%d AND `order_id`='%s'", consts.SHOP_ORDER_TYPE_WAIT, order.OrderId)
			cols := "pay_status, transaction, pay_time, update_at"
			order.PayStatus = consts.SHOP_ORDER_TYPE_PAID
			now := int(time.Now().Unix())
			order.UpdateAt = now
			order.Transaction = tradeNo
			order.PayTime = int(payTime.Unix())
			// 更新订单状态
			if _, err := svc.shop.UpdateOrderInfo(condition, cols, order); err != nil {
				log.Log.Errorf("shop_trace: update order info fail, orderId:%s, err:%v", order.OrderId, err)
				return errdef.SHOP_ORDER_UPDATE_FAIL
			}
			
		default:
			return errdef.ERROR
		}
	
	
	// 全额退款
	case consts.TradeClosed:
		log.Log.Debug("trade closed, order_id:%v", orderId)
	case consts.WaitBuyerPay:
		log.Log.Debug("wait buyer pay, order_id:%v", orderId)
	case consts.TradeFinished:
		log.Log.Debug("trade finished, order_id:%v", orderId)
	}
	
	
	return errdef.SUCCESS
}

// 支付宝通知 场馆订单包含[支付成功、部分退款成功、全额退款成功]
func (svc *OrderModule) AliPayNotifyByVenue(params url.Values, body string) int {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("aliNotify_trace: session begin fail, err:%s", err)
		return errdef.ERROR
	}

	orderId := params.Get("out_trade_no")
	order, err := svc.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Error("aliNotify_trace: order not found, orderId:%s, err:%s", orderId, err)
		svc.engine.Rollback()
		return errdef.ERROR
	}

	status := strings.Trim(params.Get("trade_status"), " ")
	payTime, _ := time.ParseInLocation("2006-01-02 15:04:05", params.Get("gmt_payment"), time.Local)
	tradeNo := params.Get("trade_no")

	switch status {
	// 成功需区分 部分退款和支付成功
	case consts.TradeSuccess:
		switch order.Status {
		// 待支付状态
		case consts.ORDER_TYPE_WAIT:
			amount, err := strconv.ParseFloat(strings.Trim(params.Get("total_amount"), " "), 64)
			if err != nil {
				log.Log.Errorf("aliNotify_trace: parse float fail, err:%s", err)
				svc.engine.Rollback()
				return errdef.ERROR
			}

			if int(amount * 100) != order.Amount {
				log.Log.Error("aliNotify_trace: amount not match, orderAmount:%d, amount:%d",
					order.Amount, amount * 100)
				svc.engine.Rollback()
				return errdef.ERROR
			}

			if err := svc.OrderProcess(orderId, body, tradeNo, payTime.Unix(), consts.PAY_NOTIFY, order.RefundAmount,
				order.RefundFee); err != nil {
				log.Log.Errorf("aliNotify_trace: order process fail, orderId:%s, err:%s", orderId, err)
				svc.engine.Rollback()
				return errdef.ERROR
			}

		// 退款中[部分退款]
		case consts.ORDER_TYPE_REFUND_WAIT:
			refundAmount, err := strconv.ParseFloat(strings.Trim(params.Get("refund_fee"), " "), 64)
			if err != nil {
				log.Log.Errorf("aliNotify_trace: parse float fail, err:%s", err)
				svc.engine.Rollback()
				return errdef.ERROR
			}

			if int(refundAmount* 100) != order.RefundAmount {
				log.Log.Errorf("aliNotify_trace: refundAmount not match, refundAmount:%d, order.RefundAmount:%d",
					refundAmount* 100, order.RefundAmount)
				svc.engine.Rollback()
				return errdef.ERROR
			}

			if err := svc.OrderProcess(order.PayOrderId, body, tradeNo, int64(order.PayTime), consts.REFUND_NOTIFY,
				order.RefundAmount, order.RefundFee); err != nil {
				log.Log.Errorf("aliNotify_trace: order process fail, orderId:%s, err:%s", orderId, err)
				svc.engine.Rollback()
				return errdef.ERROR
			}

			refundTradeNo := params.Get("out_biz_no")
			svc.order.RefundRecord.Status = 1
			svc.order.RefundRecord.UpdateAt = int(time.Now().Unix())
			refundTime, _ := time.ParseInLocation("2006-01-02 15:04:05", params.Get("gmt_refund"), time.Local)
			svc.order.RefundRecord.RefundTime = int(refundTime.Unix())
			affected, err := svc.order.UpdateRefundRecordStatus(refundTradeNo)
			if affected != 1 || err != nil {
				log.Log.Errorf("aliNotify_trace: update refund record fail, refundTradeNo:%s, err:%s",
					refundTradeNo, err)

				svc.engine.Rollback()
				return errdef.ERROR
			}


		default:
			log.Log.Errorf("invalid order status, orderId:%s, status:%d", orderId, order.Status)
			svc.engine.Rollback()
			return errdef.ERROR
		}


	// 全额退款
	case consts.TradeClosed:
		log.Log.Debug("trade closed, order_id:%v", orderId)
		refundAmount, err := strconv.ParseFloat(strings.Trim(params.Get("refund_fee"), " "), 64)
		if err != nil {
			log.Log.Errorf("aliNotify_trace: parse float fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.ERROR
		}

		if int(refundAmount* 100) != order.RefundAmount {
			log.Log.Errorf("aliNotify_trace: refundAmount not match, refundAmount:%d, order.RefundAmount:%d",
				refundAmount* 100, order.RefundAmount)
			svc.engine.Rollback()
			return errdef.ERROR
		}

		if err := svc.OrderProcess(order.PayOrderId, body, tradeNo, int64(order.PayTime), consts.REFUND_NOTIFY, order.RefundAmount, order.RefundFee); err != nil {
			log.Log.Errorf("aliNotify_trace: order process fail, orderId:%s, err:%s", orderId, err)
			svc.engine.Rollback()
			return errdef.ERROR
		}

		refundTradeNo := params.Get("out_biz_no")
		refundTime, _ := time.ParseInLocation("2006-01-02 15:04:05", params.Get("gmt_refund"), time.Local)
		// 状态为已退款
		svc.order.RefundRecord.Status = 1
		svc.order.RefundRecord.RefundTime = int(refundTime.Unix())
		svc.order.RefundRecord.UpdateAt = int(time.Now().Unix())
		affected, err := svc.order.UpdateRefundRecordStatus(refundTradeNo)
		if affected != 1 || err != nil {
			log.Log.Errorf("aliNotify_trace: update refund record fail, refundTradeNo:%s, err:%s",
				refundTradeNo, err)

			svc.engine.Rollback()
			return errdef.ERROR
		}

	case consts.WaitBuyerPay:
		log.Log.Debug("wait buyer pay, order_id:%v", orderId)
		svc.engine.Rollback()
		return errdef.SUCCESS

	case consts.TradeFinished:
		log.Log.Debug("trade finished, order_id:%v", orderId)
		svc.engine.Rollback()
		return errdef.SUCCESS
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

func (svc *OrderModule) WechatPayNotify(orderId, body, tradeNo, totalFee, refundNo string, payTm, refundTm int64, changeType int) error {
	length := len(orderId)
	switch length {
	// 场馆订单
	case 16:
		return svc.WechatPayNotifyByVenue(orderId, body, tradeNo, totalFee, refundNo, payTm, refundTm, changeType)
	// 商城订单
	case 18:
		return svc.WechatPayNotifyByShop(orderId, tradeNo, totalFee, payTm)
	}
	
	return nil
}

func (svc *OrderModule) WechatPayNotifyByShop(orderId, tradeNo, totalFee string, payTm int64) error {
	order, err := svc.shop.GetOrder(orderId)
	if order == nil || err != nil {
		log.Log.Errorf("wxNotify_trace: get order fail, orderId:%s, err:%s", orderId, err)
		return errors.New("order not found")
	}
	
	switch order.PayStatus {
	// 待支付状态
	case consts.SHOP_ORDER_TYPE_WAIT:
		amount, err := strconv.Atoi(totalFee)
		if err != nil {
			log.Log.Errorf("wxNotify_trace: parse float fail, err:%s", err)
			return err
		}
		
		if amount != order.PayAmount {
			log.Log.Error("wxNotify_trace: amount not match, orderAmount:%d, amount:%d",
				order.PayAmount, amount)
			return errors.New("amount not match")
		}
		
		
		condition := fmt.Sprintf("`pay_status`=%d AND `order_id`='%s'", consts.SHOP_ORDER_TYPE_WAIT, order.OrderId)
		cols := "pay_status, transaction, pay_time, update_at"
		order.PayStatus = consts.SHOP_ORDER_TYPE_PAID
		now := int(time.Now().Unix())
		order.UpdateAt = now
		order.Transaction = tradeNo
		order.PayTime = int(payTm)
		// 更新订单状态
		if _, err := svc.shop.UpdateOrderInfo(condition, cols, order); err != nil {
			log.Log.Errorf("shop_trace: update order info fail, orderId:%s, err:%v", order.OrderId, err)
			return err
		}
	
	default:
		return errors.New("invalid order type")
	}
	
	return nil
}

// 微信支付/退款 回调
func (svc *OrderModule) WechatPayNotifyByVenue(orderId, body, tradeNo, totalFee, refundNo string, payTm, refundTm int64, changeType int) error {
	if err := svc.engine.Begin(); err != nil {
		return err
	}

	ok, err := svc.order.GetOrder(orderId)
	if !ok || err != nil {
		log.Log.Errorf("wxNotify_trace: order not exists, orderId:%s", orderId)
		svc.engine.Rollback()
		return errors.New("order not exists")
	}
	
	fee, err := strconv.Atoi(totalFee)
	if err != nil {
		log.Log.Error("wxNotify_trace: amount fail, orderId:%s, err:%s", orderId, err)
		svc.engine.Rollback()
		return err
	}
	
	if fee != svc.order.Order.Amount {
		log.Log.Error("wxNotify_trace: amount not match, orderAmount:%d, amount:%d", svc.order.Order.Amount, fee)
		svc.engine.Rollback()
		return errors.New("amount not match")
	}

	if err := svc.OrderProcess(orderId, body, tradeNo, payTm, changeType, svc.order.Order.RefundAmount, svc.order.Order.RefundFee); err != nil {
		log.Log.Errorf("wxNotify_trace: order Process fail, orderId:%s, err:%s", orderId, err)
		svc.engine.Rollback()
		return err
	}

	if changeType == consts.REFUND_NOTIFY {
		// 状态为已退款
		svc.order.RefundRecord.Status = 1
		svc.order.RefundRecord.RefundTime = int(refundTm)
		svc.order.RefundRecord.UpdateAt = int(time.Now().Unix())
		affected, err := svc.order.UpdateRefundRecordStatus(refundNo)
		if affected != 1 || err != nil {
			log.Log.Errorf("wxNotify_trace: update refund record fail, refundTradeNo:%s, err:%s",
				refundNo, err)

			svc.engine.Rollback()
			return errors.New("update fail")
		}
	}

	svc.engine.Commit()

	return nil
}

// 订单处理流程 1 支付成功 2 退款流程 [用户申请退款 第三方回调成功时执行] 3 退款申请 4 取消订单
func (svc *OrderModule) OrderProcess(orderId, body, tradeNo string, payTm int64, changeType, refundAmount, refundFee int) error {
	// tips: 不可直接更新状态 并发情况下会有问题
	// 订单当前状态, 需更新的状态, 快照记录需更新的状态
	var curStatus, status, recordStatus int
	switch changeType {
	case consts.PAY_NOTIFY:
		// 如果是支付成功回调 则订单当前状态应是 待支付 需更新状态为 已支付
		curStatus = consts.ORDER_TYPE_WAIT
		status = consts.ORDER_TYPE_PAID
		recordStatus = 0
		svc.order.Order.IsCallback = 1
	case consts.REFUND_NOTIFY:
		curStatus = consts.ORDER_TYPE_REFUND_WAIT
		status = consts.ORDER_TYPE_REFUND_SUCCESS
		recordStatus = 1
		svc.order.Order.IsCallback = 1
	case consts.APPLY_REFUND:
		// 如果是申请退款 则订单当前状态 应是已付款 需更新状态为 退款中
		curStatus = consts.ORDER_TYPE_PAID
		status = consts.ORDER_TYPE_REFUND_WAIT
		recordStatus = 1
		// 如果退款金额和退款手续费均为0 则表示退款单金额为0 直接将退款状态置为成功
		if refundFee == 0 && refundAmount == 0 {
			status = consts.ORDER_TYPE_REFUND_SUCCESS
		}

	case consts.CANCEL_ORDER:
		// 如果是取消订单 则订单当前状态 应是待支付 需更新状态为 未支付
		curStatus = consts.ORDER_TYPE_WAIT
		status = consts.ORDER_TYPE_UNPAID
		recordStatus = 1
	}

	now := int(time.Now().Unix())
	svc.order.Order.Status = status
	svc.order.Order.PayTime = int(payTm)
	svc.order.Order.Transaction = tradeNo
	svc.order.Order.UpdateAt = now
	// 退款金额
	svc.order.Order.RefundAmount = refundAmount
	// 退款手续费
	svc.order.Order.RefundFee = refundFee
	// 更新订单状态
	affected, err := svc.order.UpdateOrderStatus(orderId, curStatus)
	if affected != 1 || err != nil {
		log.Log.Errorf("payNotify_trace: update order status fail, orderId:%s", orderId)
		return errors.New("update order status fail")
	}

	svc.order.OrderProduct.Status = status
	svc.order.OrderProduct.UpdateAt = now
	// 更新订单商品流水状态
	if _, err = svc.order.UpdateOrderProductStatus(orderId, curStatus); err != nil {
		log.Log.Errorf("payNotify_trace: update order product status fail, err:%s, orderId:%s", err, orderId)
		return errors.New("update order product status fail")
	}

	switch svc.order.Order.ProductType {
	case consts.ORDER_TYPE_APPOINTMENT_VENUE, consts.ORDER_TYPE_APPOINTMENT_COACH, consts.ORDER_TYPE_APPOINTMENT_COURSE:
		if err := svc.AppointmentOrderProcess(changeType, now, recordStatus, orderId); err != nil {
			return err
		}
		// 申请退款 / 取消订单 需归还库存 及 抵扣的会员时长
		//if changeType == consts.APPLY_REFUND || changeType == consts.CANCEL_ORDER {
		//	if err := svc.UpdateAppointmentInfo(orderId, now); err != nil {
		//		log.Log.Errorf("payNotify_trace: update appointment info fail, err:%s, orderId:%s", err, orderId)
		//		return err
		//	}
		//}


		//if svc.order.Order.ProductType == consts.ORDER_TYPE_APPOINTMENT_VENUE {
		//	// 更新标签状态[废弃]
		//	svc.appointment.Labels.Status = 1
		//	if _, err = svc.appointment.UpdateLabelsStatus(orderId, 0); err != nil {
		//		log.Log.Errorf("order_trace: update labels status fail, orderId:%s, err:%s", orderId, err)
		//		return errors.New("update label status fail")
		//	}
		//}

	case consts.ORDER_TYPE_MONTH_CARD, consts.ORDER_TYPE_SEANSON_CARD, consts.ORDER_TYPE_HALF_YEAR_CARD, consts.ORDER_TYPE_YEAR_CARD:
		//ok, err := svc.order.GetCardRecordByOrderId(orderId)
		//if !ok || err != nil {
		//	log.Log.Errorf("payNotify_trace: get card record by id fail, orderId:%s, err:%s", orderId, err)
		//	return errors.New("get card record fail")
		//}

		// 支付成功 需更新会员数据 [会员卡不可退款]
		//if changeType == consts.PAY_NOTIFY {
		//	// 更新会员可用时长 及 过期时长
		//	if err := svc.UpdateVipInfo(svc.order.Order.UserId, svc.order.CardRecord.VenueId, svc.order.CardRecord.ProductType,
		//		now, svc.order.CardRecord.PurchasedNum, svc.order.CardRecord.ExpireDuration, svc.order.CardRecord.Duration,
		//		changeType); err != nil {
		//		log.Log.Errorf("payNotify_trace: update vip info fail, orderId:%s, err:%s", orderId, err)
		//		return err
		//	}
		//}

		if err := svc.CardOrderProcess(changeType, now, orderId); err != nil {
			return err
		}

	// 混合型订单
	case consts.ORDER_TYPE_MIXED:
		products, err := svc.order.GetOrderProductsById(svc.order.Order.PayOrderId)
		if len(products) == 0 || err != nil {
			log.Log.Errorf("payNotify_trace: get order products fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			return errors.New("get order products fail")
		}

		var (
			isProcess, hasProcess bool
		)
		for _, item := range products {
			switch item.ProductType {
			// 预约类 包含 场馆、课程预约
			case consts.ORDER_TYPE_APPOINTMENT_VENUE, consts.ORDER_TYPE_APPOINTMENT_COACH, consts.ORDER_TYPE_APPOINTMENT_COURSE:
				if isProcess {
					continue
				}
				if err := svc.AppointmentOrderProcess(changeType, now, recordStatus, orderId); err != nil {
					return err
				}

				isProcess = true
			// 会员卡类
			case consts.ORDER_TYPE_MONTH_CARD, consts.ORDER_TYPE_SEANSON_CARD, consts.ORDER_TYPE_HALF_YEAR_CARD, consts.ORDER_TYPE_YEAR_CARD:
				if hasProcess {
					continue
				}

				if err := svc.CardOrderProcess(changeType, now, orderId); err != nil {
					return err
				}

				hasProcess = true
			}
		}
	}

	if changeType == consts.PAY_NOTIFY || changeType == consts.REFUND_NOTIFY {
		// 记录回调信息
		if err := svc.RecordNotifyInfo(now, changeType, orderId, body, tradeNo); err != nil {
			log.Log.Errorf("payNotify_trace: record notify info fail, orderId:%s, err:%s", orderId, err)
			return err
		}
	}

	log.Log.Debug("payNotify_trace: 订单成功， changeType:%d, orderId: %s", changeType, orderId)
	return nil
}

// 预约类 订单处理 tips: 直接load所属当前订单的所有预约快照 所以在循环订单商品列表时 执行一次即可
func (svc *OrderModule) AppointmentOrderProcess(changeType, now, recordStatus int, orderId string) error {
	// 更新订单对应的预约流水状态
	//if err := svc.appointment.UpdateAppointmentRecordStatus(orderId, now, status, curStatus); err != nil {
	//	log.Log.Errorf("payNotify_trace: update order product status fail, err:%s, orderId:%s", err, orderId)
	//	svc.engine.Rollback()
	//	return err
	//}

	// 申请退款 / 取消订单 需归还库存 及 抵扣的会员时长
	if changeType == consts.APPLY_REFUND || changeType == consts.CANCEL_ORDER {
		if err := svc.UpdateAppointmentInfo(orderId, now); err != nil {
			log.Log.Errorf("payNotify_trace: update appointment info fail, err:%s, orderId:%s", err, orderId)
			return err
		}

		if svc.order.Order.ProductType == consts.ORDER_TYPE_APPOINTMENT_VENUE {
			// 更新标签状态[废弃]
			svc.appointment.Labels.Status = 1
			if _, err := svc.appointment.UpdateLabelsStatus(orderId, 0); err != nil {
				log.Log.Errorf("order_trace: update labels status fail, orderId:%s, err:%s", orderId, err)
				return errors.New("update label status fail")
			}
		}
	}

	// 申请退款 [不可用] / 支付成功 [可用] 修改预约流水状态
	if changeType == consts.APPLY_REFUND || changeType == consts.PAY_NOTIFY {
		// 更新订单对应的预约流水状态
		if err := svc.appointment.UpdateAppointmentRecordStatus(orderId, now, recordStatus); err != nil {
			log.Log.Errorf("order_trace: update order product status fail, err:%s, orderId:%s", err, orderId)
			return err
		}

	}

	return nil
}

// 会员卡类订单处理 tips: 直接load所属当前订单的所有会员卡类快照 所以在循环订单商品列表时 执行一次即可
func (svc *OrderModule) CardOrderProcess(changeType, now int, orderId string) error {
	list, err := svc.order.FindCardRecordsByOrderId(orderId)
	if len(list) == 0 || err != nil {
		log.Log.Errorf("payNotify_trace: get card record by id fail, orderId:%s, err:%s", orderId, err)
		return errors.New("get card record fail")
	}

	// 支付成功 需更新会员数据 [会员卡不可退款]
	if changeType == consts.PAY_NOTIFY {
		for _, card := range list {
			// 更新会员可用时长 及 过期时长
			if err := svc.UpdateVipInfo(svc.order.Order.UserId, card.VenueId, card.ProductType,
				now, card.PurchasedNum, card.ExpireDuration, card.Duration,
				changeType); err != nil {
				log.Log.Errorf("payNotify_trace: update vip info fail, orderId:%s, err:%s", orderId, err)
				return err
			}
			
			cols := "use_user_id"
			card.UseUserId = svc.order.Order.UserId
			if _, err := svc.order.UpdateCardRecordInfo(cols, card); err != nil {
				log.Log.Errorf("payNotify_trace: update record info fail, id:%d, err:%s", card.Id, err)
				return err
			}
		}
		
		
	}

	return nil
}

// 更新预约信息
func (svc *OrderModule) UpdateAppointmentInfo(orderId string, now int) error {
	// 获取订单对应的预约流水
	list, err := svc.appointment.GetAppointmentRecordByOrderId(orderId)
	if err != nil {
		log.Log.Errorf("order_trace: get appointment record by orderId fail, orderId:%s, err:%s", orderId, err)
		return err
	}

	for _, record := range list {
		switch record.AppointmentType {
		case consts.APPOINTMENT_VENUE:
			// 归还场馆预约对应节点的冻结库存
			affected, err := svc.appointment.RevertStockNum(record.TimeNode, record.Date, record.PurchasedNum * -1, now,
				record.AppointmentType, int(record.VenueId))
			if affected != 1 || err != nil {
				log.Log.Errorf("order_trace: update stock info fail, orderId:%s, err:%s, affected:%d, id:%d", orderId, err, affected, record.Id)
				return errors.New("update stock info fail")
			}

			// 归还抵扣的会员时长
			if record.DeductionTm > 0 {
				affected, err := svc.appointment.UpdateVenueVipInfo(int(record.DeductionTm), record.VenueId, record.UserId)
				if affected != 1 || err != nil {
					log.Log.Errorf("order_trace: revert vip duration fail, orderId:%s, err:%s", record.PayOrderId, err)
					return err
				}

			}

		case consts.APPOINTMENT_COACH,consts.APPOINTMENT_COURSE:
			// 归还课程对应节点的冻结库存
			affected, err := svc.appointment.RevertCourseStockNum(record.TimeNode, record.Date,  record.PurchasedNum * -1, now,
				record.AppointmentType, int(record.VenueId), int(record.CourseId), int(record.CoachId))
			if affected != 1 || err != nil {
				log.Log.Errorf("orderJob_trace: update stock info fail, orderId:%s, err:%s, affected:%d, id:%d", orderId, err, affected, record.Id)
				return errors.New("update stock info fail")
			}
		}

	}

	return nil
}


func (svc *OrderModule) RecordNotifyInfo(now, notifyType int, orderId, body, tradeNo string) error {
	svc.order.Notify.CreateAt = now
	svc.order.Notify.UpdateAt = now
	svc.order.Notify.NotifyType = notifyType
	svc.order.Notify.PayChannelId = svc.order.Order.PayChannelId
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
		condition = fmt.Sprintf("order_type=1001 AND is_delete=0 AND status >= 0 AND user_id=%s", userId)
	case "0":
		// 0 待支付
		condition = fmt.Sprintf("order_type=1001 AND is_delete=0 AND status = 0 AND user_id=%s", userId)
	case "1":
		// 1 可使用 (仅展示 1001 场馆预约 2201 次卡 3001 私教预约 3002 课程预约)
		condition = fmt.Sprintf("order_type=1001 AND is_delete=0  AND status = 2 AND user_id=%s AND product_type in(1001,2201,3001,3002)", userId)
	case "2":
		// 2 退款/售后 包含[3 已完成 4 退款中 5 已退款 6 已过期]
		condition = fmt.Sprintf("order_type=1001 AND is_delete=0 AND status >= 3 AND user_id=%s", userId)
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
	res := make([]*morder.OrderInfo, 0)
	for _, order := range list {
		info := new(morder.OrderInfo)
		info.OrderType = int32(order.ProductType)
		cstSh, _ := time.LoadLocation("Asia/Shanghai")
		info.CreatAt = time.Unix(int64(order.CreateAt), 0).In(cstSh).Format(consts.FORMAT_TM)
		info.OrderStatus = int32(order.Status)
		info.OrderId = order.PayOrderId
		info.UserId = order.UserId
		info.Amount = fmt.Sprintf("%.2f", float64(order.Amount)/100)
		info.TotalAmount = order.Amount
		info.Title = order.Subject
		info.IsGift = order.IsGift
		info.GiftStatus = order.GiftStatus
		// 是否可退款
		can, _, _, _, err := svc.CanRefund(order.Amount, order.Status, order.ProductType, order.PayTime,
			order.PayOrderId, order.Extra)
		if err != nil {
			log.Log.Errorf("order_trace: refund fail, orderId:%s, can:%v, err:%s", svc.order.Order.PayOrderId, can, err)
		}

		info.CanRefund = can
		extra := &mappointment.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(order.Extra, extra); err != nil {
			log.Log.Errorf("order_trace: unmarshal extra fail, err:%s, orderId:%s", err, order.PayOrderId)
			continue
		}
		
		// 如果是赠品 判断是否可赠送
		if info.IsGift == 1 && info.GiftStatus == 0 {
			records, err := svc.appointment.GetAppointmentRecordByOrderId(order.PayOrderId)
			if err != nil || len(records) == 0 {
				log.Log.Errorf("order_trace: get appointment record by orderId fail, err:%s", err)
				continue
			}
			
			for _, item := range records {
				// 结束时间 > 当前时间 则礼物已过期
				if item.EndTm > int(time.Now().Unix()) {
					order.GiftStatus = 1
				}
			}
		}
		
		info.GiftStatus = order.GiftStatus
		// 赠品已过期 更新 订单赠品状态
		if info.GiftStatus == 1 {
			svc.order.Order = order
			cols := "gift_status"
			// 更新赠品状态
			if _, err := svc.order.UpdateOrderInfo(cols); err != nil {
				log.Log.Errorf("order_trace: update order info fail, orderId:%s, err:%s", order.PayOrderId, err)
				continue
			}
		}
		
		info.Count = extra.Count
		info.ProductImg = tencentCloud.BucketURI(extra.ProductImg)

		switch info.OrderType {
		// 预约场馆、私教、大课
		case consts.ORDER_TYPE_APPOINTMENT_VENUE:


		case consts.ORDER_TYPE_APPOINTMENT_COURSE:
			// 课程名称 + 老师名称
			info.Title = fmt.Sprintf("%s %s", extra.CourseName, extra.CoachName)
			if len(extra.TimeNodeInfo) > 0 {
				info.TimeNode = fmt.Sprintf("%s %s", extra.TimeNodeInfo[0].Date, extra.TimeNodeInfo[0].TimeNode)
			}

		case consts.ORDER_TYPE_APPOINTMENT_COACH:
			// 私教名称 + 课程名称
			info.Title = fmt.Sprintf("%s %s", extra.CoachName, extra.CourseName)
			if order.Status == consts.ORDER_TYPE_COMPLETED {
				// 查询是否评价
				ok, err := svc.coach.HasEvaluateByUserId(order.UserId, order.PayOrderId)
				if ok || err != nil {
					log.Log.Errorf("order_trace: already evaluate, userId:%s, orderId:%s", order.UserId, order.PayOrderId)
				}

				info.HasEvaluate = ok
			}

			if len(extra.TimeNodeInfo) > 0 {
				info.TimeNode = fmt.Sprintf("%s %s", extra.TimeNodeInfo[0].Date, extra.TimeNodeInfo[0].TimeNode)
			}

		case consts.ORDER_TYPE_MONTH_CARD, consts.ORDER_TYPE_SEANSON_CARD, consts.ORDER_TYPE_YEAR_CARD:
			products, err := svc.order.GetOrderProductsById(order.PayOrderId)
			if len(products) == 0 || err != nil {
				continue
			}

		}

		res = append(res, info)
	}

	return res
}

// 更新会员信息
func (svc *OrderModule) UpdateVipInfo(userId string, venueId int64, productType, now, count, expireDuration, duration, notifyType int) error {
	ok, err := svc.venue.GetVenueVipInfo(userId, venueId)
	if err != nil {
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
	//	svc.venue.Vip.Condition += int64(duration * -1)
	//	svc.venue.Vip.EndTm = svc.venue.Vip.EndTm - int64(expireDuration * count)
	//	cols = "end_tm, duration, update_at"
	// 支付成功回调通知
	case consts.PAY_NOTIFY:
		// 如果vip结束时间 >= 当前时间戳 则为续费
		if int(svc.venue.Vip.EndTm) >= now {
			svc.venue.Vip.Duration += int64(duration)
			svc.venue.Vip.EndTm = svc.venue.Vip.EndTm + int64(expireDuration * count)
			// 当前购买的会员level更高
			if svc.venue.Vip.VipType < productType {
				svc.venue.Vip.VipType = productType
			}
			cols = "vip_type, end_tm, duration, update_at"
		} else {
			// 否则 为 重新购买
			svc.venue.Vip.StartTm = int64(now)
			// 过期时间 叠加
			svc.venue.Vip.EndTm = int64(now + expireDuration * count)
			// 可用时长
			svc.venue.Vip.Duration = int64(duration)
			// 会员类型
			svc.venue.Vip.VipType = productType
			cols = "vip_type, start_tm, end_tm, duration, update_at"
		}

	default:
		log.Log.Errorf("order_trace: unsupported notify type:%d", notifyType)
		return errors.New("unsupported notify type")
	}

	// 不存在 新增
	if !ok {
		svc.venue.Vip.VenueId = venueId
		svc.venue.Vip.CreateAt = now
		svc.venue.Vip.UserId = userId
		if _, err := svc.venue.AddVenueVipInfo(); err != nil {
			log.Log.Errorf("venue_trace: add vip info fail, userId:%s, err:%s", userId, err)
			return errors.New("add vip info fail")
		}
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

	if svc.order.Order.Status == consts.ORDER_TYPE_COMPLETED && svc.order.Order.ProductType == consts.ORDER_TYPE_APPOINTMENT_COACH  {
		// 查询是否评价
		ok, err := svc.coach.HasEvaluateByUserId(svc.order.Order.UserId, svc.order.Order.PayOrderId)
		if !ok || err != nil {
			log.Log.Errorf("order_trace: already evaluate, userId:%s, orderId:%s", svc.order.Order.UserId, svc.order.Order.PayOrderId)
		}

		rsp.HasEvaluate = ok
	}

	rsp.PayDuration = 0
	// 待支付订单 剩余支付时长
	if svc.order.Order.Status == consts.ORDER_TYPE_WAIT {
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
	rsp.RefundFee = svc.order.Order.RefundFee
	rsp.RefundAmount = svc.order.Order.RefundAmount
	rsp.IsGift = svc.order.Order.IsGift
	rsp.GiftStatus = svc.order.Order.GiftStatus

	// 是否可退款
	can, _, _, _, err := svc.CanRefund(svc.order.Order.Amount, svc.order.Order.Status, svc.order.Order.ProductType, svc.order.Order.PayTime,
		svc.order.Order.PayOrderId, svc.order.Order.Extra)
	if err != nil {
		log.Log.Errorf("order_trace: refund fail, orderId:%s, can:%v, err:%s", svc.order.Order.PayOrderId, can, err)
	}

	rsp.CanRefund = can

	return errdef.SUCCESS, rsp
}

// 订单退款流程
func (svc *OrderModule) OrderRefund(param *morder.ChangeOrder, executeType int) (int, int, int, int) {
	if err := svc.engine.Begin(); err != nil {
		log.Log.Errorf("order_trace: session begin fail, err:%s", err)
		return errdef.ERROR, 0, 0, 0
	}

	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("order_trace: user not exists, userId:%s", param.UserId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS, 0, 0, 0
	}

	ok, err := svc.order.GetOrder(param.OrderId)
	if !ok || err != nil {
		log.Log.Errorf("order_trace: order not exists, orderId:%s", param.OrderId)
		svc.engine.Rollback()
		return errdef.ORDER_NOT_EXISTS, 0, 0, 0
	}

	if svc.order.Order.UserId != user.UserId {
		log.Log.Errorf("order_trace: user not match, userId:%s, curUser:%s", svc.order.Order.UserId, user.UserId)
		svc.engine.Rollback()
		return errdef.ORDER_USER_NOT_MATCH, 0, 0, 0
	}

	// 是否可退款
	can, refundFee, ruleId, _, err := svc.CanRefund(svc.order.Order.Amount, svc.order.Order.Status, svc.order.Order.ProductType, svc.order.Order.PayTime,
		svc.order.Order.PayOrderId, svc.order.Order.Extra)
	if err != nil {
		log.Log.Errorf("order_trace: refund fail, orderId:%s, can:%v, err:%s", svc.order.Order.PayOrderId, can, err)
		svc.engine.Rollback()
		return errdef.ORDER_REFUND_FAIL, 0, 0, 0
	}

	if !can {
		log.Log.Errorf("order_trace: user can't refund, orderId:%s, can:%v, err:%s", svc.order.Order.PayOrderId, can, err)
		svc.engine.Rollback()
		return errdef.ORDER_NOT_ALLOW_REFUND, 0, 0, 0
	}

	// 可退款金额 = 订单实付金额 - 手续费
	refundAmount := svc.order.Order.Amount - refundFee
	// 如果是查询退款金额 及 手续费
	if executeType == consts.EXECUTE_TYPE_QUERY {
		svc.engine.Rollback()
		return errdef.SUCCESS, refundAmount, refundFee, ruleId
	}

	if err := svc.OrderProcess(svc.order.Order.PayOrderId, "", svc.order.Order.Transaction,
		int64(svc.order.Order.PayTime), consts.APPLY_REFUND, refundAmount, refundFee); err != nil {
		svc.engine.Rollback()
		return errdef.ORDER_PROCESS_FAIL, 0, 0, 0
	}

	// 可退款金额 > 0
	if refundAmount > 0 {
		outRequestNo := util.NewOrderId()
		svc.order.RefundRecord.RefundChannelId = svc.order.Order.PayChannelId
		svc.order.RefundRecord.RefundType = 1001
		svc.order.RefundRecord.RefundAmount = refundAmount
		svc.order.RefundRecord.RefundFee = refundFee
		svc.order.RefundRecord.CreateAt = int(time.Now().Unix())
		svc.order.RefundRecord.RefundTradeNo = outRequestNo
		svc.order.RefundRecord.UserId = svc.order.Order.UserId
		svc.order.RefundRecord.PayOrderId = svc.order.Order.PayOrderId
		svc.order.RefundRecord.Remark = fmt.Sprintf("%s%s", svc.order.Order.Subject, "退款")
		affected, err := svc.order.AddRefundRecord()
		if affected != 1 || err != nil {
			log.Log.Errorf("order_trace: add refund record fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			svc.engine.Rollback()
			return errdef.ORDER_ADD_REFUND_RECORD_FAIL, 0, 0, 0
		}

		// 第三方退款
		if _, err := svc.TradeRefund(refundAmount, outRequestNo); err != nil {
			log.Log.Errorf("order_trace: trade refund err:%s", err)
			svc.engine.Rollback()
			return errdef.ORDER_REFUND_FAIL, 0, 0, 0
		}
	}


	svc.engine.Commit()

	return errdef.SUCCESS, refundAmount, refundFee, ruleId
}

// 删除订单
func (svc *OrderModule) DeleteOrder(param *morder.ChangeOrder) int {
	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("order_trace: user not exists, userId:%s", param.UserId)
		return errdef.USER_NOT_EXISTS
	}

	ok, err := svc.order.GetOrder(param.OrderId)
	if !ok || err != nil {
		log.Log.Errorf("order_trace: order not exists, orderId:%s", param.OrderId)
		return errdef.ORDER_NOT_EXISTS
	}

	if svc.order.Order.UserId != user.UserId {
		log.Log.Errorf("order_trace: user not match, userId:%s, curUser:%s", svc.order.Order.UserId, user.UserId)
		return errdef.ORDER_DELETE_FAIL
	}

	productType := fmt.Sprint(svc.order.Order.ProductType)
	if len(productType) > 2 {
		productType = productType[0:2]
	}

	// 退款中 / 已支付的订单[不包含 购买会员卡类订单] 会员卡类 类型区间为2300-2399 / 已过期的订单 不可删除
	if svc.order.Order.Status == consts.ORDER_TYPE_PAID && productType != "23" ||
		svc.order.Order.Status == consts.ORDER_TYPE_REFUND_WAIT ||
		svc.order.Order.Status == consts.ORDER_TYPE_EXPIRE {
		log.Log.Errorf("order_trace: order can't delete, orderId:%s", svc.order.Order.PayOrderId)
		return errdef.ORDER_DELETE_FAIL
	}

	svc.order.Order.IsDelete = 1
	svc.order.Order.UpdateAt = int(time.Now().Unix())
	cols := "is_delete, update_at"
	affected, err := svc.order.UpdateOrderInfo(cols)
	if affected != 1 || err != nil {
		log.Log.Errorf("order_trace: update order info fail, err:%s,affected:%d", err, affected)
		return errdef.ORDER_DELETE_FAIL
	}

	return errdef.SUCCESS
}

// 校验签名
func (svc *OrderModule) VerifySign(payChannel, body, sign string, bm gopay.BodyMap) bool {
	//ok, err := svc.GetPaymentChannel(payChannel)
	//if !ok || err != nil {
	//	log.Log.Errorf("notify_trace: get payment channel fail, ok:%v, err:%s", ok, err)
	//	return false
	//}

	switch payChannel {
	case consts.ALIPAY:
		cli := alipay.NewAliPay(true, svc.pay.PayChannel.AppId, svc.pay.PayChannel.PrivateKey)
		ok, err := cli.VerifyData(body, "RSA2", sign)
		if !ok || err != nil {
			log.Log.Errorf("notify_trace: verify ali data fail, err:%s", err)
			return false
		}

	case consts.WEICHAT:
		cli := wechat.NewWechatPay(true, svc.pay.PayChannel.AppId, svc.pay.PayChannel.AppKey, svc.pay.PayChannel.AppSecret)
		ok, err := cli.VerifySign(bm)
		if !ok || err != nil {
			log.Log.Error("notify_trace: wx sign not match, err:%s", err)
			return false
		}

	default:
		log.Log.Errorf("notify_trace: unsupported payChannel:%d", payChannel)
		return false
	}

	return true
}

// 获取支付渠道配置
func (svc *OrderModule) GetPaymentChannel(payChannelId int) (bool, error) {
	return svc.pay.GetPaymentChannel(payChannelId)
}

// 交易退款 todo:计算手续费
func (svc *OrderModule) TradeRefund(refundAmount int, outRequestNo string) (string, error) {
	ok, err := svc.GetPaymentChannel(svc.order.Order.PayChannelId)
	if !ok || err != nil {
		log.Log.Errorf("order_trace: get payment channel fail, orderId:%s, ok:%v, err:%s", svc.order.Order.PayOrderId,
			ok, err)
		return "", errors.New("channel not found")
	}

	var body string
	switch svc.pay.PayChannel.Identifier {
	case consts.ALIPAY:
		// 支付宝
		resp, err := svc.AliRefund(refundAmount, svc.pay.PayChannel.AppId, svc.pay.PayChannel.PrivateKey, outRequestNo)
		if err != nil {
			log.Log.Errorf("order_trace: alipay refund fail, orderId:%s, payChannelId:%d", svc.order.Order.PayOrderId,
				svc.order.Order.PayChannelId)
			return "", err
		}

		log.Log.Infof("order_trace: orderId:%s, alipay refund resp:%+v", svc.order.Order.PayOrderId, resp)
		body, _ = util.JsonFast.MarshalToString(resp)

	case consts.WEICHAT:
		// 微信
		resp, err := svc.WechatRefund(refundAmount, svc.pay.PayChannel.AppId, svc.pay.PayChannel.AppKey,
			svc.pay.PayChannel.AppSecret, outRequestNo)
		if err != nil {
			log.Log.Errorf("order_trace: get wechatPay param fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			return "", err
		}

		log.Log.Infof("order_trace: orderId:%s, wechat refund resp:%+v", svc.order.Order.PayOrderId, resp)
		body, _ = util.JsonFast.MarshalToString(resp)

	default:
		log.Log.Errorf("order_trace: unsupported pay channel id:%d", svc.order.Order.PayChannelId)
		return "", errors.New("unsupported payType")
	}


	return body, nil
}

// 支付宝退款
func (svc *OrderModule) AliRefund(refundAmount int, appId, privateKey, outRequestNo string) (*alipayCli.TradeRefundResponse, error) {
	client := alipay.NewAliPay(true, appId, privateKey)
	client.OutRequestNo = outRequestNo
	client.OutTradeNo = svc.order.Order.PayOrderId
	client.RefundAmount = fmt.Sprintf("%.2f", float64(refundAmount)/100)
	client.RefundReason = fmt.Sprintf("%s%s", svc.order.Order.Subject, "退款")
	resp, err := client.TradeRefund()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 微信退款
func (svc *OrderModule) WechatRefund(refundAmount int, appId, merchantId, secret, outRequestNo string) (*wxCli.RefundResponse, error) {
	client := wechat.NewWechatPay(true, appId, merchantId, secret)
	client.OutTradeNo = svc.order.Order.PayOrderId
	client.TotalAmount = svc.order.Order.Amount
	client.RefundAmount = refundAmount
	client.RefundNotify = config.Global.WechatRefundNotify
	client.OutRefundNo = outRequestNo
	resp, err := client.TradeRefund()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 是否可退款
// 返回值 是否可退款 [true表示可退] 及 退款手续费
func (svc *OrderModule) CanRefund(amount, status, orderType, payTime int, orderId, extra string) (bool, int, int, int64, error) {
	// 如果订单状态不等于已支付
	if status != consts.ORDER_TYPE_PAID {
		return false, 0, 0, 0, nil
	}

	// 只有预约场馆/次卡/私教/课程可申请退款 同时需要判断 订单是否过期
	if orderType == consts.ORDER_TYPE_YEAR_CARD || orderType == consts.ORDER_TYPE_MONTH_CARD ||
		orderType == consts.ORDER_TYPE_SEANSON_CARD || orderType == consts.ORDER_TYPE_HALF_YEAR_CARD {
		return false, 0, 0, 0, errors.New("invalid order type")
	}

	var (
		refundFee, ruleId int
		// 预约类型订单 最早节点 开始时间 及 最后节点 结束时间
		startTime, endTime int64
	)
	now := time.Now().Unix()
	switch orderType {
	case consts.ORDER_TYPE_APPOINTMENT_VENUE, consts.ORDER_TYPE_APPOINTMENT_COACH, consts.ORDER_TYPE_APPOINTMENT_COURSE:
		// 获取预约流水[可用状态]
		infos, err := svc.appointment.GetAppointmentRecordByOrderId(orderId)
		if len(infos) == 0 || err != nil {
			log.Log.Errorf("order_trace: get appointment record by orderId fail, orderId:%s, err:%s", orderId, err)
			return false, 0, 0, 0, errors.New("get record fail")
		}

		rsp := &mappointment.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(extra, rsp); err != nil {
			log.Log.Errorf("order_trace: unmarshal extra fail, orderId:%s, err:%s", orderId, err)
			return false, 0, 0, 0, errors.New("unmarshal fail")
		}

		if len(rsp.TimeNodeInfo) == 0 {
			log.Log.Errorf("order_trace: time node is empty, orderId:%s", orderId)
			return false, 0, 0, 0, errors.New("time node empty")
		}

		// 默认预约可退
		can := true
		for _, node := range rsp.TimeNodeInfo {
			// todo: 可能会有开场前15分钟无法退款
			// 如果预约中 某个节点的开始时间 <= 当前时间 表示不能退款
			if node.StartTm - 15 * 60 <= now {
				can = false
				//return false, 0, nil
			}

			if startTime == 0 || startTime > node.StartTm {
				// 获取多个节点中 最早的开始时间
				startTime = node.StartTm
			}

			if endTime == 0 || endTime < node.EndTm {
				// 获取多个节点中 最后的结束时间
				endTime = node.EndTm
			}
		}

		// 不能退款
		if !can {
			return false, 0, 0, endTime, nil
		}

		// 可退款 则计算手续费
		refundFee, ruleId, err = svc.CalculationRefundFee(amount, int(startTime))
		if err != nil {
			log.Log.Errorf("order_trace: calculation refund fee fail, orderId:%s, err:%s", orderId, err)
			return false, 0, 0, endTime, nil
		}

	case consts.ORDER_TYPE_EXPERIENCE_CARD:
		// todo: 只处理app端次卡退款 一次购买 对应一条快照
		ok, err := svc.order.GetCardRecordByOrderId(orderId)
		if !ok || err != nil {
			log.Log.Errorf("order_trace: get vip card by id fail, orderId:%s, err:%s", orderId, err)
			return false, 0, 0, 0, errors.New("get vip card fail")
		}

		endTime = int64(payTime + svc.order.CardRecord.ExpireDuration)
		// 次卡  订单完成时间 + 过期时长 <= 当前时间戳 表示不能退款
		if endTime <= now {
			return false, 0, 0, endTime, nil
		}

		// 未过期 可退款, 次卡 全额退 则 手续费为0
		refundFee = 0

	default:
		log.Log.Errorf("order_trace: invalid order type, orderId:%s", orderId)
		return false, 0, 0, 0, errors.New("invalid order type")
	}


	return true, refundFee, ruleId, endTime, nil
}

// 获取券码信息
func (svc *OrderModule) GetCouponCodeInfo(userId, orderId string) (int, *morder.CouponCodeInfo) {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("order_trace: user not exists, userId:%s", userId)
		return errdef.USER_NOT_EXISTS, nil
	}

	ok, err := svc.order.GetOrder(orderId)
	if !ok || err != nil {
		log.Log.Errorf("order_trace: order not exists, orderId:%s", orderId)
		return errdef.ORDER_COUPON_CODE_FAIL, nil
	}

	if svc.order.Order.UserId != user.UserId {
		log.Log.Errorf("order_trace: user not match, userId:%s, curUser:%s", svc.order.Order.UserId, user.UserId)
		return errdef.ORDER_COUPON_CODE_FAIL, nil
	}

	// 只有支付成功之后才可查看券码
	if svc.order.Order.Status < consts.ORDER_TYPE_PAID {
		log.Log.Errorf("order_trace: invalid order status, status:%d", svc.order.Order.Status)
		return errdef.ORDER_COUPON_CODE_FAIL, nil
	}

	extra := &mappointment.OrderResp{}
	if err := util.JsonFast.UnmarshalFromString(svc.order.Order.Extra, extra); err != nil {
		log.Log.Errorf("order_trace: unmarshal extra fail, orderId:%s, err:%s", orderId, err)
		return errdef.ORDER_COUPON_CODE_FAIL, nil
	}

	resp := &morder.CouponCodeInfo{}
	resp.VenueName = extra.VenueName
	resp.Subject = svc.order.Order.Subject
	//resp.Code = fmt.Sprintf("o%s", strings.ToLower(util.GenSecret(util.CHAR_MODE, 18)))
	resp.Code = fmt.Sprintf("o%s", util.GenQrcodeInfo())
	resp.Count = extra.Count
	resp.TotalAmount = svc.order.Order.Amount
	//resp.QrCodeInfo = fmt.Sprintf("O%s", util.GenSecret(util.MIX_MODE, 16))
	resp.QrCodeInfo = resp.Code
	expire := int64(rdskey.KEY_EXPIRE_MIN * 60)
	resp.QrCodeExpireDuration = expire - 30
	cstSh, _ := time.LoadLocation("Asia/Shanghai")

	// 只有次卡/预约场馆可以查看券码
	switch svc.order.Order.ProductType {
	case consts.ORDER_TYPE_APPOINTMENT_VENUE:
		// 到期时间 展示最近时间
		var minTm int64
		for _, val := range extra.TimeNodeInfo {
			if minTm < val.StartTm {
				minTm = val.StartTm
			}
		}

		resp.ExpireTm = time.Unix(minTm + int64(extra.ExpireDuration), 0).In(cstSh).Format(consts.FORMAT_DATE)

	case consts.ORDER_TYPE_EXPERIENCE_CARD:

		resp.ExpireTm = time.Unix(int64(svc.order.Order.PayTime + int(extra.ExpireDuration)), 0).In(cstSh).Format(consts.FORMAT_DATE)

	default:
		log.Log.Errorf("order_trace: unsupported product type, type:%d", svc.order.Order.ProductType)
		return errdef.ORDER_COUPON_CODE_FAIL, nil
	}

	if err = svc.order.SaveQrCodeInfo(resp.QrCodeInfo, orderId, expire); err != nil {
		log.Log.Errorf("order_trace: save qrcode info fail, err:%s", err)
		return errdef.ORDER_COUPON_CODE_FAIL, nil
	}

	return errdef.SUCCESS, resp
}

// 订单取消
func (svc *OrderModule) OrderCancel(param *morder.ChangeOrder) int {
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
		return errdef.ORDER_CANCEL_FAIL
	}

	// 只有待支付状态订单可以取消
	if svc.order.Order.Status != consts.ORDER_TYPE_WAIT {
		log.Log.Errorf("order_trace: order not allow cancel, orderId:%s, status:%d", svc.order.Order.PayOrderId, svc.order.Order.Status)
		svc.engine.Rollback()
		return errdef.ORDER_NOT_ALLOW_CANCEL
	}

	// 取消订单流程
	if err := svc.OrderProcess(svc.order.Order.PayOrderId, "", svc.order.Order.Transaction, 0, consts.CANCEL_ORDER,
		svc.order.Order.RefundAmount, svc.order.Order.RefundFee); err != nil {
		svc.engine.Rollback()
		log.Log.Errorf("order_trace: order process fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
		return errdef.ORDER_CANCEL_FAIL
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}

// 获取已支付的订单列表
func (svc *OrderModule) GetOrderListByPaid() ([]*models.VenuePayOrders, error){
	// 成功支付 未消费/未过期/未退款的订单
	condition := fmt.Sprintf("status = 2 AND is_delete=0")
	list, err := svc.order.GetOrderListByStatus(condition, 0, 100)
	if err != nil {
		log.Log.Errorf("order_trace: get order list by status fail, err:%s", err)
		return nil, err
	}

	return list, nil
}

// 检查订单是否过期
func (svc *OrderModule) CheckOrderExpire() error {
	list, err := svc.GetOrderListByPaid()
	if err != nil {
		log.Log.Errorf("order_trace: get job order list fail, err:%s", err)
		return err
	}

	if len(list) == 0 {
		log.Log.Error("order_trace: list empty")
		return nil
	}

	for _, order := range list {
		if err := svc.engine.Begin(); err != nil {
			log.Log.Errorf("order_trace: session begin err:%s, orderId:%s", err)
			return err
		}

		svc.order.Order = order
		// 预约场馆/次卡 查看订单是否可退款
		can, _, _, endTm, err := svc.CanRefund(order.Amount, order.Status, order.ProductType, order.PayTime, order.PayOrderId, order.Extra)
		// 可以退款 表示未过期 或 出现错误 不处理
		if can || err != nil {
			log.Log.Errorf("order_trace: canRefund err:%s", err)
			svc.engine.Rollback()
			continue
		}

		loc, _ := time.LoadLocation("Asia/Shanghai")
		now := time.Now().In(loc).Unix()
		log.Log.Errorf("order_trace: orderId:%s, endTm:%d, now:%d", order.PayOrderId, endTm, now)
		// 节点结束时间 > 当前时间  表示未过期
		if endTm > now {
			svc.engine.Rollback()
			continue
		}

		// 已过期 只需要更新订单状态
		switch order.ProductType {
		// 场馆/次卡 变更为已过期状态
		case consts.ORDER_TYPE_APPOINTMENT_VENUE, consts.ORDER_TYPE_EXPERIENCE_CARD:
			svc.order.Order.Status = consts.ORDER_TYPE_EXPIRE
		// 课程/私教 变更为已完成
		case consts.ORDER_TYPE_APPOINTMENT_COACH, consts.ORDER_TYPE_APPOINTMENT_COURSE:
			svc.order.Order.Status = consts.ORDER_TYPE_COMPLETED
		default:
			svc.engine.Rollback()
			continue
		}

		// 更新订单状态
		affected, err := svc.order.UpdateOrderStatus(order.PayOrderId, consts.ORDER_TYPE_PAID)
		if affected != 1 || err != nil {
			log.Log.Errorf("payNotify_trace: update order status fail, orderId:%s", order.PayOrderId)
			svc.engine.Rollback()
			return errors.New("update order status fail")
		}

		svc.order.OrderProduct.Status = svc.order.Order.Status
		svc.order.OrderProduct.UpdateAt = int(time.Now().Unix())
		// 更新订单商品流水状态
		if _, err = svc.order.UpdateOrderProductStatus(order.PayOrderId, consts.ORDER_TYPE_PAID); err != nil {
			log.Log.Errorf("order_trace: update order product status fail, err:%s, affected:%d, orderId:%s", err, affected, order.PayOrderId)
			svc.engine.Rollback()
			return errors.New("update order product status fail")
		}

		svc.engine.Commit()
	}

	return nil

}

// 订单退款规则
func (svc *OrderModule) OrderRefundRules() (int, []*models.VenueRefundRules) {
	rules, err := svc.order.GetRefundRules()
	if rules == nil || err != nil {
		log.Log.Errorf("order_trace: get refund rules fail, err:%s", err)
		return errdef.ORDER_REFUND_FAIL, []*models.VenueRefundRules{}
	}

	return errdef.SUCCESS, rules
}

// 计算退票手续费
// amount 订单实付金额
// lastStartTime 场次/预约 最近开始时间
// 返回手续费 [单位 分]
func (svc *OrderModule) CalculationRefundFee(amount, lastStartTime int) (int, int, error) {
	rules, err := svc.order.GetRefundRules()
	if len(rules) == 0 || err != nil {
		log.Log.Errorf("order_trace: get refund rules fail, err:%s", err)
		return 0, 0, errors.New("get refund rules fail")
	}

	if amount == 0 {
		return 0, 0, nil
	}

	now := int(time.Now().Unix())
	var refundFee, ruleId int
	for _, rule := range rules {
		// 规则校验最大时长 != 0 表示需要校验时长区间
		if rule.RuleMaxDuration > 0 {
			// 最近开始时间 - 当前时间戳 >= 最小时长 &&  最近开始时间 - 当前时间戳 < 最大时长
			if lastStartTime - now >= rule.RuleMinDuration && lastStartTime - now < rule.RuleMaxDuration && rule.FeeRate > 0 {
				// 按此时间区间 计算手续费
				// 手续费 = 订单实付金额[分] * 退款费率[数据库比例已乘以100] 单位[分]
				if refundFee = amount * rule.FeeRate / 1e4; refundFee < rule.MinimumCharge {
					// 最低手续费
					refundFee = rule.MinimumCharge
					ruleId = rule.Id
				}

				svc.order.RefundRecord.MinimumCharge = rule.MinimumCharge
				svc.order.RefundRecord.FeeRate = rule.FeeRate

				break
			}
		}

		// 如果 规则校验最大时长 = 0 表示只需校验最小时长
		if rule.RuleMaxDuration == 0 {
			if lastStartTime - now >= rule.RuleMinDuration && rule.FeeRate > 0 {
				if refundFee = amount * rule.FeeRate / 1e4; refundFee < rule.MinimumCharge {
					// 最低手续费
					refundFee = rule.MinimumCharge
					ruleId = rule.Id
				}

				svc.order.RefundRecord.MinimumCharge = rule.MinimumCharge
				svc.order.RefundRecord.FeeRate = rule.FeeRate
				break
			}

		}
	}

	log.Log.Infof("amount:%d,refundFee:%d", amount, refundFee)

	// 如果订单金额 < 手续费 [场景：搞活动 1元抢购 假设退款最低为1元 则可能出现]
	if amount <= refundFee {
		// 全款退
		refundFee = 0
		// 属于特殊情况 规则id为0
		ruleId = 0
	}

	return refundFee, ruleId, nil
}

// 领取赠品
func (svc *OrderModule) ReceiveGift(param *morder.ReceiveGiftReq) int {
	if param.UserId == "" {
		return errdef.INVALID_PARAMS
	}
	
	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("order_trace: user not found, userId:%s", param.UserId)
		return errdef.USER_NOT_EXISTS
	}
	
	ok, err := svc.order.GetOrder(param.OrderId)
	if !ok || err != nil {
		log.Log.Errorf("order_trace: order not found, orderId:%s, err:%s", param.OrderId, err)
		return errdef.ORDER_NOT_EXISTS
	}
	
	// 是否为赠品
	if svc.order.Order.IsGift != 1 {
		log.Log.Errorf("order_trace: not gift, orderId:%s", param.OrderId)
		return errdef.ORDER_GIFT_NOT_ALLOW_RECEIVE
	}
	
	// 赠品已过期
	if svc.order.Order.GiftStatus == 1 {
		log.Log.Errorf("order_trace: gitf has expired, orderId:%s", param.OrderId)
		return errdef.ORDER_GIFT_HAS_EXPIRED
	}
	
	// 赠品已被领取
	if svc.order.Order.GiftStatus == 2 {
		log.Log.Errorf("order_trace: gift has received, orderId:%s", param.OrderId)
		return errdef.ORDER_GIFT_HAS_RECEIVED
	}
	
	if svc.order.Order.UserId == user.UserId {
		log.Log.Errorf("order_trace: gift not allow receive, orderId:%s", param.OrderId)
		return errdef.ORDER_GIFT_NOT_ALLOW_RECEIVE
	}
	
	records, err := svc.appointment.GetAppointmentRecordByOrderId(param.OrderId)
	if records == nil || err != nil {
		log.Log.Errorf("order_trace: get appointment record fail, orderId:%s, err:%s", param.OrderId, err)
		return errdef.ORDER_APPOINTMENT_RECORD_FAIL
	}
	
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR
	}
	
	for _, item := range records {
		if item.UseUserId != "" {
			log.Log.Errorf("order_trace: gift has received, orderId:%s", param.OrderId)
			svc.engine.Rollback()
			return errdef.ORDER_GIFT_HAS_RECEIVED
		}
		
		// 结束时间 > 当前时间 不可领取 赠品已过期
		if item.EndTm > int(time.Now().Unix()) {
			log.Log.Errorf("order_trace: gift has expired, orderId:%s", param.OrderId)
			svc.engine.Rollback()
			return errdef.ORDER_GIFT_HAS_EXPIRED
		}
		
		cols := "use_user_id"
		item.UseUserId = param.UserId
		condition := fmt.Sprintf("id=%d AND use_user_id=''", item.Id)
		// 更新使用者
		if _, err := svc.appointment.UpdateAppointmentRecordInfo(condition, cols, item); err != nil {
			log.Log.Errorf("order_trace: update appointment record fail, orderId:%s, useUserId:%s, err:%s",
				item.PayOrderId, param.UserId, err)
			svc.engine.Rollback()
			return errdef.ORDER_GIFT_RECEIVE_FAIL
		}
	}
	
	// 状态更新为已赠送/已领取
	svc.order.Order.GiftStatus = 2
	svc.order.Order.UpdateAt = int(time.Now().Unix())
	cols := "gift_status, update_at"
	if _, err := svc.order.UpdateOrderInfo(cols); err != nil {
		log.Log.Errorf("order_trace: update gift status fail, orderId:%s, useUserId:%s, err:%s",
			svc.order.Order.PayOrderId, param.UserId, err)
		svc.engine.Rollback()
		return errdef.ORDER_GIFT_RECEIVE_FAIL
	}
	
	svc.engine.Commit()
	
	//code := svc.sms.GetSmsCode()
	// todo: 发送短信通知
	
	
	return errdef.SUCCESS
}
