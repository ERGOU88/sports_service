package venue

import "github.com/gin-gonic/gin"

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	venue := api.Group("/venue")
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
	}
}
