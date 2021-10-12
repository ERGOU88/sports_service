package routers

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/backend/routers/api/v1/finance"
	"sports_service/server/backend/routers/api/v1/venue"
	"sports_service/server/backend/routers/api/v1/information"
	"sports_service/server/backend/routers/api/v1/comment"
	"sports_service/server/backend/routers/api/v1/doc"
	"sports_service/server/backend/routers/api/v1/post"
	"sports_service/server/backend/routers/api/v1/swag"
	"sports_service/server/backend/config"
	"sports_service/server/backend/routers/api/v1/user"
	"sports_service/server/backend/routers/api/v1/video"
	"sports_service/server/backend/routers/api/v1/admin"
	"sports_service/server/backend/routers/api/v1/notify"
	"sports_service/server/global/consts"
	"sports_service/server/backend/routers/api/v1/configure"
	"sports_service/server/middleware"
	"sports_service/server/global/backend/log"
	"sports_service/server/backend/routers/api/v1/stat"
)

func InitRouters(engine *gin.Engine) {
	// 初始化中间件
	middleware.InitMiddleware(engine, log.Log, config.Global.Log.ShowColor)
	// 生产环境 不展示api文档 及 错误码文档
	if config.Global.Mode != string(consts.ModeProd) {
		// swag文档
		swag.Router(engine)
		// 错误码文档
		doc.Router(engine)
	}

	// 后台点播模块
	video.Router(engine)
	// 后台评论模块
	comment.Router(engine)
	// 后台配置模块（banner等）
	configure.Router(engine)
	// 后台用户管理模块
	user.Router(engine)
	// 后台管理员模块
	admin.Router(engine)
	// 站内信模块
	notify.Router(engine)
	// 帖子模块
	post.Router(engine)
	// 资讯模块
	information.Router(engine)
	// 场馆模块
	venue.Router(engine)
	// 统计模块
	stat.Router(engine)
	// 财务模块
	finance.Router(engine)
}
