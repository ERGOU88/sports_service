package routers

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/backend/routers/api/v1/doc"
	"sports_service/server/backend/routers/api/v1/swag"
	"sports_service/server/backend/config"
	"sports_service/server/backend/routers/api/v1/video"
	"sports_service/server/global/consts"
)

func InitRouters(engine *gin.Engine) {
	// 生产环境 不展示api文档 及 错误码文档
	if config.Global.Mode != string(consts.ModeProd) {
		// swag文档
		swag.Router(engine)
		// 错误码文档
		doc.Router(engine)
	}

	// 后台点播模块
	video.Router(engine)

}
