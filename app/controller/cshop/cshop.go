package cshop

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mshop"
	"sports_service/server/models/muser"
	"time"
)

type ShopModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	shop        *mshop.ShopModel
	//resp        *mshop.OrderResp
}

func New(c *gin.Context) ShopModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return ShopModule{
		context: c,
		user: muser.NewUserModel(socket),
		shop: mshop.NewShop(socket),
		//resp: &mshop.OrderResp{Products: make([]*mshop.Product, 0)},
		engine: socket,
	}
}

func (svc *ShopModule) GetProducts(categoryId, sortType string, page, size int) (int, []mshop.ProductSimpleInfo) {
	if categoryId == "0" {
		return svc.GetProductList(sortType, "", page, size)
	}

	return svc.GetProductListByCategory(categoryId, sortType,  page, size)
}

func (svc *ShopModule) GetProductList(sortType, keyword string, page, size int) (int, []mshop.ProductSimpleInfo) {
	offset := (page - 1) * size
	list, err := svc.shop.GetAllSpu(sortType, keyword, offset, size)
	if err != nil {
		log.Log.Errorf("shop_trace: get all spu fail, err:%s", err)
		return errdef.SHOP_GET_ALL_SPU_FAIL, nil
	}

	return errdef.SUCCESS, list
}

func (svc *ShopModule) GetProductListByCategory(categoryId, sortType string, page, size int) (int, []mshop.ProductSimpleInfo) {
	offset := (page - 1) * size
	list, err := svc.shop.GetSpuListByCategory(categoryId, sortType, offset, size)
	if err != nil {
		log.Log.Errorf("shop_trace: get spu list by category fail, err:%s", err)
		return errdef.SHOP_GET_SPU_BY_CATEGORY_FAIL, nil
	}

	return errdef.SUCCESS, list
}

func (svc *ShopModule) GetProductCategoryConf() []*mshop.Category {
	conf := svc.shop.GetProductCategory()
	if conf == nil {
		return []*mshop.Category{}
	}

	info := &mshop.Category{
	   CategoryId: 0,
	   CategoryName: "综合",
	}

	conf = append([]*mshop.Category{info}, conf...)
	return conf
}

// 推荐的商品
func (svc *ShopModule) RecommendProduct(productId string) (int, []mshop.ProductSimpleInfo) {
	list, err := svc.shop.GetRecommendProducts(productId, consts.RECOMMEND_DEFAULT_LIMIT)
	if err != nil {
		log.Log.Errorf("shop_trace: get recommend products fail, err:%s", err)
		return errdef.SHOP_RECOMMEND_FAIL, nil
	}

	return errdef.SUCCESS, list
}

// 商品详情 获取商品实体信息
func (svc *ShopModule) GetProductDetail(productId, indexes, userId string) (int, *mshop.ProductDetailInfo) {
	if _, err := svc.shop.GetProductSpu(productId); err != nil {
		log.Log.Errorf("shop_trace: get product spu fail, productId:%s, err:%s", productId, err)
		return errdef.SHOP_PRODUCT_SPU_FAIL, nil
	}

	detail, err := svc.shop.GetProductDetail(productId, indexes)
	if err != nil {
		log.Log.Errorf("shop_trace: get product detail fail, productId:%s, err:%s", productId, err)
		return errdef.SHOP_PRODUCT_SKU_FAIL, nil
	}

	now := time.Now().Unix()
	if now >= detail.StartTime && now < detail.EndTime {
		detail.HasActivities = 1
		detail.RemainDuration = detail.EndTime - now
	}

	stockInfo, err := svc.shop.GetProductSkuStock(fmt.Sprint(detail.Id))
	if err == nil {
		detail.Stock = stockInfo.Stock
		detail.MinBuy = stockInfo.MinBuy
		detail.MaxBuy = stockInfo.MaxBuy
	}
	
	if detail.Specifications == nil {
		detail.Specifications = make([]mshop.SpecInfo, 0)
	}
	
	if detail.SpecTemplate == nil {
		detail.SpecTemplate = make([]mshop.SpecTemplate, 0)
	}
	
	if detail.OwnSpec == nil {
		detail.OwnSpec = make([]mshop.OwnSpec, 0)
	}
	
	if userId != "" {
		count, err := svc.shop.GetProductCartNum(userId)
		if err != nil {
			log.Log.Errorf("shop_trace: get product cart num fail, err:%s", err)
		}
		
		detail.ProductCartNum = count
	}

	return errdef.SUCCESS, detail
}

// 获取地区信息
func (svc *ShopModule) GetAreaConf() []*mshop.AreaInfo {
	conf := svc.shop.GetArea()
	if conf == nil {
		return []*mshop.AreaInfo{}
	}

	return conf
}

// 用户添加/更新 地址信息
func (svc *ShopModule) AddOrUpdateUserAddr(info *models.UserAddress) (int, int) {
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR, 0
	}
	// 校验手机号合法性
	if b := svc.user.CheckCellPhoneNumber(info.Mobile); !b {
		log.Log.Errorf("shop_trace: invalid mobile num %v", info.Mobile)
		svc.engine.Rollback()
		return errdef.USER_INVALID_MOBILE_NUM, 0
	}

	// 如果要将当前地址设置为默认
	if info.IsDefault == 1 {
		addrList, err := svc.shop.GetUserAddrByUserId(info.UserId, 0, 1)
		if err != nil {
			log.Log.Errorf("shop_trace: get user addr by userId fail, userId:%d", info.UserId)
			svc.engine.Rollback()
			return errdef.SHOP_USER_ADDR_NOT_FOUND, 0
		}

		if len(addrList) > 0 {
			// 先将其他地址置为不默认
			if _, err := svc.shop.UpdateUserDefaultAddr("", info.UserId, 0); err != nil {
				log.Log.Errorf("shop_trace: update user default addr fail, err:%s", err)
				svc.engine.Rollback()
				return errdef.SHOP_UPDATE_USER_ADDR_FAIL, 0
			}
		}
	}

	now := int(time.Now().Unix())
	info.UpdateAt = now
	if info.Id <= 0 {
		info.CreateAt = now
		// 添加地址
		if _, err := svc.shop.AddUserAddr(info); err != nil {
			log.Log.Errorf("shop_trace: add user addr fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.SHOP_ADD_USER_ADDR_FAIL, 0
		}

		svc.engine.Commit()
		return errdef.SUCCESS, info.Id
	}

	addr, err := svc.shop.GetUserAddr(fmt.Sprint(info.Id), info.UserId)
	if err != nil || addr == nil {
		log.Log.Errorf("shop_trace: user addr not found, id:%d", info.Id)
		svc.engine.Rollback()
		return errdef.SHOP_USER_ADDR_NOT_FOUND, 0
	}

	// 修改地址
	if _, err := svc.shop.UpdateUserAddr(info); err != nil {
		svc.engine.Rollback()
		return errdef.SHOP_UPDATE_USER_ADDR_FAIL, 0
	}

	svc.engine.Commit()
	return errdef.SUCCESS, info.Id
}

func (svc *ShopModule) GetUserAddrList(page, size int, userId string) (int, []*models.UserAddress) {
	offset := (page - 1) * size
	addrList, err := svc.shop.GetUserAddrByUserId(userId, offset, size)
	if err != nil {
		log.Log.Errorf("shop_trace: get user addr by userId fail, userId:%s, err:%s", userId, err)
		return errdef.SHOP_GET_USER_ADDR_FAIL, nil
	}
	
	//for _, item := range addrList {
	//	item.Mobile = util.HideMobileNum(item.Mobile)
	//}

	return errdef.SUCCESS, addrList
}

// 添加商品购物车
func (svc *ShopModule) AddProductCart(info *models.ProductCart) (int, int64) {
	if info.Count <= 0 {
		return errdef.INVALID_PARAMS, 0
	}

	if _, err := svc.shop.GetProductSpu(fmt.Sprint(info.ProductId)); err != nil {
		log.Log.Errorf("shop_trace: get product spu fail, productId:%d, err:%s", info.ProductId, err)
		return errdef.SHOP_PRODUCT_SPU_FAIL, 0
	}

	condition := fmt.Sprintf("id=%d AND product_id=%d", info.SkuId, info.ProductId)
	if _, err := svc.shop.GetProductSku(condition); err != nil {
		log.Log.Errorf("shop_trace: get product sku by id fail, skuId:%d, err:%s", info.SkuId, err)
		return errdef.SHOP_PRODUCT_SKU_FAIL, 0
	}

	info.CreateAt = int(time.Now().Unix())
	condition = fmt.Sprintf("product_id=%d AND user_id=%s AND sku_id=%d", info.ProductId, info.UserId, info.SkuId)
	cartInfo, err := svc.shop.GetProductCart(condition)
	if err != nil {
		log.Log.Errorf("shop_trace: get product cart failm err:%s", err)
		return errdef.SHOP_GET_PRODUCT_CART_FAIL, 0
	}
	
	total, err := svc.shop.GetProductCartNum(info.UserId)
	if err != nil {
		log.Log.Errorf("shop_trace: get product cart num fail, err:%s", err)
	}
	
	if cartInfo == nil {
		if _, err := svc.shop.AddProductCart(info); err != nil {
			log.Log.Errorf("shop_trace: add product cart fail, err:%s", err)
			return errdef.SHOP_ADD_PRODUCT_CART_FAIL, total
		}
		
		total += 1

	} else {
		affected, err := svc.shop.UpdateProductCart(info)
		if affected <= 0 || err != nil {
			log.Log.Errorf("shop_trace: update product cart fail, err:%s", err)
			return errdef.SHOP_ADD_PRODUCT_CART_FAIL, total
		}
	}


	return errdef.SUCCESS, total
}

// 获取商品购物车列表
func (svc *ShopModule) GetProductCartList(userId string) (int, []*mshop.ProductCartInfo) {
	//offset := (page - 1) * size
	list, err := svc.shop.GetProductCartList(userId)
	if err != nil {
		log.Log.Errorf("shop_trace: get product cart list fail, err:%s", err)
		return errdef.SHOP_GET_PRODUCT_CART_FAIL, nil
	}
	
	if len(list) == 0 {
		return errdef.SUCCESS, []*mshop.ProductCartInfo{}
	}

	for _, item := range list {
		stockInfo, err := svc.shop.GetProductSkuStock(fmt.Sprint(item.Id))
		if err != nil {
			log.Log.Errorf("shop_trace: get product sku stock fail, skuId:%d, err:%s", item.Id, err)
		}
		
		log.Log.Infof("stockInfo:%+v", stockInfo)
		
		if stockInfo != nil {
			item.Stock = stockInfo.Stock
			item.MinBuy = stockInfo.MinBuy
			item.MinBuy = stockInfo.MaxBuy
		}
		
		now := time.Now().Unix()
		if now >= item.StartTime && now < item.EndTime {
			item.HasActivities = 1
			item.RemainDuration = item.EndTime - now
		}
	}

	return errdef.SUCCESS, list
}

func (svc *ShopModule) SearchProduct(sortType, keyword string, page, size int) (int, []mshop.ProductSimpleInfo) {
	if keyword == "" {
		return errdef.SUCCESS, []mshop.ProductSimpleInfo{}
	}

	code, list := svc.GetProductList(sortType, keyword, page, size)
	if list == nil {
		return code, []mshop.ProductSimpleInfo{}
	}
	
	return code, list
}

// 更新用户商品购物车
func (svc *ShopModule) UpdateProductCart(list []*models.ProductCart) (int, []int) {
	if len(list) == 0 {
		log.Log.Errorf("shop_trace: invalid params, len:%d", len(list))
		return errdef.INVALID_PARAMS, []int{}
	}

	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR, []int{}
	}

	ids := make([]int, 0)
	for _, item := range list {
		if item.Count <= 0 {
			svc.engine.Rollback()
			return errdef.INVALID_PARAMS, []int{}
		}

		if _, err := svc.shop.GetProductSpu(fmt.Sprint(item.ProductId)); err != nil {
			log.Log.Errorf("shop_trace: get product spu fail, productId:%d, err:%s", item.ProductId, err)
			svc.engine.Rollback()
			return errdef.SHOP_PRODUCT_SPU_FAIL, []int{}
		}

		condition := fmt.Sprintf("id=%d AND product_id=%d", item.SkuId, item.ProductId)
		if _, err := svc.shop.GetProductSku(condition); err != nil {
			log.Log.Errorf("shop_trace: get product sku by id fail, skuId:%d, err:%s", item.SkuId, err)
			svc.engine.Rollback()
			return errdef.SHOP_PRODUCT_SKU_FAIL, []int{}
		}

		stockInfo, err := svc.shop.GetProductSkuStock(fmt.Sprint(item.SkuId))
		if err != nil {
			log.Log.Errorf("shop_trace: get product sku stock fail, skuId:%d, err:%s", item.SkuId, err)
			svc.engine.Rollback()
			return errdef.ERROR, []int{}
		}

		if item.Count > stockInfo.MaxBuy && stockInfo.MaxBuy != 0 || item.Count < stockInfo.MinBuy {
			log.Log.Errorf("shop_trace: invalid count, item.Count:%d, max:%d, min:%d", item.Count, stockInfo.MaxBuy, stockInfo.MinBuy)
			ids = append(ids, item.Id)
		}

		if len(ids) == 0 {
			if _, err := svc.shop.UpdateProductCartById(item); err != nil {
				log.Log.Errorf("shop_trace: update product cart by id fail, id:%d, err:%s", item.Id, err)
				svc.engine.Rollback()
				return errdef.SHOP_UPDATE_PRODUCT_CART_FAIL, []int{}
			}
		}
	}

	if len(ids) > 0 {
		log.Log.Errorf("shop_trace: invalid count, ids:%+v", ids)
		svc.engine.Rollback()
		return errdef.INVALID_PARAMS, ids
	}

	svc.engine.Commit()

	return errdef.SUCCESS, []int{}
}

func (svc *ShopModule) DeleteProductCart(ids []int, userId string) int {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		return errdef.USER_NOT_EXISTS
	}
	
	affected, err := svc.CleanProductCart(ids, userId)
	if err != nil || int(affected) != len(ids) {
		return errdef.SHOP_DEL_PRODUCT_CART_FAIL
	}
	
	return errdef.SUCCESS
}
