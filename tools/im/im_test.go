package im

import (
	"testing"
)


// 文本检测
func TestGenSign(t *testing.T) {
	sig, err := GenSig(86400)
	if err != nil {
		return
	}

	url := GenRequestUrl(sig, "/v4/im_open_login_svc/account_import")
	t.Logf("url:%s", url)
}
