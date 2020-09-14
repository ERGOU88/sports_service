package client

import (
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	client := api.Group("/client")
	{
		// 客户端初始化时 调用
		client.GET("/init", InitInfo)
	}
}
