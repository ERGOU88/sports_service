package venue

import (
	"github.com/gin-gonic/gin"
)

// 场馆模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	venue := api.Group("/venue")
	{
		// 场馆信息
		venue.GET("/info", VenueInfo)
	}

}

