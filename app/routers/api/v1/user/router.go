package user

import (
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	user := api.Group("/user")
	{
		// 手机一键登陆
		user.POST("/mobile/login", MobilePhoneLogin)
		// 用户微信登陆
		user.POST("/wechat/login", UserWechatLogin)
		// 用户微博登陆
		user.POST("/weibo/login", UserWeiboLogin)
		// 用户QQ登陆
		user.POST("/qq/login", UserQQLogin)
	}
}
