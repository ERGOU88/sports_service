package cshop

import (
	"errors"
	"fmt"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mshop"
	tc "sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"time"
)

func (svc *ShopModule) PlaceOrder(param *mshop.PlaceOrderReq) (int, *mshop.OrderResp) {
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR, nil
	}
	
	resp := &mshop.OrderResp{Products: make([]*mshop.Product, len(param.Products))}
	resp.UserId = param.UserId
	resp.ClientIp = param.ClientIp
	resp.OrderId = util.NewShopOrderId()
	// 默认详情页下单
	resp.ActionType = consts.ORDER_ACTION_TYPE_DETAIL
	resp.ChannelId = param.Channel
	// 存储购物车id
	cartIds := make([]int, 0)
	for index, item := range param.Products {
		// 购物车下单 需校验购物车数据 并在下单成功时 清理购物车
		if param.ReqType == 3 {
			condition := fmt.Sprintf("id=%d", item.CartId)
			cart, err := svc.shop.GetProductCart(condition)
			if err != nil || cart == nil {
				log.Log.Errorf("shop_trace: get product cart fail, cartId:%d err:%s", item.CartId, err)
				svc.engine.Rollback()
				return errdef.SHOP_GET_PRODUCT_CART_FAIL, nil
			}
			
			if cart.SkuId != item.SkuId || cart.ProductId != item.ProductId {
				log.Log.Errorf("shop_trace: invalid cartId, cartId:%d", item.CartId)
				svc.engine.Rollback()
				return errdef.INVALID_PARAMS, nil
			}
			// 购物车下单
			resp.ActionType = consts.ORDER_ACTION_TYPE_CART
			cartIds = append(cartIds, cart.Id)
		}
		
		info := &mshop.Product{
			SkuId: item.SkuId,
			ProductId: item.ProductId,
			Count: item.Count,
			CartId:  item.CartId,
			UserId: resp.UserId,
			OrderId: resp.OrderId,
		}
		
		if code := svc.OrderProcess(info); code != errdef.SUCCESS {
			svc.engine.Rollback()
			return code, nil
		}
		
		resp.IsEnough = true
		if info.IsEnough == false {
			resp.IsEnough = false
		}
		
		resp.Total += info.Count
		resp.PayAmount += info.PayAmount
		resp.OrderAmount += info.OrderAmount
		resp.DeliveryAmount += info.DeliveryAmount
		resp.DiscountAmount += info.DiscountAmount
		resp.ProductAmount += info.ProductAmount
		resp.Products[index] = info
	}
	
	switch param.ReqType {
	case 1:
		id := fmt.Sprint(param.UserAddrId)
		if param.UserAddrId == 0 {
			id = ""
		}
		
		// 查询
		addr, err := svc.shop.GetUserAddr(id, param.UserId)
		if err != nil {
			log.Log.Errorf("shop_trace: get user addr by id fail, err:%s", err)
		}
		
		if addr != nil {
			resp.UserAddr = addr
			resp.UserAddr.Mobile = util.HideMobileNum(resp.UserAddr.Mobile)
		}
		
		// 事务回滚
		svc.engine.Rollback()
	case 2, 3:
		// 详情页/购物车下单
		addr, err := svc.shop.GetUserAddr(fmt.Sprint(param.UserAddrId), param.UserId)
		if err != nil || addr == nil {
			log.Log.Errorf("shop_trace: get user addr by id fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.SHOP_USER_ADDR_NOT_FOUND, nil
		}
		
		resp.UserAddr = addr
		resp.UserAddr.Mobile = util.HideMobileNum(resp.UserAddr.Mobile)
		
		if !resp.IsEnough {
			svc.engine.Rollback()
			return errdef.SHOP_SKU_STOCK_NOT_ENOUGH, resp
		}
		
		if len(cartIds) > 0 {
			// 清理购物车对应的数据
			affected, err := svc.CleanProductCart(cartIds, param.UserId)
			if int(affected) != len(cartIds) || err != nil {
				log.Log.Errorf("shop_trace: clean product cart fail, affected:%d, len:%d, err:%s", affected, len(cartIds), err)
				svc.engine.Rollback()
				return errdef.SHOP_PLACE_ORDER_FAIL, nil
			}
		}
		
		now := time.Now()
		resp.CreateAt = int(now.Unix())
		resp.CreateTm = now.Format(consts.FORMAT_TM)
		resp.PayDuration = consts.SHOP_PAYMENT_DURATION
		// todo: 暂时写死
		resp.DeliveryTypeName = consts.DELIVERY_NAME
		resp.DeliveryTelephone = consts.DELIVERY_TELEPHONE
		if _, err := svc.AddOrder(resp); err != nil {
			log.Log.Errorf("shop_trace: add order fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.SHOP_PLACE_ORDER_FAIL, nil
		}
		
		affected, err := svc.AddOrderProduct(resp.Products)
		if int(affected) != len(resp.Products) || err != nil {
			log.Log.Errorf("shop_trace: add order product fail, affected:%d, len:%d, err:%s", affected, len(cartIds), err)
			svc.engine.Rollback()
			return errdef.SHOP_PLACE_ORDER_FAIL, nil
		}
		
		if _, err := svc.AddBuyerDeliveryInfo(resp); err != nil {
			log.Log.Errorf("shop_trace: add buyer delivery info fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.SHOP_PLACE_ORDER_FAIL, nil
		}
		
		if _, err := svc.shop.RecordOrderId(resp.OrderId); err != nil {
			log.Log.Errorf("shop_trace: record orderId fail, orderId:%s, err:%s", resp.OrderId, err)
			svc.engine.Rollback()
			return errdef.SHOP_PLACE_ORDER_FAIL, nil
		}
		
		svc.engine.Commit()
	default:
		log.Log.Errorf("shop_trace: invalid reqType, reqType:%d", param.ReqType)
		svc.engine.Rollback()
		return errdef.INVALID_PARAMS, nil
	}
	
	
	return errdef.SUCCESS, resp
}

// 添加订单
func (svc *ShopModule) AddOrder(resp *mshop.OrderResp) (int64, error) {
	str, _ := util.JsonFast.MarshalToString(resp)
	
	order := &models.Orders{
		OrderId:  resp.OrderId,
		Extra: str,
		UserId: resp.UserId,
		ProductAmount: resp.ProductAmount,
		DeliveryAmount: resp.DeliveryAmount,
		OrderAmount: resp.OrderAmount,
		DiscountAmount: resp.DiscountAmount,
		PayAmount: resp.PayAmount,
		OrderTypeName: "FPV无人机",
		OrderType: 1001,
		ChannelId: resp.ChannelId,
		ActionType: resp.ActionType,
		CreateAt: resp.CreateAt,
		UpdateAt: resp.CreateAt,
	}
	
	return svc.shop.AddOrder(order)
}

func (svc *ShopModule) AddOrderProduct(products []*mshop.Product) (int64, error) {
	return svc.shop.AddOrderProduct(products)
}

// 添加买家配送信息
func (svc *ShopModule) AddBuyerDeliveryInfo(resp *mshop.OrderResp) (int64, error) {
	info := &models.BuyerDeliveryInfo{
		OrderId: resp.OrderId,
		UserId: resp.UserId,
		Name: resp.UserAddr.Name,
		Mobile: resp.UserAddr.Mobile,
		Telephone: resp.UserAddr.Telephone,
		ProvinceCode: resp.UserAddr.ProvinceCode,
		CityCode: resp.UserAddr.CityCode,
		DistrictCode: resp.UserAddr.DistrictCode,
		CommunityCode: resp.UserAddr.CommunityCode,
		Address: resp.UserAddr.Address,
		FullAddress: resp.UserAddr.FullAddress,
		Longitude: resp.UserAddr.Longitude,
		Latitude: resp.UserAddr.Latitude,
		BuyerIp: resp.ClientIp,
		CreateAt: resp.CreateAt,
		UpdateAt: resp.CreateAt,
	}
	
	return svc.shop.AddBuyerDeliveryInfo(info)
}

// 下单流程
func (svc *ShopModule) OrderProcess(item *mshop.Product) int {
	if item.Count <= 0 {
		return errdef.INVALID_PARAMS
	}
	
	product, err := svc.shop.GetProductSpu(fmt.Sprint(item.ProductId))
	if err != nil {
		log.Log.Errorf("shop_trace: get product spu fail, productId:%d, err:%s", item.ProductId, err)
		return errdef.SHOP_PRODUCT_SPU_FAIL
	}
	
	item.ProductName = product.ProductName
	
	condition := fmt.Sprintf("id=%d AND product_id=%d", item.SkuId, item.ProductId)
	sku, err := svc.shop.GetProductSku(condition)
	if err != nil {
		log.Log.Errorf("shop_trace: get product sku by id fail, skuId:%d, err:%s", item.SkuId, err)
		return errdef.SHOP_PRODUCT_SKU_FAIL
	}
	
	item.SkuName = sku.Title
	item.SkuNo = sku.SkuNo
	item.SkuImage = tc.BucketURI(sku.SkuImage)
	item.CurPrice = sku.CurPrice
	item.MarketPrice = sku.MarketPrice
	item.IsFreeShip = sku.IsFreeShip
	item.DiscountPrice = sku.DiscountPrice
	item.StartTime = sku.StartTime
	item.EndTime = sku.EndTime
	item.Indexes = sku.Indexes
	now := int(time.Now().Unix())
	if now >= sku.StartTime && now < sku.EndTime {
		item.HasActivities = 1
		item.RemainDuration = item.EndTime - now
	}
	
	if sku.IsFreeShip == 1 {
		// todo: 不包邮的情况 记录邮费
		item.DeliveryAmount = 0
	}
	
	if sku.OwnSpec != "" {
		if err = util.JsonFast.UnmarshalFromString(sku.OwnSpec, &item.SkuSpec); err != nil {
			log.Log.Errorf("shop_trace: unmarshal own spec fail, skuId:%d, err:%s", item.SkuId, err)
		}
	} else {
		item.SkuSpec = make([]mshop.OwnSpec, 0)
	}
	
	stockInfo, err := svc.shop.GetProductSkuStock(fmt.Sprint(item.SkuId))
	if err != nil {
		log.Log.Errorf("shop_trace: get product sku stock fail, skuId:%d, err:%s", item.SkuId, err)
		return errdef.ERROR
	}
	
	item.MaxBuy = stockInfo.MaxBuy
	item.MinBuy = stockInfo.MinBuy
	item.Stock = stockInfo.Stock - stockInfo.PurchasedNum
	// 默认可够买
	item.CanBuy = true
	if item.Count > stockInfo.MaxBuy && stockInfo.MaxBuy > 0 {
		item.Count = stockInfo.MaxBuy
		item.CanBuy = false
	}
	
	if item.Count < stockInfo.MinBuy {
		item.Count = stockInfo.MinBuy
		item.CanBuy = false
	}
	
	price := item.CurPrice
	// 有活动
	if item.HasActivities == 1 {
		price = item.DiscountPrice
	}
	// 商品总价 = 数量 * 划线价格
	item.ProductAmount = item.Count * item.MarketPrice
	// 合计 = 商品总价 + 邮费
	item.OrderAmount = item.ProductAmount + item.DeliveryAmount
	// 优惠金额 = （划线价 - 当前价/活动价）* 数量
	item.DiscountAmount = (item.MarketPrice - price) * item.Count
	// 应付金额
	item.PayAmount = item.Count * price
	// 默认库存足够
	item.IsEnough = true
	// 更新库存
	affected, err := svc.shop.UpdateProductSkuStock(fmt.Sprint(item.SkuId), item.Count)
	if err != nil {
		log.Log.Errorf("shop_trace: update product sku stock fail, skuId:%s, err:%s", item.SkuId, err)
		return errdef.ERROR
	}
	
	if affected != 1 {
		item.IsEnough = false
		item.Count = stockInfo.Stock - stockInfo.PurchasedNum
	}
	
	return errdef.SUCCESS
}

// 购物车下单 成功时 清理购物车数据
func (svc *ShopModule) CleanProductCart(cartIds []int, userId string) (int64, error) {
	return svc.shop.DelProductCartByIds(cartIds, userId)
}

// 取消订单
// 1 修改订单状态
// 2 返还相应库存
// 3 从商品详情页下单 取消订单时 需添加到购物车
// 4 从购物车下单 取消时 无操作
func (svc *ShopModule) OrderCancel(param *mshop.ChangeOrderReq) int {
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR
	}
	
	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("shop_trace: user not exists, userId:%s", param.UserId)
		svc.engine.Rollback()
		return errdef.USER_NOT_EXISTS
	}
	
	order, err := svc.shop.GetOrder(param.OrderId)
	if err != nil {
		log.Log.Errorf("shop_trace: order not exists, orderId:%s", param.OrderId)
		svc.engine.Rollback()
		return errdef.SHOP_ORDER_NOT_EXISTS
	}
	
	if order.UserId != user.UserId {
		log.Log.Errorf("shop_trace: user not match, userId:%s, curUser:%s", order.UserId, user.UserId)
		svc.engine.Rollback()
		return errdef.SHOP_ORDER_CANCEL_FAIL
	}
	
	// 只有待支付状态订单可以取消
	if order.PayStatus != consts.SHOP_ORDER_TYPE_WAIT {
		log.Log.Errorf("shop_trace: order not allow cancel, orderId:%s, status:%d", order.OrderId, order.PayStatus)
		svc.engine.Rollback()
		return errdef.SHOP_ORDER_NOT_ALLOW_CANCEL
	}
	
	condition := fmt.Sprintf("`pay_status`=%d AND `order_id`='%s'", consts.SHOP_ORDER_TYPE_WAIT, order.OrderId)
	cols := "pay_status, close_time, update_at"
	order.PayStatus = consts.SHOP_ORDER_TYPE_UNPAID
	now := int(time.Now().Unix())
	order.UpdateAt = now
	order.CloseTime = now
	// 更新订单状态
	if _, err := svc.shop.UpdateOrderInfo(condition, cols, order); err != nil {
		log.Log.Errorf("shop_trace: update order info fail, orderId:%s, err:%v", order.OrderId, err)
		svc.engine.Rollback()
		return errdef.SHOP_ORDER_UPDATE_FAIL
	}
	
	if err := svc.CancelOrderProcess(order); err != nil {
		log.Log.Errorf("shop_trace: cancel order fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.SHOP_ORDER_CANCEL_FAIL
	}
	
	svc.engine.Commit()
	return errdef.SUCCESS
}

// 取消订单流程
func (svc *ShopModule) CancelOrderProcess(order *models.Orders) error {
	// 获取订单对应的商品流水
	list, err := svc.shop.GetOrderProductList(order.OrderId)
	if err != nil {
		log.Log.Errorf("shop_trace: get order products by orderId fail, orderId:%s, err:%s", order.OrderId, err)
		return err
	}
	
	for _, item := range list {
		// 归还商品sku库存
		affected, err := svc.shop.UpdateProductSkuStock(fmt.Sprint(item.SkuId), item.Count * -1)
		if affected != 1 || err != nil {
			log.Log.Errorf("shop_trace: update stock info fail, orderId:%s, err:%s, affected:%d, skuId:%s",
				order.OrderId, err, affected, item.SkuId)
			return errors.New("update stock info fail")
		}
		
	
		// 详情页下单 取消订单时 需添加到购物车
		//if order.ActionType == consts.ORDER_ACTION_TYPE_DETAIL {
			info := &models.ProductCart{
				UserId:    order.UserId,
				SkuId:     item.SkuId,
				Count:     item.Count,
				IsCheck:   0,
				ProductId: item.ProductId,
			}
			
			if _, err := svc.shop.AddProductCart(info); err != nil {
				return errors.New("add product cart fail")
			}
		//}
		
	}
	
	return nil
}

// 订单列表
func (svc *ShopModule) OrderList(userId, reqType string, page, size int) (int, []*mshop.OrderResp) {
	if userId == "" {
		return errdef.USER_NOT_EXISTS, nil
	}
	
	if user := svc.user.FindUserByUserid(userId); user == nil {
		return errdef.USER_NOT_EXISTS, nil
	}
	
	condition := svc.GetQueryCondition(reqType, userId)
	if condition == "" {
		return errdef.INVALID_PARAMS, nil
	}
	
	offset := (page - 1) * size
	list, err := svc.shop.GetOrderList(condition, offset, size)
	if err != nil {
		log.Log.Errorf("shop_trace: get order list fail, err:%s", err)
		return errdef.SHOP_ORDER_LIST_FAIL, nil
	}
	
	if len(list) == 0 {
		return errdef.SUCCESS, []*mshop.OrderResp{}
	}
	
	res := make([]*mshop.OrderResp, len(list))
	for index, item := range list {
		info, err := svc.OrderInfo(item)
		if err != nil {
			log.Log.Errorf("shop_trace: get order info fail, err:%s", err)
			return errdef.SHOP_ORDER_LIST_FAIL, nil
		}
		
		res[index] = info
	}
	
	
	return errdef.SUCCESS, res
}

func (svc *ShopModule) OrderDetail(userId, orderId string) (int, *mshop.OrderResp) {
	if userId == "" {
		return errdef.USER_NOT_EXISTS, nil
	}
	
	if user := svc.user.FindUserByUserid(userId); user == nil {
		return errdef.USER_NOT_EXISTS, nil
	}
	
	order, err := svc.shop.GetOrder(orderId)
	if err != nil {
		return errdef.SHOP_ORDER_NOT_EXISTS, nil
	}
	
	res, err := svc.OrderInfo(*order)
	if err != nil {
		log.Log.Errorf("shop_trace: get order info fail, err:%s", err)
		return errdef.SHOP_ORDER_LIST_FAIL, nil
	}
	
	return errdef.SUCCESS, res
}

// 订单信息
func (svc *ShopModule) OrderInfo(item models.Orders) (*mshop.OrderResp, error) {
		info := &mshop.OrderResp{}
		if err := util.JsonFast.UnmarshalFromString(item.Extra, &info); err != nil {
			log.Log.Errorf("shop_trace: unmarshal fail, orderId:%s, err:%s", item.OrderId, err)
			return nil, err
		}
		
		info.PayDuration = 0
		info.Status = item.PayStatus
		// 待支付订单 剩余支付时长
		if item.PayStatus == consts.SHOP_ORDER_TYPE_WAIT {
			// 已过时长 =  当前时间戳 - 订单创建时间戳
			duration := time.Now().Unix() - int64(item.CreateAt)
			// 订单状态是待支付 且 已过时长 <= 总时差
			if duration < consts.SHOP_PAYMENT_DURATION {
				log.Log.Debugf("order_trace: duration:%v", duration)
				// 剩余支付时长 = 总时长 - 已过时长
				info.PayDuration = consts.SHOP_PAYMENT_DURATION - duration
			}
		}
		
		if item.PayStatus == consts.SHOP_ORDER_TYPE_PAID {
			info.Status = svc.GetDeliveryStatus(item.DeliveryStatus)
		}
		
		info.PayStatus = item.PayStatus
		info.DeliveryCode = item.DeliveryCode
		info.DeliveryStatus = item.DeliveryStatus
		info.DeliveryTelephone = item.DeliveryTelephone
		info.DeliveryTypeName = item.DeliveryTypeName
	
	return info, nil
}

func (svc *ShopModule) GetDeliveryStatus(deliveryStatus int) int {
	switch deliveryStatus {
	case consts.NOT_DELIVERED:
		return 2
	case consts.HAS_DELIVERED:
		return 3
	case consts.HAS_SIGNED:
		return 4
	}
	
	return 0
}

func (svc *ShopModule) GetQueryCondition(reqType, userId string) string {
	var condition string
	switch reqType {
	// 查看全部订单
	case "1":
		condition = fmt.Sprintf("user_id='%s' AND is_delete=0", userId)
	// 待付款订单
	case "2":
		condition = fmt.Sprintf("user_id='%s' AND is_delete=0 AND pay_status=0", userId)
	// 待收货订单
	case "3":
		condition = fmt.Sprintf("user_id='%s' AND is_delete=0 AND pay_status=2 AND delivery_status in(0, 1)", userId)
	// 已完成订单
	case "4":
		condition = fmt.Sprintf("user_id='%s' AND is_delete=0 AND pay_status=2 AND delivery_status=2", userId)
	default:
		return ""
	}
	
	return condition
}

// 订单确认收货
func (svc *ShopModule) ConfirmReceipt(param *mshop.ChangeOrderReq) int {
	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("shop_trace: user not exists, userId:%s", param.UserId)
		return errdef.USER_NOT_EXISTS
	}
	
	order, err := svc.shop.GetOrder(param.OrderId)
	if err != nil {
		log.Log.Errorf("shop_trace: order not exists, orderId:%s", param.OrderId)
		return errdef.SHOP_ORDER_NOT_EXISTS
	}
	
	if order.UserId != user.UserId {
		log.Log.Errorf("shop_trace: user not match, userId:%s, curUser:%s", order.UserId, user.UserId)
		return errdef.SHOP_CONFIRM_RECEIPT_FAIL
	}
	
	// 订单 != 支付成功 || 配送状态 != 已配送
	if order.PayStatus != consts.SHOP_ORDER_TYPE_PAID || order.DeliveryStatus != consts.HAS_DELIVERED {
		log.Log.Errorf("shop_trace: not allow confirm,orderId:%s, payStatus:%d, deliveryStatus:%d", order.OrderId,
			order.PayStatus, order.DeliveryStatus)
		return errdef.SHOP_NOT_ALLOW_CONFIRM
	}
	
	now := int(time.Now().Unix())
	// 支付状态=已支付 && 配送状态=已配送
	condition := fmt.Sprintf("order_id='%s' AND pay_status=%d AND delivery_status=1", order.OrderId, consts.SHOP_ORDER_TYPE_PAID)
	cols := "update_at, sign_time, finish_time, delivery_status"
	order.UpdateAt = now
	order.SignTime = now
	// todo: 暂时已收货 即表示 已完成
	order.FinishTime = now
	// 配送状态 修改为 已签收
	order.DeliveryStatus = consts.HAS_SIGNED
	if _, err := svc.shop.UpdateOrderInfo(condition, cols, order); err != nil {
		log.Log.Errorf("shop_trace: update order info fail, orderId:%s, err:%s", order.OrderId, err)
		return errdef.SHOP_CONFIRM_RECEIPT_FAIL
	}
	
	return errdef.SUCCESS
}

// 删除订单
func (svc *ShopModule) DeleteOrder(param *mshop.ChangeOrderReq) int {
	user := svc.user.FindUserByUserid(param.UserId)
	if user == nil {
		log.Log.Errorf("shop_trace: user not exists, userId:%s", param.UserId)
		return errdef.USER_NOT_EXISTS
	}
	
	order, err := svc.shop.GetOrder(param.OrderId)
	if err != nil {
		log.Log.Errorf("shop_trace: order not exists, orderId:%s", param.OrderId)
		return errdef.SHOP_ORDER_NOT_EXISTS
	}
	
	if order.UserId != user.UserId {
		log.Log.Errorf("order_trace: user not match, userId:%s, curUser:%s", order.UserId, user.UserId)
		return errdef.SHOP_ORDER_DELETE_FAIL
	}
	
	// 已支付 且 未完成的订单 不能删除
	if order.PayStatus == consts.SHOP_ORDER_TYPE_PAID && order.DeliveryStatus != consts.HAS_SIGNED {
		return errdef.SHOP_ORDER_DELETE_FAIL
	}
	
	condition := fmt.Sprintf("order_id='%s'", order.OrderId)
	order.IsDelete = 1
	order.UpdateAt = int(time.Now().Unix())
	cols := "is_delete, update_at"
	affected, err := svc.shop.UpdateOrderInfo(condition, cols, order)
	if affected != 1 || err != nil {
		log.Log.Errorf("shop_trace: update order info fail, err:%s, affected:%d", err, affected)
		return errdef.SHOP_ORDER_DELETE_FAIL
	}
	
	return errdef.SUCCESS
}
