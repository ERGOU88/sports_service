package venue

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
	"sports_service/server/middleware/token"
)

// 场馆模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	venue := api.Group("/venue")
	venue.Use(sign.CheckSign())
	{
		// 场馆信息
		venue.GET("/info", VenueInfo)
        // 购买次卡/月卡/年卡
		venue.POST("/purchase/vipCard", token.TokenAuth(), PurchaseVipCard)
        // 进出场记录
		venue.GET("/action/record", token.TokenAuth(), ActionRecord)
	}

}

