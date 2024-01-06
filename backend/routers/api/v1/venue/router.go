package venue

import (
	"github.com/gin-gonic/gin"
	"sports_service/middleware/jwt"
)

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	venue := api.Group("/venue")
	venue.Use(jwt.JwtAuth())
	{
		// 获取场馆列表
		venue.GET("/list", VenueList)
		// 场馆详情
		venue.GET("/detail", VenueDetail)
		// 编辑场馆
		venue.POST("/edit", EditVenue)
		// 更新退款费率
		venue.POST("/refund/rate", UpdateRefundRate)
		// 退款规则
		venue.GET("/refund/rules", RefundRules)
		// 添加场馆
		venue.POST("/add", AddVenue)
		// 添加场馆角标配置
		venue.POST("/add/mark", AddMark)
		// 删除场馆角标配置
		venue.POST("/del/mark", DelMark)
		// 角标列表
		venue.GET("/mark/list", MarkList)
		// 添加店长
		venue.POST("/add/store/manager", AddStoreManager)
		// 编辑店长
		venue.POST("/edit/store/manager", EditStoreManager)
		// 店长列表
		venue.GET("/store/manager/list", StoreManagerList)
	}
}
