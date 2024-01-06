package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/backend/config"
	"sports_service/backend/routers/api/v1/admin"
	"sports_service/backend/routers/api/v1/comment"
	"sports_service/backend/routers/api/v1/configure"
	"sports_service/backend/routers/api/v1/contest"
	"sports_service/backend/routers/api/v1/course"
	"sports_service/backend/routers/api/v1/doc"
	"sports_service/backend/routers/api/v1/finance"
	"sports_service/backend/routers/api/v1/information"
	"sports_service/backend/routers/api/v1/notify"
	"sports_service/backend/routers/api/v1/post"
	"sports_service/backend/routers/api/v1/pub"
	"sports_service/backend/routers/api/v1/shop"
	"sports_service/backend/routers/api/v1/stat"
	"sports_service/backend/routers/api/v1/swag"
	"sports_service/backend/routers/api/v1/user"
	"sports_service/backend/routers/api/v1/venue"
	"sports_service/backend/routers/api/v1/video"
	"sports_service/global/backend/log"
	"sports_service/global/consts"
	"sports_service/middleware"
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

	engine.Use(gin.Recovery())
	engine.Any("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
		return
	})

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
	// 赛事模块
	contest.Router(engine)
	// 发布模块
	pub.Router(engine)
	// 商城模块
	shop.Router(engine)
	// 教育课程
	course.Router(engine)
}
