package contest

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
)

// 赛事模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	contest := api.Group("/contest")
	contest.Use(sign.CheckSign())
	{
		// 获取赛事模块banner列表
		contest.GET("/banner", BannerList)
		// 获取赛事直播列表
		contest.GET("/live/list", LiveList)
		// 首页推荐直播
		contest.GET("/recommend/live", RecommendLive)
	}
}
