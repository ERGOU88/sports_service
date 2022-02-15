package shop

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/backend/controller/cshop"
	"sports_service/server/global/backend/log"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models"
	"sports_service/server/models/mshop"
	"sports_service/server/util"
)

func ProductList(c *gin.Context) {
	reply := errdef.New(c)
	
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	sortType := c.DefaultQuery("sort_type", "0")
	keyword := c.Query("keyword")
	
	svc := cshop.New(c)
    code, list := svc.GetProductList(sortType, keyword, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func ProductCategory(c *gin.Context) {
	reply := errdef.New(c)
	svc := cshop.New(c)
	reply.Data["list"] = svc.GetProductCategoryConf()
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

func AddCategory(c *gin.Context) {
	reply := errdef.New(c)
	params := &models.ProductCategory{}
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("shop_trace: invalid param, err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.AddCategory(params))
}

func EditCategory(c *gin.Context) {
	reply := errdef.New(c)
	params := &models.ProductCategory{}
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("shop_trace: invalid param, err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.EditCategory(params))
}

func ServiceList(c *gin.Context) {
	reply := errdef.New(c)
	svc := cshop.New(c)
	code, list := svc.GetServiceList()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func AddService(c *gin.Context) {
	reply := errdef.New(c)
	params := &models.ShopServiceConf{}
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("shop_trace: invalid param, err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.AddService(params))
}

func EditService(c *gin.Context) {
	reply := errdef.New(c)
	params := &models.ShopServiceConf{}
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("shop_trace: invalid param, err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.UpdateService(params))
}

func AddSpecification(c *gin.Context) {
	reply := errdef.New(c)
	params := &mshop.AddOrEditCategorySpecReq{}
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("shop_trace: invalid param, err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.AddCategorySpec(params))
}

func EditSpecification(c *gin.Context) {
	reply := errdef.New(c)
	params := &mshop.AddOrEditCategorySpecReq{}
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("shop_trace: invalid param, err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.EditCategorySpec(params))
}

func DelSpecification(c *gin.Context) {
	reply := errdef.New(c)
	categoryId := c.Query("category_id")
	if categoryId == "" {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.DelCategorySpec(categoryId))
}
