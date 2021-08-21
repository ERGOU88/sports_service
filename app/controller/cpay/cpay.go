package cpay

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/morder"
	"fmt"
	"sports_service/server/models/muser"
	"sports_service/server/tools/alipay"
	"sports_service/server/tools/wechat"
)

type PayModule struct {
	context     *gin.Context
	engine      *xorm.Session
	order       *morder.OrderModel
	user        *muser.UserModel
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
	if svc.order.Order.Status != consts.PAY_TYPE_WAIT {
		log.Log.Errorf("pay_trace: order status fail, orderId:%s, status:%d", param.OrderId, svc.order.Order.Status)
		return errdef.ORDER_NOT_EXISTS, nil
	}

	switch param.PayType {
	case consts.ALIPAY:
		// 支付宝
		payParam, err := svc.AliPay()
		if err != nil {
			log.Log.Errorf("pay_trace: get alipay param fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			return errdef.PAY_ALI_PARAM_FAIL, nil
		}

		info := make(map[string]interface{}, 0)
		info["sign"] = payParam
		return errdef.SUCCESS, info

	case consts.WEICHAT:
		// 微信
		mp, err := svc.WechatPay()
		if err != nil {
			log.Log.Errorf("pay_trace: get wechatPay param fail, orderId:%s, err:%s", svc.order.Order.PayOrderId, err)
			return errdef.PAY_WX_PARAM_FAIL, nil
		}

		return errdef.SUCCESS, mp
	default:
		log.Log.Errorf("pay_trace: unsupported payType:%d", param.PayType)
	}


	return errdef.PAY_INVALID_TYPE, nil
}

func (svc *PayModule) AliPay() (string, error) {
	client := alipay.NewAliPay(true)
	client.OutTradeNo = svc.order.Order.PayOrderId
	client.TotalAmount = fmt.Sprintf("%.2f", float64(svc.order.Order.Amount)/100)
	client.Subject = svc.order.Order.Subject
	payParam, err := client.TradeAppPay()
	if err != nil {
		return "", err
	}

	return payParam, nil
}

func (svc *PayModule) WechatPay() (map[string]interface{}, error) {
	client := wechat.NewWechatPay(true)
	client.OutTradeNo = svc.order.Order.PayOrderId
	client.TotalAmount = svc.order.Order.Amount
	client.Subject = svc.order.Order.Subject
	client.NotifyUrl = "http://fpv-app-api-qa.bluetrans.cn/api/v1/pay/wechat/notify"
	client.CreateIp = "127.0.0.1"
	mp, err := client.TradeAppPay()
	if err != nil {
		return nil, err
	}

	return mp, nil

}
