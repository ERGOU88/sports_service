package errdef

const (
	SUCCESS            = 200
	ERROR              = 500
	INVALID_PARAMS     = 400

	// 用户相关错误码 1000-2000
	INVALID_MOBILE_NUM   = 1000
	USER_ALREADY_EXISTS  = 1001
	USER_REPEAT_REG      = 1002
	USER_REGISTER_FAIL   = 1003
	USER_ADD_INFO_FAIL   = 1004
	USER_GET_INFO_FAIL   = 1005

	WX_ACCESS_TOKEN_FAIL = 1101
	WX_REGISTER_FAIL     = 1102
	WX_ADD_ACCOUNT_FAIL  = 1103
)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	INVALID_MOBILE_NUM:  "非法的手机号",
	USER_ALREADY_EXISTS: "用户已存在",
	USER_REPEAT_REG:     "用户重复注册",
	USER_REGISTER_FAIL:  "用户注册失败",
	USER_ADD_INFO_FAIL:  "添加用户信息失败",
	USER_GET_INFO_FAIL:  "获取用户信息失败",

	WX_ACCESS_TOKEN_FAIL: "获取微信授权token失败",
	WX_REGISTER_FAIL:     "微信注册帐号失败",
	WX_ADD_ACCOUNT_FAIL:  "微信帐号添加失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}


