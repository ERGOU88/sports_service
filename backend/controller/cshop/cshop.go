package cshop

import (
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
	"fmt"
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
