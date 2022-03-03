package cshop

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/backend/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mpay"
	"sports_service/server/models/mshop"
	"sports_service/server/models/muser"
	"sports_service/server/tools/wechat"
	"sports_service/server/util"
	"strconv"
	"strings"
	"time"
	"sports_service/server/tools/alipay"
)

type ShopModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	shop        *mshop.ShopModel
	pay         *mpay.PayModel
}

func New(c *gin.Context) ShopModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	venueSocket := dao.AppEngine.NewSession()
	defer socket.Close()
	return ShopModule{
		context: c,
		user: muser.NewUserModel(socket),
		shop: mshop.NewShop(socket),
		pay: mpay.NewPayModel(venueSocket),
		engine: socket,
	}
}

func (svc *ShopModule) GetProductList(sortType, keyword string, page, size int) (int, int64, []*mshop.ProductSimpleInfo) {
	offset := (page - 1) * size
	list, count, err := svc.shop.GetSpuList(sortType, keyword, offset, size)
	if err != nil {
		log.Log.Errorf("shop_trace: get all spu fail, err:%s", err)
		return errdef.SHOP_GET_ALL_SPU_FAIL, count, nil
	}
	
	if len(list) == 0 {
		return errdef.SUCCESS, count, []*mshop.ProductSimpleInfo{}
	}
	
	for _, item := range list {
		skuList, err := svc.shop.GetProductSkuList(fmt.Sprint(item.Id))
		if err != nil {
			log.Log.Errorf("shop_trace: get product sku by spuId fail, spuId:%d, err:%s", item.Id, err)
			continue
		}
		
		item.SkuNum = len(skuList)
		condition := fmt.Sprintf("product_id=%d", item.Id)
		related, err := svc.shop.GetProductCategoryRelated(condition)
		if related != nil && err == nil {
			item.CategoryId = related.CategoryId
			item.CategoryName = related.CategoryName
		}
	}
	
	return errdef.SUCCESS, count, list
}

// 获取spu总数
func (svc *ShopModule) GetSpuTotal() {

}

func (svc *ShopModule) GetProductCategoryConf() []*mshop.Category {
	err, conf := svc.shop.GetProductCategoryByBackend()
	if err != nil || conf == nil {
		return []*mshop.Category{}
	}
	
	return conf
}

func (svc *ShopModule) AddCategory(params *models.ProductCategory) int {
	if _, err := svc.shop.AddProductCategory(params); err != nil {
		log.Log.Errorf("shop_trace: add category fail, err:%s", err)
		return errdef.SHOP_ADD_CATEGORY_FAIL
	}
	
	svc.shop.CleanCategoryInfoByMem()
	
	return errdef.SUCCESS
}

func (svc *ShopModule) EditCategory(params *models.ProductCategory) int {
	if _, err := svc.shop.UpdateProductCategory(params); err != nil {
		log.Log.Errorf("shop_trace: update category fail, err:%s", err)
		return errdef.SHOP_EDIT_CATEGORY_FAIL
	}
	
	svc.shop.CleanCategoryInfoByMem()
	return errdef.SUCCESS
}

func (svc *ShopModule) GetServiceList() (int, []mshop.AfterService) {
	list, err := svc.shop.GetServiceList()
	if err != nil {
		log.Log.Errorf("shop_trace: get service list fail, err:%s", err)
		return errdef.SHOP_GET_SERVICE_FAIL, list
	}
	
	return errdef.SUCCESS, list
}

func (svc *ShopModule) AddService(info *models.ShopServiceConf) int {
	if _, err := svc.shop.AddService(info); err != nil {
		return errdef.SHOP_ADD_SERVICE_FAIL
	}
	
	return errdef.SUCCESS
}

func (svc *ShopModule) UpdateService(info *models.ShopServiceConf) int {
	if _, err := svc.shop.UpdateService(info); err != nil {
		return errdef.SHOP_UPDATE_SERVICE_FAIL
	}
	
	return errdef.SUCCESS
}

func (svc *ShopModule) DelService(id string) int {
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR
	}
	
	if _, err := svc.shop.DelService(id); err != nil {
		svc.engine.Rollback()
		return errdef.SHOP_DEL_SERVICE_FAIL
	}
	
	condition := fmt.Sprintf("service_id=%s", id)
	if _, err := svc.shop.DelProductService(condition); err != nil {
		log.Log.Errorf("shop_trace: del product service fail, id:%s, err:%s", id, err)
		svc.engine.Rollback()
		return errdef.SHOP_DEL_SERVICE_FAIL
	}
	
	svc.engine.Commit()
	return errdef.SUCCESS
}

// 商品详情 获取商品实体信息
func (svc *ShopModule) GetProductDetail(productId string) (int, *mshop.ProductDetail) {
	detail, err := svc.shop.GetProductById(productId)
	if err != nil {
		log.Log.Errorf("shop_trace: get product spu fail, productId:%s, err:%s", productId, err)
		return errdef.SHOP_PRODUCT_SPU_FAIL, nil
	}
	
	condition := fmt.Sprintf("product_id=%s", productId)
	related, err := svc.shop.GetProductCategoryRelated(condition)
	if related != nil && err == nil {
		detail.CategoryId = related.CategoryId
		detail.CategoryName = related.CategoryName
	} else {
		log.Log.Errorf("shop_trace: get product category related fail, err:%s", err)
	}
	
	skuList, err := svc.shop.GetProductSkuList(fmt.Sprint(detail.Id))
	if err != nil {
		log.Log.Errorf("shop_trace: get product sku by spuId fail, spuId:%d, err:%s", detail.Id, err)
		return errdef.SHOP_PRODUCT_SKU_FAIL, nil
	}
	
	if skuList != nil {
		detail.SkuList = skuList
		for _, item := range detail.SkuList {
			stockInfo, err := svc.shop.GetProductSkuStock(fmt.Sprint(item.Id))
			if err == nil {
				item.Stock = stockInfo.Stock
				item.MinBuy = stockInfo.MinBuy
				item.MaxBuy = stockInfo.MaxBuy
			}
			
			if item.OwnSpec == nil {
				item.OwnSpec = make([]*mshop.OwnSpec, 0)
			}
		}
		
	} else {
		detail.SkuList = make([]*mshop.ProductSkuInfo, 0)
	}
	
	if detail.Specifications == nil {
		detail.Specifications = make([]*mshop.SpecInfo, 0)
	}
	
	if detail.SpecTemplate == nil {
		detail.SpecTemplate = make([]*mshop.SpecTemplate, 0)
	}
	
	services, err := svc.shop.GetProductServiceInfo(fmt.Sprint(detail.Id))
	if services != nil && err == nil {
		detail.AfterService = services
	} else {
		log.Log.Errorf("shop_trace: get product service fail, err:%s", err)
		detail.AfterService = []*mshop.AfterService{}
	}
	
	return errdef.SUCCESS, detail
}


// 添加商品分类规格
func (svc *ShopModule) AddCategorySpec(params *mshop.AddOrEditCategorySpecReq) int {
	now := int(time.Now().Unix())
	info := &models.ProductSpecification{}
	info.CategoryId = params.CategoryId
	info.CreateAt = now
	info.UpdateAt = now
	spec, err := util.JsonFast.MarshalToString(params.SpecInfo)
	if err != nil {
		log.Log.Errorf("shop_trace: marshal to string fail, err:%s", err)
		return errdef.SHOP_ADD_CATEGORY_SPEC_FAIL
	}
	
	info.Specifications = spec
	if _, err := svc.shop.AddCategorySpec(info); err != nil {
		log.Log.Errorf("shop_trace: add category spec fail, err:%s", err)
		return errdef.SHOP_ADD_CATEGORY_SPEC_FAIL
	}
	
	return errdef.SUCCESS
}

func (svc *ShopModule) EditCategorySpec(params *mshop.AddOrEditCategorySpecReq) int {
	now := int(time.Now().Unix())
	info := &models.ProductSpecification{}
	info.CategoryId = params.CategoryId
	info.UpdateAt = now
	spec, err := util.JsonFast.MarshalToString(params.SpecInfo)
	if err != nil {
		log.Log.Errorf("shop_trace: marshal to string fail, err:%s", err)
		return errdef.SHOP_EDIT_CATEGORY_SPEC_FAIL
	}
	
	info.Specifications = spec
	if _, err := svc.shop.UpdateCategorySpec(info); err != nil {
		log.Log.Errorf("shop_trace: update category spec fail, err:%s", err)
		return errdef.SHOP_EDIT_CATEGORY_SPEC_FAIL
	}

	return errdef.SUCCESS
}

func (svc *ShopModule) DelCategorySpec(categoryId string) int {
	if _, err := svc.shop.DelCategorySpec(categoryId); err != nil {
		return errdef.SHOP_DEL_CATEGORY_SPEC_FAIL
	}
	
	return errdef.SUCCESS
}

// 获取分类规格
func (svc *ShopModule) GetCategorySpec(categoryId string) (*mshop.CategorySpecInfo, int) {
	spec, err := svc.shop.GetCategorySpec(categoryId)
	if err != nil {
		log.Log.Errorf("shop_trace: get category spec fail, err:%s", err)
		return nil, errdef.SHOP_GET_SPEC_FAIL
	}
	
	return spec, errdef.SUCCESS
}

func (svc *ShopModule) GetCategorySpecList() (int, []*mshop.CategorySpecInfo) {
	list, err := svc.shop.GetCategorySpecList()
	if err != nil {
		return errdef.SHOP_GET_SPEC_FAIL, nil
	}
	
	if len(list) == 0 {
		return errdef.SUCCESS, []*mshop.CategorySpecInfo{}
	}
	
	for _, item := range list {
		category := svc.shop.GetProductCategoryInfoById(fmt.Sprint(item.CategoryId))
		if category != nil {
			item.CategoryName = category.CategoryName
		}
	}
	
	return errdef.SUCCESS, list
}

func (svc *ShopModule) AddProduct(params *mshop.AddOrEditProductReq) int {
	if len(params.SkuList) == 0 {
		return errdef.INVALID_PARAMS
	}
	
	productDetail, _ := util.JsonFast.MarshalToString(params.ProductDetail)
	introduction, _ := util.JsonFast.MarshalToString(params.Introduction)
	globalSpec, _ := util.JsonFast.MarshalToString(params.Specifications)
	spec, _ := util.JsonFast.MarshalToString(params.SpecTemplate)
	//service, _ := util.JsonFast.MarshalToString(params.AfterService)
	now := int(time.Now().Unix())
	spu := &models.Products{
		ProductName:  params.ProductName,
		ProductImage:  params.ProductImage,
		ProductDetail: productDetail,
		Status: params.Status,
		IsFreeShip: params.IsFreeShip,
		IsDelete: params.IsDelete,
		Introduction: introduction,
		Keywords: params.Keywords,
		Sortorder: params.Sortorder,
		VideoUrl: params.VideoUrl,
		CurPrice: params.CurPrice,
		MarketPrice: params.MarketPrice,
		IsRecommend: params.IsRecommend,
		IsTop: params.IsTop,
		IsCream: params.IsCream,
		Specifications: globalSpec,
		SpecTemplate: spec,
		//AfterService: service,
		CreateAt: now,
		UpdateAt: now,
	}
	
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR
	}
	
	if _, err := svc.shop.AddProductSpu(spu); err != nil {
		log.Log.Errorf("shop_trace: add product spu fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.SHOP_ADD_SPU_FAIL
	}
	
	category := svc.shop.GetProductCategoryInfoById(fmt.Sprint(params.CategoryId))
	if category == nil {
		log.Log.Errorf("shop_trace: get category fail, categoryId:%d", params.CategoryId)
		svc.engine.Rollback()
		return errdef.SHOP_GET_CATEGORY_FAIL
	}
	
	code := svc.AddSkuListInfo(params, spu.Id, category.CategoryName)
	if code != errdef.SUCCESS {
		svc.engine.Rollback()
		return code
	}
	
	if len(params.AfterService) > 0 {
		service := make([]*models.ProductService, len(params.AfterService))
		for index, id := range params.AfterService {
			info := &models.ProductService{
				ProductId: spu.Id,
				ServiceId: id,
			}
			
			service[index] = info
		}
		
		if _, err := svc.shop.AddProductService(service); err != nil {
			log.Log.Errorf("shop_trace: add product service fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.SHOP_ADD_PRODUCT_SVC_FAIL
		}
	}
	
	svc.engine.Commit()
	return code
}

func (svc *ShopModule) AddSkuListInfo(params *mshop.AddOrEditProductReq, productId int, categoryName string) int {
	now := int(time.Now().Unix())
	stockInfo := make([]*models.ProductSkuStock, len(params.SkuList))
	for index, item := range params.SkuList {
		images, _ := util.JsonFast.MarshalToString(item.Images)
		ownSpec, _ := util.JsonFast.MarshalToString(item.OwnSpec)
		sku := &models.ProductSku{
			ProductId: int64(productId),
			Title: item.Title,
			SkuImage: string(item.SkuImage),
			SkuNo: item.SkuNo,
			Images: images,
			CurPrice: item.CurPrice,
			MarketPrice: item.MarketPrice,
			IsFreeShip: item.IsFreeShip,
			Indexes: item.Indexes,
			OwnSpec: ownSpec,
			Status: item.Status,
			CreateAt: now,
			UpdateAt: now,
			Sortorder: item.Sortorder,
		}
		
		if _, err := svc.shop.AddProductSku(sku); err != nil {
			log.Log.Errorf("shop_trace: add product sku fail, err:%s", err)
			return errdef.SHOP_ADD_SKU_FAIL
		}
		
		stock := &models.ProductSkuStock{
			SkuId: sku.Id,
			Stock: item.Stock,
			MinBuy: 1,
			MaxBuy: 200,
			CreateAt: now,
			UpdateAt: now,
			ProductId: productId,
		}
		
		stockInfo[index] = stock
	}
	
	related := &models.ProductCategoryRelated{
		ProductId: productId,
		CategoryId: params.CategoryId,
		CreateAt: now,
		CategoryName: categoryName,
	}
	if _, err := svc.shop.AddProductCategoryRelated(related); err != nil {
		log.Log.Errorf("shop_trace: add related fail, err:%s", err)
		return errdef.SHOP_ADD_RELATED_FAIL
	}
	
	affected, err := svc.shop.AddProductSkuStock(stockInfo)
	if int(affected) != len(params.SkuList) || err != nil {
		log.Log.Errorf("shop_trace: add product sku stock fail, err:%s", err)
		return errdef.SHOP_ADD_SKU_STOCK_FAIL
	}
	
	return errdef.SUCCESS
}

func (svc *ShopModule) EditProduct(params *mshop.AddOrEditProductReq) int {
	if len(params.SkuList) == 0 {
		return errdef.INVALID_PARAMS
	}
	
	productDetail, _ := util.JsonFast.MarshalToString(params.ProductDetail)
	introduction, _ := util.JsonFast.MarshalToString(params.Introduction)
	globalSpec, _ := util.JsonFast.MarshalToString(params.Specifications)
	spec, _ := util.JsonFast.MarshalToString(params.SpecTemplate)
	now := int(time.Now().Unix())
	spu := &models.Products{
		Id: params.Id,
		ProductName:  params.ProductName,
		ProductImage:  params.ProductImage,
		ProductDetail: productDetail,
		Status: params.Status,
		IsFreeShip: params.IsFreeShip,
		IsDelete: params.IsDelete,
		Introduction: introduction,
		Keywords: params.Keywords,
		Sortorder: params.Sortorder,
		VideoUrl: params.VideoUrl,
		CurPrice: params.CurPrice,
		MarketPrice: params.MarketPrice,
		IsRecommend: params.IsRecommend,
		IsTop: params.IsTop,
		IsCream: params.IsCream,
		Specifications: globalSpec,
		SpecTemplate: spec,
		CreateAt: now,
		UpdateAt: now,
	}
	
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR
	}
	
	if _, err := svc.shop.UpdateProductSpu(spu); err != nil {
		log.Log.Errorf("shop_trace: update product spu fail, err:%s", err)
		svc.engine.Rollback()
		return errdef.SHOP_UPDATE_SPU_FAIL
	}
	
	if len(params.AfterService) > 0 {
		service := make([]*models.ProductService, len(params.AfterService))
		for index, id := range params.AfterService {
			info := &models.ProductService{
				ProductId: spu.Id,
				ServiceId: id,
			}
			
			service[index] = info
		}
		
		condition := fmt.Sprintf("product_id=%d", params.Id)
		if _, err := svc.shop.DelProductService(condition); err != nil {
			log.Log.Errorf("shop_trace: del product service fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.SHOP_ADD_PRODUCT_SVC_FAIL
		}
		
		if _, err := svc.shop.AddProductService(service); err != nil {
			log.Log.Errorf("shop_trace: add product service fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.SHOP_ADD_PRODUCT_SVC_FAIL
		}
	}
	
	category := svc.shop.GetProductCategoryInfoById(fmt.Sprint(params.CategoryId))
	if category == nil {
		log.Log.Errorf("shop_trace: get category fail, categoryId:%d", params.CategoryId)
		svc.engine.Rollback()
		return errdef.SHOP_GET_CATEGORY_FAIL
	}
	
	// 是否重置sku true为重置
	if params.IsReset == true {
		// 软删除废弃的sku列表
		if _, err := svc.shop.DelSkuByProductId(fmt.Sprint(params.Id)); err != nil {
			svc.engine.Rollback()
			return errdef.ERROR
		}
		
		// 软删除关联的分类
		if _, err := svc.shop.DelProductCategoryRelated(fmt.Sprint(params.Id)); err != nil {
			svc.engine.Rollback()
			return errdef.ERROR
		}
		
		// 重新添加新数据
		code := svc.AddSkuListInfo(params, params.Id, category.CategoryName)
		if code != errdef.SUCCESS {
			svc.engine.Rollback()
			return code
		}
		
		svc.engine.Commit()
		return code
	}
	
	code := svc.UpdateSkuListInfo(params, now, category.CategoryName)
	if code != errdef.SUCCESS {
		svc.engine.Rollback()
		return code
	}
	
	svc.engine.Commit()
	return errdef.SUCCESS
}

func (svc *ShopModule) UpdateSkuListInfo(params *mshop.AddOrEditProductReq, now int, categoryName string) int {
	for _, item := range params.SkuList {
		images, _ := util.JsonFast.MarshalToString(item.Images)
		ownSpec, _ := util.JsonFast.MarshalToString(item.OwnSpec)
		
		sku := &models.ProductSku{
			Id: item.Id,
			ProductId: int64(params.Id),
			Title: item.Title,
			SkuImage: string(item.SkuImage),
			SkuNo: item.SkuNo,
			Images: images,
			CurPrice: item.CurPrice,
			MarketPrice: item.MarketPrice,
			IsFreeShip: item.IsFreeShip,
			Indexes: item.Indexes,
			OwnSpec: ownSpec,
			Status: item.Status,
			UpdateAt: now,
			Sortorder: item.Sortorder,
		}
		
		
		if _, err := svc.shop.UpdateProductSku(sku); err != nil {
			log.Log.Errorf("shop_trace: update product sku fail, err:%s", err)
			return errdef.SHOP_UPDATE_SKU_FAIL
		}
		
		stock := &models.ProductSkuStock{
			SkuId: item.Id,
			ProductId: params.Id,
			MinBuy: 1,
			MaxBuy: 200,
			UpdateAt: now,
			Stock: item.Stock,
		}
		
		if _, err := svc.shop.UpdateProductSkuStockInfo(stock); err != nil {
			log.Log.Errorf("shop_trace: update product sku stock info fail, err:%s", err)
			return errdef.SHOP_UPDATE_SKU_STOCK_FAIL
		}
	}
	
	
	condition := fmt.Sprintf("product_id=%d AND category_id=%d", params.Id, params.CategoryId)
	relatedInfo, err := svc.shop.GetProductCategoryRelated(condition)
	if err != nil {
		log.Log.Errorf("shop_trace: get product category related fail, err:%s", err)
	}
	
	if relatedInfo == nil {
		related := &models.ProductCategoryRelated{
			ProductId: params.Id,
			CategoryId: params.CategoryId,
			CreateAt: now,
			CategoryName: categoryName,
		}
		
		if _, err := svc.shop.AddProductCategoryRelated(related); err != nil {
			log.Log.Errorf("shop_trace: add related fail, err:%s", err)
			return errdef.SHOP_ADD_RELATED_FAIL
		}
	}
	
	return errdef.SUCCESS
}

// 订单列表
func (svc *ShopModule) OrderList(reqType, keyword string, page, size int) (int, int64, []*mshop.OrderResp) {
	condition := svc.GetQueryCondition(reqType, keyword)
	if condition == "" {
		return errdef.INVALID_PARAMS, 0, nil
	}
	
	offset := (page - 1) * size
	list, err := svc.shop.GetOrderList(condition, offset, size)
	if err != nil {
		log.Log.Errorf("shop_trace: get order list fail, err:%s", err)
		return errdef.SHOP_ORDER_LIST_FAIL, 0, nil
	}
	
	if len(list) == 0 {
		return errdef.SUCCESS, 0, []*mshop.OrderResp{}
	}
	
	count, err := svc.shop.GetOrderTotal(condition)
	if err != nil {
		log.Log.Errorf("shop_trace: get order total fail, err:%s", err)
	}
	
	res := make([]*mshop.OrderResp, len(list))
	for index, item := range list {
		info, err := svc.OrderInfo(item)
		if err != nil {
			log.Log.Errorf("shop_trace: get order info fail, err:%s", err)
			return errdef.SHOP_ORDER_LIST_FAIL, 0, nil
		}
		
		res[index] = info
	}
	
	
	return errdef.SUCCESS, count, res
}

// 订单信息
func (svc *ShopModule) OrderInfo(item models.Orders) (*mshop.OrderResp, error) {
	info := &mshop.OrderResp{}
	if err := util.JsonFast.UnmarshalFromString(item.Extra, &info); err != nil {
		log.Log.Errorf("shop_trace: unmarshal fail, orderId:%s, err:%s", item.OrderId, err)
		return nil, err
	}
	
	info.UserId = item.UserId
	user := svc.user.FindUserByUserid(item.UserId)
	if user != nil {
		info.MobileNum = util.HideMobileNum(fmt.Sprint(user.MobileNum))
		info.RealMobileNum = fmt.Sprint(user.MobileNum)
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
	info.DeliveryTime = item.DeliveryTime
	
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

func (svc *ShopModule) GetQueryCondition(reqType, keyword string) string {
	var condition string
	switch reqType {
	// 查看全部订单
	case "1":
		condition = "1=1"
	// 待付款订单
	case "2":
		condition = "pay_status=0"
	// 待收货订单
	case "3":
		condition = "pay_status=2 AND delivery_status in(0, 1)"
	// 已完成订单
	case "4":
		condition = "pay_status=2 AND delivery_status=2"
	default:
		return "1=1"
	}
	
	if keyword != "" {
		condition += " AND order_id like '%" + keyword + "%'"
	}
	
	return condition
}

// 订单确认收货
func (svc *ShopModule) ConfirmReceipt(param *mshop.ChangeOrderReq) int {
	order, err := svc.shop.GetOrder(param.OrderId)
	if err != nil {
		log.Log.Errorf("shop_trace: order not exists, orderId:%s", param.OrderId)
		return errdef.SHOP_ORDER_NOT_EXISTS
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

func (svc *ShopModule) DeliverProduct(param *mshop.DeliverProductReq) int {
	order, err := svc.shop.GetOrder(param.OrderId)
	if err != nil {
		log.Log.Errorf("shop_trace: order not exists, orderId:%s", param.OrderId)
		return errdef.SHOP_ORDER_NOT_EXISTS
	}
	
	// 订单 != 支付成功 || 配送状态 != 未配送
	if order.PayStatus != consts.SHOP_ORDER_TYPE_PAID || order.DeliveryStatus != consts.NOT_DELIVERED {
		log.Log.Errorf("shop_trace: not allow deliver, orderId:%s, payStatus:%d, deliveryStatus:%d", order.OrderId,
			order.PayStatus, order.DeliveryStatus)
		return errdef.SHOP_NOT_ALLOW_DELIVER
	}
	
	now := int(time.Now().Unix())
	// 支付状态=已支付 && 配送状态=未配送
	condition := fmt.Sprintf("order_id='%s' AND pay_status=%d AND delivery_status=0", order.OrderId, consts.SHOP_ORDER_TYPE_PAID)
	cols := "update_at, delivery_time, delivery_code, delivery_type_name, delivery_telephone, delivery_status"
	order.UpdateAt = now
	order.DeliveryTime = now
	order.DeliveryTelephone = param.DeliveryTelephone
	order.DeliveryTypeName = param.DeliveryTypeName
	order.DeliveryCode = param.DeliveryCode
	// 配送状态 修改为 已配送
	order.DeliveryStatus = consts.HAS_DELIVERED
	if _, err := svc.shop.UpdateOrderInfo(condition, cols, order); err != nil {
		log.Log.Errorf("shop_trace: update order info fail, orderId:%s, err:%s", order.OrderId, err)
		return errdef.SHOP_CONFIRM_RECEIPT_FAIL
	}
	
	return errdef.SUCCESS
}

func (svc *ShopModule) OrderCallback(orderId string) int {
	order, err := svc.shop.GetOrder(orderId)
	if err != nil || order.PayStatus != consts.SHOP_ORDER_TYPE_WAIT {
		return errdef.SHOP_ORDER_NOT_EXISTS
	}
	
	ok, err := svc.pay.GetPaymentChannel(order.PayChannelId)
	if !ok || err != nil {
		log.Log.Errorf("order_trace: get payment channel fail, orderId:%s, ok:%v, err:%s", order.OrderId,
			ok, err)
		return errdef.ERROR
	}
	
	switch svc.pay.PayChannel.Identifier {
	case "alipay":
		alipayCli := alipay.NewAliPay(true, svc.pay.PayChannel.AppId, svc.pay.PayChannel.PrivateKey)
		alipayCli.OutTradeNo = order.OrderId
		rsp, err := alipayCli.TradeQuery()
		if err != nil {
			log.Log.Errorf("order_trace: alipay trade query fail, orderId:%s, err:%s", order.OrderId, err)
			return errdef.ERROR
		}
		
		amount, err := strconv.ParseFloat(strings.Trim(rsp.Response.TotalAmount, " "), 64)
		if err != nil {
			log.Log.Errorf("order_trace: parse float fail, err:%s", err)
			return errdef.ERROR
		}
		
		if int(amount * 100) != order.PayAmount {
			log.Log.Error("order_trace: amount not match, orderAmount:%d, amount:%d",
				order.PayAmount, amount * 100)
			return errdef.ERROR
		}
		
		if rsp.Response.Code != "10000" || rsp.Response.TradeStatus != consts.TradeSuccess {
			log.Log.Errorf("order_trace: request fail, orderId:%s, response:%+v", order.OrderId, rsp.Response)
			return errdef.ERROR
		}
		
		order.Transaction = rsp.Response.TradeNo
		payTime, _ := time.ParseInLocation("2006-01-02 15:04:05", rsp.Response.SendPayDate, time.Local)
		order.PayTime = int(payTime.Unix())
		
	case "weixin":
		wxCli := wechat.NewWechatPay(true, svc.pay.PayChannel.AppId, svc.pay.PayChannel.AppKey, svc.pay.PayChannel.AppSecret)
		wxCli.OutTradeNo = order.OrderId
		rsp, _, err := wxCli.TradeQuery()
		if err != nil {
			log.Log.Errorf("order_trace: wx trade query fail, orderId:%s, err:%s", order.OrderId, err)
			return errdef.ERROR
		}
		
		totalFee, err := strconv.Atoi(rsp.TotalFee)
		if err != nil {
			return errdef.ERROR
		}
		
		if order.PayAmount != totalFee {
			log.Log.Error("order_trace: amount not match, orderAmount:%d, amount:%d",
				order.PayAmount, rsp.TotalFee)
			return errdef.ERROR
		}
		
		if rsp.ReturnCode != "SUCCESS" || rsp.ResultCode != "SUCCESS" || rsp.TradeState != "SUCCESS" {
			log.Log.Errorf("order_trace: request fail, orderId:%s, response:%+v", order.OrderId, rsp)
			return errdef.ERROR
		}
		
		order.Transaction = rsp.TransactionId
		payTime, _ := time.ParseInLocation("20060102150405", rsp.TimeEnd, time.Local)
		order.PayTime = int(payTime.Unix())
		
	default:
		return errdef.ERROR
	}
	
	return svc.SetOrderSuccess(order)
}

func (svc *ShopModule) SetOrderSuccess(order *models.Orders) int {
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR
	}
	
	condition := fmt.Sprintf("`pay_status`=%d AND `order_id`='%s'", consts.SHOP_ORDER_TYPE_WAIT, order.OrderId)
	cols := "pay_status, transaction, pay_time, update_at"
	order.PayStatus = consts.SHOP_ORDER_TYPE_PAID
	now := int(time.Now().Unix())
	order.UpdateAt = now
	// 更新订单状态
	if _, err := svc.shop.UpdateOrderInfo(condition, cols, order); err != nil {
		log.Log.Errorf("shop_trace: update order info fail, orderId:%s, err:%v", order.OrderId, err)
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	list, err := svc.shop.GetOrderProductList(order.OrderId)
	if err != nil {
		log.Log.Errorf("shop_trace: get order product list fail, orderId:%s, err:%s", order.OrderId, err)
		svc.engine.Rollback()
		return errdef.ERROR
	}
	
	for _, item := range list {
		product, err := svc.shop.GetProductSpu(fmt.Sprint(item.ProductId))
		if err != nil || product == nil {
			log.Log.Errorf("order_trace: get product spu fail, productId:%d,  err:%s", item.ProductId, err)
			continue
		}
		
		condition := fmt.Sprintf("id=%d", item.ProductId)
		cols := "sale_num"
		// 更新产品销量
		if _, err := svc.shop.UpdateProductInfo(condition, cols, product); err != nil {
			log.Log.Errorf("order_trace: update product info fail, productId:%d, err:%s", product.Id, err)
			svc.engine.Rollback()
			return errdef.ERROR
		}
	}
	
	svc.engine.Commit()
	return errdef.SUCCESS
}
