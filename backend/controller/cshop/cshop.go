package cshop

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/backend/log"
	"sports_service/server/models"
	"sports_service/server/models/mshop"
	"sports_service/server/models/muser"
	"sports_service/server/util"
	"time"
)

type ShopModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	shop        *mshop.ShopModel
}

func New(c *gin.Context) ShopModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return ShopModule{
		context: c,
		user: muser.NewUserModel(socket),
		shop: mshop.NewShop(socket),
		engine: socket,
	}
}

func (svc *ShopModule) GetProductList(sortType, keyword string, page, size int) (int, []*mshop.ProductSimpleInfo) {
	offset := (page - 1) * size
	list, err := svc.shop.GetAllSpu(sortType, keyword, offset, size)
	if err != nil {
		log.Log.Errorf("shop_trace: get all spu fail, err:%s", err)
		return errdef.SHOP_GET_ALL_SPU_FAIL, nil
	}
	
	if len(list) == 0 {
		return errdef.SUCCESS, []*mshop.ProductSimpleInfo{}
	}
	
	for _, item := range list {
		skuList, err := svc.shop.GetProductSkuList(fmt.Sprint(item.Id))
		if err != nil {
			log.Log.Errorf("shop_trace: get product sku by spuId fail, spuId:%d, err:%s", item.Id, err)
			continue
		}
		
		item.SkuNum = len(skuList)
	}
	
	return errdef.SUCCESS, list
}

func (svc *ShopModule) GetProductCategoryConf() []*mshop.Category {
	conf := svc.shop.GetProductCategory()
	if conf == nil {
		return []*mshop.Category{}
	}
	
	return conf
}

func (svc *ShopModule) AddCategory(params *models.ProductCategory) int {
	if _, err := svc.shop.AddProductCategory(params); err != nil {
		log.Log.Errorf("shop_trace: add category fail, err:%s", err)
		return errdef.SHOP_ADD_CATEGORY_FAIL
	}
	
	return errdef.SUCCESS
}

func (svc *ShopModule) EditCategory(params *models.ProductCategory) int {
	if _, err := svc.shop.UpdateProductCategory(params); err != nil {
		log.Log.Errorf("shop_trace: update category fail, err:%s", err)
		return errdef.SHOP_EDIT_CATEGORY_FAIL
	}
	
	return errdef.SUCCESS
}

func (svc *ShopModule) GetServiceList() (int, []models.ShopServiceConf) {
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
	service, _ := util.JsonFast.MarshalToString(params.AfterService)
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
		AfterService: service,
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
			SkuImage: item.SkuImage,
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
	service, _ := util.JsonFast.MarshalToString(params.AfterService)
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
		AfterService: service,
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
			ProductId: int64(params.Id),
			Title: item.Title,
			SkuImage: item.SkuImage,
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
	
	
	relatedInfo, err := svc.shop.GetProductCategoryRelated(fmt.Sprint(params.Id), fmt.Sprint(params.CategoryId))
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
