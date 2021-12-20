package cshop

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mshop"
	"sports_service/server/models/muser"
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

func (svc *ShopModule) GetProductDetail(productId string) {

}
