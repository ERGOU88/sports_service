package shop

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
)

// 分享/转发模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	shop := api.Group("/shop")
	shop.Use(sign.CheckSign())
	{
		// 商品分类配置
		shop.GET("/product/category", ProductCategory)
		// 获取分类下的商品
		shop.GET("/products", ProductList)
		// 推荐的商品
		shop.GET("/product/recommend", RecommendProduct)
		// 商品详情
		shop.GET("/product/detail", ProductDetail)
		// 地址配置
		shop.GET("/area/config", AreaConfig)
		// 用户添加/更新 地址信息
		shop.POST("/edit/area", EditArea)
	}
}
