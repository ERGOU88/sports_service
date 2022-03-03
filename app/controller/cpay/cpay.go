package cpay

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/app/config"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/morder"
	"fmt"
	"sports_service/server/models/mpay"
	"sports_service/server/models/mshop"
	"sports_service/server/models/muser"
	"sports_service/server/tools/alipay"
	"sports_service/server/tools/wechat"
	"time"
	"errors"
)

type PayModule struct {
	context     *gin.Context
	engine      *xorm.Session
	order       *morder.OrderModel
	user        *muser.UserModel
	pay         *mpay.PayModel
	shop        *mshop.ShopModel
	social      *muser.SocialModel
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
		shop: mshop.NewShop(appSocket),
		social: muser.NewSocialPlatform(appSocket),
		engine: venueSocket,
	}
}

func (svc *PayModule) InitiatePayment(param *morder.PayReqParam) (int, interface{}) {
	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("pay_trace: user not found, userId:%s", param.UserId)
		return errdef.USER_NOT_EXISTS, nil
	}
	
	ok, err := svc.pay.GetPaymentChannel(param.PayType)
	if !ok || err != nil {
		log.Log.Errorf("pay_trace: get payment channel fail, orderId:%s, ok:%v, err:%s", param.OrderId,
			ok, err)
		return errdef.PAY_CHANNEL_NOT_EXISTS, nil
	}
	
	length := len(param.OrderId)
	switch length {
	// 场馆订单
	case 16:
		ok, err := svc.order.GetOrder(param.OrderId)
		if !ok || err != nil {
			log.Log.Errorf("pay_trace: order not found, orderId:%s", param.OrderId)
			return errdef.ORDER_NOT_EXISTS, nil
		}
		
		param.Subject = svc.order.Order.Subject
		param.Amount = svc.order.Order.Amount
		param.CreateAt = svc.order.Order.CreateAt
		
		// 订单如果不是待支付状态
		if svc.order.Order.Status != consts.ORDER_TYPE_WAIT {
			log.Log.Errorf("pay_trace: order status fail, orderId:%s, status:%d", param.OrderId, svc.order.Order.Status)
			return errdef.ORDER_NOT_EXISTS, nil
		}
		
		if _, err = svc.UpdateOrderPayTypeByVenue(param.OrderId); err != nil {
			log.Log.Errorf("pay_trace: update order payType fail, orderId:%s, err:%s", param.OrderId, err)
			return errdef.ORDER_UPDATE_FAIL, nil
		}
		
	// 商城订单
	case 18:
		order, err := svc.shop.GetOrder(param.OrderId)
		if order == nil || err != nil {
			log.Log.Errorf("pay_trace: shop order not found, orderId:%s, err:%s", param.OrderId, err)
			return errdef.SHOP_ORDER_NOT_EXISTS, nil
		}
		
		param.Subject = order.OrderTypeName
		param.Amount = order.PayAmount
		param.CreateAt = order.CreateAt
		
		if order.PayStatus != consts.SHOP_ORDER_TYPE_WAIT {
			log.Log.Errorf("pay_trace: shop order status fail, orderId:%s, status:%d", param.OrderId, order.PayStatus)
			return errdef.SHOP_ORDER_NOT_EXISTS, nil
		}
		
		if _, err := svc.UpdateOrderPayTypeByShop(order); err != nil {
			log.Log.Errorf("pay_trace: update order payType fail, orderId:%s, err:%s", param.OrderId, err)
			return errdef.SHOP_ORDER_UPDATE_FAIL, nil
		}
	}

	return svc.GetPaymentParams(param)
}

// 获取支付参数
func (svc *PayModule) GetPaymentParams(param *morder.PayReqParam) (int, interface{}) {
	switch svc.pay.PayChannel.Identifier {
	case consts.ALIPAY:
		// 支付宝
		payParam, err := svc.AliPay(svc.pay.PayChannel.AppId, svc.pay.PayChannel.PrivateKey, param)
		if err != nil {
			log.Log.Errorf("pay_trace: get alipay param fail, orderId:%s, err:%s", param.OrderId, err)
			return errdef.PAY_ALI_PARAM_FAIL, nil
		}
		
		info := make(map[string]interface{}, 0)
		info["sign"] = payParam
		
		return errdef.SUCCESS, info
	
	case consts.WEICHAT:
		openId := ""
		switch param.Platform {
		// app
		case 0:
			// 微信
			mp, err := svc.WechatPay(svc.pay.PayChannel.AppId, svc.pay.PayChannel.AppKey, svc.pay.PayChannel.AppSecret, openId, param)
			if err != nil {
				log.Log.Errorf("pay_trace: get wechatPay param fail, orderId:%s, err:%s", param.OrderId, err)
				return errdef.PAY_WX_PARAM_FAIL, nil
			}
			
			return errdef.SUCCESS, mp
		// 小程序
		case 1:
			ok, err := svc.social.GetSocialAccount(consts.TYPE_APPLET, svc.user.User.UserId)
			if ok && err == nil {
				openId = svc.social.SocialAccount.OpenId
			}
			
			// 微信
			mp, err := svc.WechatPay(wechat.APPLET_APPID, svc.pay.PayChannel.AppKey, svc.pay.PayChannel.AppSecret, openId, param)
			if err != nil {
				log.Log.Errorf("pay_trace: get wechatPay param fail, orderId:%s, err:%s", param.OrderId, err)
				return errdef.PAY_WX_PARAM_FAIL, nil
			}
			
			return errdef.SUCCESS, mp
		case 2:
			mp, err := svc.WechatPay(svc.pay.PayChannel.AppId, svc.pay.PayChannel.AppKey, svc.pay.PayChannel.AppSecret, openId, param)
			if err != nil {
				log.Log.Errorf("pay_trace: get wechatPay param fail, orderId:%s, err:%s", param.OrderId, err)
				return errdef.PAY_WX_PARAM_FAIL, nil
			}
			
			return errdef.SUCCESS, mp
		}
		
	default:
		log.Log.Errorf("pay_trace: unsupported payType:%d", param.PayType)
	}
	
	return errdef.PAY_INVALID_TYPE, nil
}

// 更新场馆订单支付类型
func (svc *PayModule) UpdateOrderPayTypeByVenue(orderId string) (int64, error) {
	svc.order.Order.PayChannelId = svc.pay.PayChannel.Id
	svc.order.Order.PayOrderId = orderId
	cols := "pay_channel_id"
	return svc.order.UpdateOrderInfo(cols)
}

// 更新商城订单支付类型
func (svc *PayModule) UpdateOrderPayTypeByShop(order *models.Orders) (int64, error) {
	condition := fmt.Sprintf("`pay_status`=%d AND `order_id`='%s'", consts.SHOP_ORDER_TYPE_WAIT, order.OrderId)
	cols := "pay_type, update_at"
	order.PayStatus = consts.SHOP_ORDER_TYPE_PAID
	now := int(time.Now().Unix())
	order.UpdateAt = now
	order.PayChannelId = svc.pay.PayChannel.Id
	// 更新订单状态
	return svc.shop.UpdateOrderInfo(condition, cols, order)
}

func (svc *PayModule) AliPay(appId, privateKey string, param *morder.PayReqParam) (string, error) {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	client := alipay.NewAliPay(true, appId, privateKey)
	client.OutTradeNo = param.OrderId
	client.TotalAmount = fmt.Sprintf("%.2f", float64(param.Amount)/100)
	client.Subject = param.Subject
	client.TimeExpire = time.Unix(int64(param.CreateAt + consts.PAYMENT_DURATION), 0).In(cstSh).Format(consts.FORMAT_TM)
	payParam, err := client.TradeAppPay()
	if err != nil {
		return "", err
	}

	return payParam, nil
}

func (svc *PayModule) WechatPay(appId, merchantId, secret, openId string, param *morder.PayReqParam) (map[string]interface{}, error) {
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	client := wechat.NewWechatPay(true, appId, merchantId, secret)
	client.OutTradeNo = param.OrderId
	client.TotalAmount = param.Amount
	client.Subject = param.Subject
	client.NotifyUrl = config.Global.WechatNotifyUrl
	client.OpenId = openId
	client.CreateIp = svc.context.ClientIP()
	client.TimeStart = time.Unix(int64(param.CreateAt), 0).In(cstSh).Format(consts.FORMAT_WX_TM)
	client.TimeExpire = time.Unix(int64(param.CreateAt + consts.PAYMENT_DURATION), 0).In(cstSh).Format(consts.FORMAT_WX_TM)
	
	// 小程序支付
	switch param.Platform {
	case 0:
		return client.TradeAppPay()
	case 1:
		return client.TradeJsAPIPay()
	case 2:
		return client.TradeH5Pay()
		
	}
	
	return nil, errors.New("invalid platform")
}
