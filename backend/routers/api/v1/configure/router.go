package configure

import "github.com/gin-gonic/gin"

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	configure := api.Group("/config")
	{
		// 添加banner
		configure.POST("/add/banner", AddBanner)
		// 删除banner
		configure.POST("/del/banner", DelBanner)
		// 获取banner列表
		configure.GET("/banners", GetBanners)
	}
}
