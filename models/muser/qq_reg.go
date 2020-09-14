package muser

type qqRegister struct {
	*base
}

// qq登陆请求参数
type QQLoginParams struct {
	Openid       string     `json:"openid"`
	AccessToken  string     `json:"access_token" example:"授权token"`
	Platform     string     `json:"platform" example:"qzone、pengyou或qplus"`       // qzone、pengyou或qplus
}

// 实例
func NewQQRegister() *qqRegister {
	return &qqRegister{
		&base{},
	}
}


