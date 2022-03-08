package shop

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/jwt"
)

// 商城模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	shop := api.Group("/shop")
	shop.Use(jwt.JwtAuth())
	{
		// 商品分类配置
		shop.GET("/product/category", ProductCategory)
		// 商品列表
		shop.GET("/product/list", ProductList)
		// 商品详情
		shop.GET("/product/detail", ProductDetail)
		// 添加分类
		shop.POST("/add/category", AddCategory)
		// 编辑分类
		shop.POST("/edit/category", EditCategory)
		// 服务列表
		shop.GET("/service/list", ServiceList)
		// 添加服务
		shop.POST("/add/service", AddService)
		// 编辑服务
		shop.POST("/edit/service", EditService)
		// 删除服务
		shop.DELETE("/del/service", DelService)
		// 添加分类规格
		shop.POST("/add/specification", AddSpecification)
		// 编辑分类规格
		shop.POST("/edit/specification", EditSpecification)
		// 删除分类规格
		shop.DELETE("/del/specification", DelSpecification)
		// 规格信息
		shop.GET("/specification/info", SpecificationInfo)
		// 分类规格列表
		shop.GET("/specification/list", SpecificationList)
		// 添加商品
		shop.POST("/add/product", AddProduct)
		// 编辑商品
		shop.POST("/edit/product", EditProduct)
		// 订单列表
		shop.GET("/order/list", OrderList)
		// 发货
		shop.POST("/deliver/product", DeliverProduct)
		// 确认收货
		shop.POST("/confirm/receipt", ConfirmReceipt)
		// 订单回调 [掉单时 人工操作]
		shop.POST("/order/callback", OrderCallback)
	}
}
