package information

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
)

// 资讯模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	information := api.Group("/information")
	information.Use(sign.CheckSign())
	{
		// 获取资讯列表
		information.GET("/list", InformationList)
		// 获取资讯详情
		information.GET("/detail", InformationDetail)
	}
}
