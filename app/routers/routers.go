package routers

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware"
	"sports_service/server/global/consts"
	"sports_service/server/app/routers/api/v1/user"
	"sports_service/server/app/routers/api/v1/swag"
	"sports_service/server/app/routers/api/v1/doc"
	"sports_service/server/app/routers/api/v1/client"
	"sports_service/server/app/config"
	"sports_service/server/global/app/log"
)

// 路由初始化
func InitRouters(engine *gin.Engine) {
	// 初始化中间件
	middleware.InitMiddleware(engine, log.Log, config.Global.Log.ShowColor)
	// 生成环境 不展示api文档 及 错误码文档
	if config.Global.Mode != string(consts.ModeProd) {
		// swag文档
		swag.Router(engine)
		// 错误码文档
		doc.Router(engine)
	}

	// 初始化接口
	client.Router(engine)
	// 用户账户
	user.Router(engine)

}
