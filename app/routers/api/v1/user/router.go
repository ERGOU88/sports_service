package user

import (
  "github.com/gin-gonic/gin"
  "sports_service/server/middleware/sign"
  "sports_service/server/middleware/token"
)

// 用户账户模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	user := api.Group("/user")
	// todo: 先不校验签名
	{
		// 获取短信验证码
		user.POST("/smscode", sign.CheckSign(), SmsCode)
		// 短信验证码登陆
		user.POST("/smscode/login", sign.CheckSign(), SmsCodeLogin)
		// 手机一键登陆
		user.POST("/mobile/login", sign.CheckSign(), MobilePhoneLogin)
		// 用户微信登陆
		user.POST("/wechat/login", sign.CheckSign(), UserWechatLogin)
		// 用户微博登陆
		user.POST("/weibo/login", sign.CheckSign(), UserWeiboLogin)
		// 用户QQ登陆
		user.POST("/qq/login", sign.CheckSign(), UserQQLogin)
		// 用户信息
		user.GET("/info", sign.CheckSign(), token.TokenAuth(), UserInfo)
		// 修改用户信息
		user.POST("/edit/info", sign.CheckSign(), token.TokenAuth(), EditUserInfo)
		// 用户反馈
		user.POST("/feedback", token.TokenAuth(), UserFeedback)
		// 个人空间信息
		user.GET("/zone/info", sign.CheckSign(), UserZoneInfo)
    // 绑定设备token
    user.POST("/bind/deviceToken", sign.CheckSign(), token.TokenAuth(), BindDeviceToken)
    // 版本更新（数据库控制）
    user.GET("/version/up", sign.CheckSign(), VersionUp)
	}

}
