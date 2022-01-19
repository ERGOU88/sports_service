package shop

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cshop"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mshop"
	"sports_service/server/util"
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
	userId := c.Query("user_id")

	svc := cshop.New(c)
	code, detail := svc.GetProductDetail(productId, indexes, userId)
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

func EditAddr(c *gin.Context) {
	reply := errdef.New(c)
	params := &models.UserAddress{}
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("shop_trace: invalid param, params:%+v, err:%s", params, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	userId, _ := c.Get(consts.USER_ID)
	params.UserId = userId.(string)
	svc := cshop.New(c)
	code, id := svc.AddOrUpdateUserAddr(params)
	reply.Data["id"] = id
	reply.Response(http.StatusOK, code)
}

func UserAddrList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	userId, _ := c.Get(consts.USER_ID)
	svc := cshop.New(c)
	code, list := svc.GetUserAddrList(page, size, userId.(string))
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func AddProductCart(c *gin.Context) {
	reply := errdef.New(c)
	params := &models.ProductCart{}
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("shop_trace: invalid param, params:%+v, err:%s", params, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	userId, _ := c.Get(consts.USER_ID)
	params.UserId = userId.(string)
	svc := cshop.New(c)
	code, total := svc.AddProductCart(params)
	reply.Data["total"] = total
	reply.Response(http.StatusOK, code)
}

func ProductCart(c *gin.Context) {
	reply := errdef.New(c)
	//page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	userId, _ := c.Get(consts.USER_ID)

	svc := cshop.New(c)
	code, list := svc.GetProductCartList(userId.(string))
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func SearchProduct(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	sortType := c.DefaultQuery("sort_type", "0")
	keyword := c.Query("keyword")
	svc := cshop.New(c)
	code, list := svc.SearchProduct(sortType, keyword, page, size)
	reply.Data["list"] = list
	if len(list) == 0 {
		_, recommend := svc.RecommendProduct("")
		reply.Data["recommend"] = recommend
	}

	reply.Response(http.StatusOK, code)
}

func UpdateProductCart(c *gin.Context) {
	reply := errdef.New(c)
	param := &mshop.UpdateProductCartParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("shop_trace: invalid param, params:%+v, err:%s", param, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cshop.New(c)
	code, id := svc.UpdateProductCart(param.Params)
	reply.Data["id"] = id
	reply.Response(http.StatusOK, code)
}

func DeleteProductCart(c *gin.Context) {
	reply := errdef.New(c)
	param := &mshop.DeleteProductCartParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("shop_trace: invalid param, params:%+v, err:%s", param, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	
	userId, _ := c.Get(consts.USER_ID)
	svc := cshop.New(c)
	code := svc.DeleteProductCart(param.Ids, userId.(string))
	reply.Response(http.StatusOK, code)
}

func PlaceOrder(c *gin.Context) {
	reply := errdef.New(c)
	param := &mshop.PlaceOrderReq{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("shop_trace: invalid param, err:%s", err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	
	userId, _ := c.Get(consts.USER_ID)
	param.UserId = userId.(string)
	svc := cshop.New(c)
	code, resp := svc.PlaceOrder(param)
	reply.Data["detail"] = resp
	reply.Response(http.StatusOK, code)
}

func OrderCancel(c *gin.Context) {
	reply := errdef.New(c)
	param := &mshop.ChangeOrderReq{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("shop_trace: invalid param, err:%s, param:%+v", err, param)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	
	userId, _ := c.Get(consts.USER_ID)
	param.UserId = userId.(string)
	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.OrderCancel(param))
}

func OrderList(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	reqType := c.Query("req_type")
	
	svc := cshop.New(c)
	code, list := svc.OrderList(userId.(string), reqType, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func ConfirmReceipt(c *gin.Context) {
	reply := errdef.New(c)
	param := &mshop.ChangeOrderReq{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("shop_trace: invalid param, param:%+v, err:%s", param, err)
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	
	userId, _ := c.Get(consts.USER_ID)
	param.UserId = userId.(string)
	svc := cshop.New(c)
	reply.Response(http.StatusOK, svc.ConfirmReceipt(param))
}

// 删除订单
func OrderDelete(c *gin.Context) {
	reply := errdef.New(c)
	param := &mshop.ChangeOrderReq{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("order_trace: invalid param, param:%+v, err:%s", param, err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}
	
	userId, _ := c.Get(consts.USER_ID)
	svc := cshop.New(c)
	param.UserId = userId.(string)
	reply.Response(http.StatusOK, svc.DeleteOrder(param))
}

// 订单详情
func OrderDetail(c *gin.Context) {
	reply := errdef.New(c)
	orderId := c.Query("order_id")
	userId, _ := c.Get(consts.USER_ID)
	
	svc := cshop.New(c)
	code, detail := svc.OrderDetail(userId.(string), orderId)
	
	
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, code)
}
