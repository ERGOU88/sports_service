package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"sports_service/server/app/routers/api/v1/doc"
)

func InitRouter(engine *gin.Engine) {
	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "x-requested-with")
		c.Header("Access-Control-Allow-Headers", "Cookie")
		c.Header("Access-Control-Allow-Headers", "Authorization")
		c.Header("Access-Control-Allow-Headers", "auth")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		// 允许请求带有验证信息
		c.Header("Access-Control-Allow-Credentials", "true")
	})

	// swag生成接口文档
	swagger := engine.Group("/swagger")
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api/v1")
	// 错误码文档
	api.GET("/doc", doc.ApiCode)
}
