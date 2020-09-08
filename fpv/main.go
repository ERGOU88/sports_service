package main

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"sports_service/server/fpv/routers"
)

func init() {
	// todo 初始化配置、日志、mysql等
}

func main() {
	// 启动服务
	engine := gin.New()
	routers.InitRouter(engine)
	if err := engine.Run(":11010"); err != nil {
		fmt.Printf("engine.Run err:%v", err)
		return
	}
}
