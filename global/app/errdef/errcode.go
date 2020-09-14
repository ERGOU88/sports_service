package errdef

const (
	SUCCESS            = 200
	ERROR              = 500
	INVALID_PARAMS     = 400
	UNAUTHORIZED       = 401

	// 用户相关错误码 1000-2000
	FREE_LOGIN_FAIL         = 1000
	INVALID_MOBILE_NUM      = 1001
	USER_ALREADY_EXISTS     = 1002
	USER_REPEAT_REG         = 1003
	USER_REGISTER_FAIL      = 1004
	USER_ADD_INFO_FAIL      = 1005
	USER_GET_INFO_FAIL      = 1006

	WX_USER_INFO_FAIL       = 1100
	WX_ACCESS_TOKEN_FAIL    = 1101
	WX_REGISTER_FAIL        = 1102
	WX_ADD_ACCOUNT_FAIL     = 1103

	WEIBO_USER_INFO_FAIL    = 1200
	WEIBO_ADD_ACCOUNT_FAIL  = 1201
	WEIBO_REGISTER_FAIL     = 1202

	QQ_UNIONID_FAIL         = 1301
	QQ_USER_INFO_FAIL       = 1302
	QQ_REGISTER_FAIL        = 1303
	QQ_ADD_ACCOUNT_FAIL     = 1304

)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",
	UNAUTHORIZED:   "未经授权",

	FREE_LOGIN_FAIL:     "一键登陆失败",
	INVALID_MOBILE_NUM:  "非法的手机号",
	USER_ALREADY_EXISTS: "用户已存在",
	USER_REPEAT_REG:     "用户重复注册",
	USER_REGISTER_FAIL:  "用户注册失败",
	USER_ADD_INFO_FAIL:  "添加用户信息失败",
	USER_GET_INFO_FAIL:  "获取用户信息失败",

	WX_USER_INFO_FAIL:    "获取微信用户信息失败",
	WX_ACCESS_TOKEN_FAIL: "获取微信授权token失败",
	WX_REGISTER_FAIL:     "微信注册帐号失败",
	WX_ADD_ACCOUNT_FAIL:  "微信帐号添加失败",

	WEIBO_USER_INFO_FAIL:   "获取微博用户信息失败",
	WEIBO_ADD_ACCOUNT_FAIL: "记录微博登陆信息失败",
	WEIBO_REGISTER_FAIL:    "微博用户注册失败",

	QQ_UNIONID_FAIL:       "获取QQ授权信息失败",
	QQ_USER_INFO_FAIL:     "获取QQ用户信息失败",
	QQ_REGISTER_FAIL:      "QQ注册账户失败",
	QQ_ADD_ACCOUNT_FAIL:   "QQ账户添加失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}



