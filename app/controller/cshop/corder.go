package cshop

import (
	"fmt"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mshop"
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
	resp.OrderId = fmt.Sprint("T", util.NewOrderId())
	// 订单商品流水mp
	//productMp := make(map[string]*models.OrderProduct, 0)
	// 存储购物车id
	cartIds := make([]int, 0)
	for index, item := range param.Products {
		// 购物车下单 需校验购物车数据 并在下单成功时 清理购物车
		if param.ReqType == 3 {
			condition := fmt.Sprintf("id=%d", item.CartId)
			cart, err := svc.shop.GetProductCart(condition)
			if err != nil {
				log.Log.Errorf("shop_trace: get product cart fail, cartId:%d err:%s", item.CartId, err)
				svc.engine.Rollback()
				return errdef.SHOP_GET_PRODUCT_CART_FAIL, nil
			}
			
			if cart.SkuId != item.SkuId || cart.ProductId != item.ProductId {
				log.Log.Errorf("shop_trace: invalid cartId, cartId:%d", item.CartId)
				svc.engine.Rollback()
				return errdef.INVALID_PARAMS, nil
			}
			
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
		
		// 订单商品快照
		//record := &models.OrderProduct{
		//	OrderId: resp.OrderId,
		//	UserId: resp.UserId,
		//	ProductId: item.ProductId,
		//	SkuId: item.SkuId,
		//
		//}
	}
	
	switch param.ReqType {
	case 1:
		// 查询
		addr, err := svc.shop.GetUserAddrById("")
		if err != nil {
			log.Log.Errorf("shop_trace: get user addr by id fail, err:%s", err)
		}
		
		if addr != nil && addr.IsDefault == 1 {
			resp.UserAddr = addr
			resp.UserAddr.Mobile = util.HideMobileNum(resp.UserAddr.Mobile)
		}
		
		// 事务回滚
		svc.engine.Rollback()
	case 2, 3:
		// 详情页/购物车下单
		addr, err := svc.shop.GetUserAddrById(fmt.Sprint(param.UserAddrId))
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
			affected, err := svc.CleanProductCart(cartIds)
			if int(affected) != len(cartIds) || err != nil {
				log.Log.Errorf("shop_trace: clean product cart fail, affected:%d, len:%d, err:%s", affected, len(cartIds), err)
				svc.engine.Rollback()
				return errdef.SHOP_PLACE_ORDER_FAIL, nil
			}
		}
		
		now := time.Now()
		resp.CreateAt = int(now.Unix())
		resp.CreateTm = now.Format(consts.FORMAT_TM)
		resp.PayDuration = consts.PAYMENT_DURATION
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
		
		svc.engine.Commit()
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
	
	if err = util.JsonFast.UnmarshalFromString(sku.OwnSpec, &item.SkuSpec); err != nil {
		log.Log.Errorf("shop_trace: unmarshal own spec fail, skuId:%d, err:%s", item.SkuId, err)
	}
	
	stockInfo, err := svc.shop.GetProductSkuStock(fmt.Sprint(item.SkuId))
	if err != nil {
		log.Log.Errorf("shop_trace: get product sku stock fail, skuId:%d, err:%s", item.SkuId, err)
		return errdef.ERROR
	}
	
	item.MaxBuy = stockInfo.MaxBuy
	item.MinBuy = stockInfo.MinBuy
	item.Stock = stockInfo.Stock
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
		item.Count = stockInfo.Stock
	}
	
	return errdef.SUCCESS
}

// 下单成功 清理购物车数据
func (svc *ShopModule) CleanProductCart(cartIds []int) (int64, error) {
	return svc.shop.DelProductCartByIds(cartIds)
}
