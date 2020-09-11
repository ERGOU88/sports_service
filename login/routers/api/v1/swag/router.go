package swag

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "sports_service/server/login/docs"
)

func Router(engine *gin.Engine) {
	swagger := engine.Group("/swagger")
	// swag生成接口文档
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
