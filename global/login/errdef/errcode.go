package errdef

const (
	SUCCESS            = 200
	ERROR              = 500
	INVALID_PARAMS     = 400

	// 用户相关错误码 1000-2000
	INVALID_MOBILE_NUM  = 1000
	USER_ALREADY_EXISTS = 1001
	USER_REPEAT_REG     = 1002
	USER_REGISTER_FAIL  = 1003
)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	INVALID_MOBILE_NUM:  "非法的手机号",
	USER_ALREADY_EXISTS: "用户已存在",
	USER_REPEAT_REG:     "用户重复注册",
	USER_REGISTER_FAIL:  "用户注册失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}


