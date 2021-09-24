package cpay

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/app/config"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/morder"
	"fmt"
	"sports_service/server/models/mpay"
	"sports_service/server/models/muser"
	"sports_service/server/tools/alipay"
	"sports_service/server/tools/wechat"
	"time"
)

type PayModule struct {
	context     *gin.Context
	engine      *xorm.Session
	order       *morder.OrderModel
	user        *muser.UserModel
	pay         *mpay.PayModel
}

func New(c *gin.Context) PayModule {
	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
	appSocket := dao.AppEngine.NewSession()
	defer appSocket.Close()
	return PayModule{
		context: c,
		order: morder.NewOrderModel(venueSocket),
		user: muser.NewUserModel(appSocket),
		pay: mpay.NewPayModel(venueSocket),
		engine: venueSocket,
	}
}

func (svc *PayModule) AppPay(param *morder.PayReqParam) (int, interface{}) {
	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("pay_trace: user not found, userId:%s", param.UserId)
		return errdef.USER_NOT_EXISTS, nil
	}

	ok, err := svc.order.GetOrder(param.OrderId)
	if !ok || err != nil {
		log.Log.Errorf("pay_trace: order not found, orderId:%s", param.OrderId)
		return errdef.ORDER_NOT_EXISTS, nil
	}

	// 订单如果不是待支付状态
	if svc.order.Order.Status != consts.ORDER_TYPE_WAIT {
		log.Log.Errorf("pay_trace: order status fail, orderId:%s, status:%d", param.OrderId, svc.order.Order.Status)
		return errdef.ORDER_NOT_EXISTS, nil
	}

	ok, err = svc.pay.GetPaymentChannel(param.PayType)
	if !ok || err != nil {
		log.Log.Errorf("pay_trace: get payment channel fail, orderId:%s, ok:%v, err:%s", svc.order.Order.PayOrderId,
			ok, err)
		return errdef.PAY_CHANNEL_NOT_EXISTS, nil
	}

	switch svc.pay.PayChannel.Identifier {
	case consts.ALIPAY:
		// 支付宝
		payParam, err := svc.AliPay(svc.pay.PayChannel.AppId, svc.pay.PayChannel.PrivateKey)
		if err != nil {
			log.Log.Errorf("pay_trace: get alipay param fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			return errdef.PAY_ALI_PARAM_FAIL, nil
		}

		info := make(map[string]interface{}, 0)
		info["sign"] = payParam

		if _, err = svc.UpdateOrderPayType(param.OrderId); err != nil {
			log.Log.Errorf("pay_trace: update order payType by ali fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			return errdef.ORDER_UPDATE_FAIL, nil
		}

		return errdef.SUCCESS, info

	case consts.WEICHAT:
		// 微信
		mp, err := svc.WechatPay(svc.pay.PayChannel.AppId, svc.pay.PayChannel.AppKey, svc.pay.PayChannel.AppSecret)
		if err != nil {
			log.Log.Errorf("pay_trace: get wechatPay param fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			return errdef.PAY_WX_PARAM_FAIL, nil
		}

		if _, err = svc.UpdateOrderPayType(param.OrderId); err != nil {
			log.Log.Errorf("pay_trace: update order payType by wx fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			return errdef.ORDER_UPDATE_FAIL, nil
		}

		return errdef.SUCCESS, mp
	default:
		log.Log.Errorf("pay_trace: unsupported payType:%d", param.PayType)
	}


	return errdef.PAY_INVALID_TYPE, nil
}

// 更新订单支付类型
func (svc *PayModule) UpdateOrderPayType(orderId string) (int64, error) {
	svc.order.Order.PayChannelId = svc.pay.PayChannel.Id
	svc.order.Order.PayOrderId = orderId
	cols := "pay_channel_id"
	return svc.order.UpdateOrderInfo(cols)
}

func (svc *PayModule) AliPay(appId, privateKey string) (string, error) {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	client := alipay.NewAliPay(true, appId, privateKey)
	client.OutTradeNo = svc.order.Order.PayOrderId
	client.TotalAmount = fmt.Sprintf("%.2f", float64(svc.order.Order.Amount)/100)
	client.Subject = svc.order.Order.Subject
	client.TimeExpire = time.Unix(int64(svc.order.Order.CreateAt + consts.PAYMENT_DURATION), 0).In(cstSh).Format(consts.FORMAT_TM)
	payParam, err := client.TradeAppPay()
	if err != nil {
		return "", err
	}

	return payParam, nil
}

func (svc *PayModule) WechatPay(appId, merchantId, secret string) (map[string]interface{}, error) {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	client := wechat.NewWechatPay(true, appId, merchantId, secret)
	client.OutTradeNo = svc.order.Order.PayOrderId
	client.TotalAmount = svc.order.Order.Amount
	client.Subject = svc.order.Order.Subject
	client.NotifyUrl = config.Global.WechatNotifyUrl
	client.CreateIp = svc.context.ClientIP()
	client.TimeStart = time.Unix(int64(svc.order.Order.CreateAt), 0).In(cstSh).Format(consts.FORMAT_WX_TM)
	client.TimeExpire = time.Unix(int64(svc.order.Order.CreateAt + consts.PAYMENT_DURATION), 0).In(cstSh).Format(consts.FORMAT_WX_TM)
	mp, err := client.TradeAppPay()
	if err != nil {
		return nil, err
	}

	return mp, nil

}
