package shop

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cshop"
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

	svc := cshop.New(c)
	svc.GetProductDetail(productId)
	reply.Response(http.StatusOK, errdef.SUCCESS)
}
