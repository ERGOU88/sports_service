package shop

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/jwt"
)

// 商城模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	shop := api.Group("/shop")
	shop.Use(jwt.JwtAuth())
	{
		// 商品分类配置
		shop.GET("/product/category", ProductCategory)
		// 商品列表
		shop.GET("/product/list", ProductList)
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
		// 添加分类规格
		shop.POST("/add/specification", AddSpecification)
		// 编辑分类规格
		shop.POST("/edit/specification", EditSpecification)
		// 删除分类规格
		shop.DELETE("/del/specification", DelSpecification)
	}
}
