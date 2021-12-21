package cshop

import (
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

func (svc *ShopModule) GetProducts(categoryId, sortType string, page, size int) (int, []mshop.ProductSimpleInfo) {
	if categoryId == "0" {
		return svc.GetProductList(sortType, page, size)
	}

	return svc.GetProductListByCategory(categoryId, sortType,  page, size)
}

func (svc *ShopModule) GetProductList(sortType string, page, size int) (int, []mshop.ProductSimpleInfo) {
	offset := (page - 1) * size
	list, err := svc.shop.GetAllSpu(sortType, offset, size)
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
func (svc *ShopModule) GetProductDetail(productId, indexes string) (int, *mshop.ProductDetailInfo) {
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
func (svc *ShopModule) AddOrUpdateUserArea(info *models.UserAddress) int {
	if err := svc.engine.Begin(); err != nil {
		return errdef.ERROR
	}
	// 校验手机号合法性
	if b := svc.user.CheckCellPhoneNumber(info.Mobile); !b {
		log.Log.Errorf("shop_trace: invalid mobile num %v", info.Mobile)
		svc.engine.Rollback()
		return errdef.USER_INVALID_MOBILE_NUM
	}

	if info.Id <= 0 {
		// 如果要将当前地址设置为默认
		if info.IsDefault == 1 {
			// 先将其他地址置为不默认
			if _, err := svc.shop.UpdateUserDefaultAddr("", info.UserId, 0); err != nil {
				log.Log.Errorf("shop_trace: update user default addr fail, err:%s", err)
				svc.engine.Rollback()
				return errdef.SHOP_ADD_USER_ADDR_FAIL
			}
		}


		if _, err := svc.shop.AddUserAddr(info); err != nil {
			log.Log.Errorf("shop_trace: add user addr fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.SHOP_ADD_USER_ADDR_FAIL
		}

		svc.engine.Commit()
		return errdef.SUCCESS
	}

	addr, err := svc.shop.GetUserAddrById(fmt.Sprint(info.Id))
	if err != nil || addr == nil {
		log.Log.Errorf("shop_trace: user addr not found, id:%d", info.Id)
		svc.engine.Rollback()
		return errdef.SHOP_USER_ADDR_NOT_FOUND
	}

	if info.IsDefault == 1 {
		// 先将其他地址置为不默认
		if _, err := svc.shop.UpdateUserDefaultAddr("", info.UserId, 0); err != nil {
			log.Log.Errorf("shop_trace: update user default addr fail, err:%s", err)
			svc.engine.Rollback()
			return errdef.SHOP_UPDATE_USER_ADDR_FAIL
		}
	}

	if _, err := svc.shop.UpdateUserAddr(info); err != nil {
		svc.engine.Rollback()
		return errdef.SHOP_UPDATE_USER_ADDR_FAIL
	}

	svc.engine.Commit()
	return errdef.SUCCESS
}
