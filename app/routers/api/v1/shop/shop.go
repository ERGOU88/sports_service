package shop

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cshop"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"sports_service/server/util"
	"sports_service/server/global/app/errdef"
)

func ProductList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	categoryId := c.DefaultQuery("category_id", "0")
	sortType := c.DefaultQuery("sort_type", "0")
	svc := cshop.New(c)
	code, list := svc.GetProducts(categoryId, sortType, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func ProductCategory(c *gin.Context) {
	reply := errdef.New(c)
	svc := cshop.New(c)
	reply.Data["list"] = svc.GetProductCategoryConf()
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

func RecommendProduct(c *gin.Context) {
	reply := errdef.New(c)
	curProductId := c.DefaultQuery("product_id", "0")

	svc := cshop.New(c)
	code, list := svc.RecommendProduct(curProductId)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func ProductDetail(c *gin.Context) {
	reply := errdef.New(c)
	productId := c.Query("product_id")
	indexes := c.Query("indexes")

	svc := cshop.New(c)
	code, detail := svc.GetProductDetail(productId, indexes)
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, code)
}

func AreaConfig(c *gin.Context) {
	reply := errdef.New(c)
	svc := cshop.New(c)
	areaList := svc.GetAreaConf()
	reply.Data["list"] = areaList
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

func EditArea(c *gin.Context) {
	reply := errdef.New(c)
	params := &models.UserAddress{}
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("shop_trace: invalid param, params:%+v, err:%s", params, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.AddOrUpdateUserArea(params))
}
