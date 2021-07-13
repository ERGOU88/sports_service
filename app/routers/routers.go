package routers

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/app/routers/api/v1/attention"
	"sports_service/server/app/routers/api/v1/barrage"
	"sports_service/server/app/routers/api/v1/collect"
	"sports_service/server/app/routers/api/v1/comment"
	"sports_service/server/app/routers/api/v1/like"
	"sports_service/server/app/routers/api/v1/notify"
	"sports_service/server/app/routers/api/v1/posting"
	"sports_service/server/app/routers/api/v1/search"
	"sports_service/server/app/routers/api/v1/video"
    "sports_service/server/app/routers/api/v1/tencentCloud"
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
	// 生产环境 不展示api文档 及 错误码文档
	if config.Global.Mode != string(consts.ModeProd) {
		// swag文档
		swag.Router(engine)
		// 错误码文档
		doc.Router(engine)
	}

	// 初始化接口
	client.Router(engine)
	// 用户体系
	user.Router(engine)
	// 关注模块
	attention.Router(engine)
	// 收藏模块
	collect.Router(engine)
	// 视频模块
	video.Router(engine)
	// 点赞模块
	like.Router(engine)
	// 通知模块
	notify.Router(engine)
	// 搜索模块
	search.Router(engine)
	// 评论模块
	comment.Router(engine)
	// 弹幕模块
	barrage.Router(engine)
	// 腾讯云
	tencentCloud.Router(engine)
	// 帖子模块
	posting.Router(engine)
}
