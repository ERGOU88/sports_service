package shop

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
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
		shop.POST("/edit/addr", token.TokenAuth(), EditAddr)
		// 用户地址列表
		shop.GET("/user/addr/list", token.TokenAuth(), UserAddrList)
		// 添加商品购物车
		shop.POST("/add/product/cart", token.TokenAuth(), AddProductCart)
		// 获取商品购物车列表
		shop.GET("/product/cart", token.TokenAuth(), ProductCart)
		// 搜索商品
		shop.GET("/product/search", SearchProduct)
		// 更新商品购物车
		shop.POST("/update/product/cart", token.TokenAuth(), UpdateProductCart)
		// 下单
		shop.POST("/place/order", token.TokenAuth(), PlaceOrder)
		// 订单取消
		shop.POST("/order/cancel", token.TokenAuth(), OrderCancel)
		// 商城订单列表
		shop.GET("/order/list", token.TokenAuth(), OrderList)
		// 确认收货
		shop.POST("/confirm/receipt", token.TokenAuth(), ConfirmReceipt)
	}
}
