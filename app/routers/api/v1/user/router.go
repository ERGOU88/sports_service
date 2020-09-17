package user

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/token"
)

// 用户账户模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	user := api.Group("/user")
	// todo: 先不校验签名
	//user.Use(sign.CheckSign())
	{
		// 手机一键登陆
		user.POST("/mobile/login", MobilePhoneLogin)
		// 用户微信登陆
		user.POST("/wechat/login", UserWechatLogin)
		// 用户微博登陆
		user.POST("/weibo/login", UserWeiboLogin)
		// 用户QQ登陆
		user.POST("/qq/login", UserQQLogin)
		// 用户信息
		user.GET("/info", token.TokenAuth(), UserInfo)
		// 修改用户信息
		user.POST("/edit/info", token.TokenAuth(), EditUserInfo)
		// 用户反馈
		user.POST("/feedback", token.TokenAuth(), UserFeedback)
	}

}
