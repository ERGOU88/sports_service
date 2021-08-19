package alipay

import (
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"sports_service/server/global/app/log"
)

const (
	ALIPAY_APP_ID        = "2021002164657197"
	ALIPAY_PRIVATE_KEY   = "MIIEpAIBAAKCAQEA29Ozn6SdfeqlCmsnNQ/rk6gvZr1jjkhDzOrYJXfOq/MGZOla3ymRSN8+6P0UV0o6QQDINNU/5ZtrJPtJzMOg7bQt8Z1ATidk/OhUjC9sTyBgwXxwjaikB73n9o/fgcgesz1ofkjDpZJeV76Cn5j4JENzaQP3xjvade3B307ReJOzEKVHjyyqfJ25DcGHqYEtdnbuYt9ijfaCfB0oEXJuUxUsIQQALzbfgTCVjdXKHjqCrRFKs1et8k12d0m2xoPTWf2YX72oFbRWqEZD13VhzL/Q2hU32/ENNXbHHrTZtTb3yFHw5Uj67pUmPzNhZFiHek0BkPVvKD89uErvIFpyWQIDAQABAoIBACoRo6iDmlhElX0e8IvpFg5V+2xQBkNudPs8Xk0dVoH1ql2Zgvh+Pf2SK7nu5Puniupxud7SiL3qNmEHbiIvthaHittYWrwaMetskvGZCcNC0QF2TRvvECUjJMc81WtC3w0yTVMNndOL5V4paVodri9ScT3BsqNPRQmYjKetr8zBLII6cpjRbT339RD7Z1FrNSKKQQWPGYH6JAd2sJR+hLHHylSvtpD2IVD3RJgVw11Ge9CKiDmCD6BkHIL5zPe3YErEWAscWsBnjF9sBdTLbkzlF83AIsDipS89dNF/WCCKejCj6A7Rl5kddu5jPSLfwCF6z0YJp7DCAzVZJXg1pgECgYEA9en/Krg9iHU5kP2d58OooHMz73srmtLhJiMINw7098r+1or309MqeSmlB4HGFUdl1D/zd5eQm3VJlixXSmlo5Ew+0xlalL2n+9fNJnbQjEieUFHKdPaBMRW7FQ7pD0UmaTWot1HxZeAt3rkNA0e6oxYPPtOjpswEP/nnvkS17mECgYEA5NfM6xdvB5aITe1uHsIn8KAoBVFuUZu9GgTZgKQgR/JW1QZ5Ba/DCRSK/gmQWbOQHhi6hyLqcUUNaQbuA/1T7ByNCCdk9Mr/zmehyh+BdYj12SQ8i3Af2H45l4xX3JDAXa9l+ba+HuaJ5W4FJcYiynGnG+uNgxNTsI+OckTxVvkCgYEAt0jWZDK5uhEU/NnqbSlJb30twlpdH6H5KYGGx/Kf5mgoFCOznu+Ogovlcnjo+EckwFOB1SrkHtoGJKWb0dxKz418bb5B4waQQ4aOYxK/US92v4qWiSKJG9qEe6eHUVhKzrOtsiSi9TlnNs9ZwY4erxrr9fmryc/Zgw1yCkAQEUECgYEA464hXzUtbmtCqeW0Tj315t4xczkVfXRprF1u2SJyS6K86a1K83FvprUdpKp3SAfzNz57NsByaMe/E+OlI6sDuEKfvqETPMpLwFwzCBpYf0wI7kWzRzgDNy4+tp0XPYd3HL7Jwq0iczQDtpTD4lVDgA+bp5ewb9zmwx/RJbeaNmECgYAaooellJxLAXT0HFsKWmIMn/ckQu+wXcg+sqFOkhKCOMXNnGiADWT5LcyEVcejWQAru3QsvTzPdundFbQR+0hSWDFiPITUopINLATDW0zvSUcbL27Z30HZUXd+LASmU/EC0I6rNbsTZ/nyImpfOul7v2T/BS+eVRjevDrFMGtoeQ=="
)

type AliPayClient struct {
	Client       *alipay.Client
	Subject      string        // 商品名称
	OutTradeNo   string        // 订单号
	TotalAmount  string        // 金额   （元）
}

// appId：应用ID
// privateKey：应用私钥，支持PKCS1和PKCS8
// isProd：是否是正式环境
func NewAliPay(isProd bool) *AliPayClient {
	ali := &AliPayClient{}
	ali.Client = alipay.NewClient(ALIPAY_APP_ID, ALIPAY_PRIVATE_KEY, isProd)
	// 配置公共参数
	ali.Client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetPrivateKeyType(alipay.PKCS1)
	return ali
}

// app支付
func (c *AliPayClient) TradeAppPay() (string, error) {
	// 请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", c.Subject)
	body.Set("out_trade_no", c.OutTradeNo)
	body.Set("total_amount", c.TotalAmount)
	// 手机APP支付参数请求
	payParam, err := c.Client.TradeAppPay(body)
	if err != nil {
		log.Log.Errorf("ali_trace: trade app pay fail, err:%s, orderId:%s", err, c.OutTradeNo)
		return "", err
	}

	return payParam, nil
}

// 校验签名
func (c *AliPayClient) VerifySign(body interface{}) (bool, error) {
	return false, nil
}
