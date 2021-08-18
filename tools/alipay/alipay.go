package alipay

import (
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/gopay/pkg/xlog"
)

func TradeAppPay() {
	//aliPayPublicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1wn1sU/8Q0rYLlZ6sq3enrPZw2ptp6FecHR2bBFLjJ+sKzepROd0bKddgj+Mr1ffr3Ej78mLdWV8IzLfpXUi945DkrQcOUWLY0MHhYVG2jSs/qzFfpzmtut2Cl2TozYpE84zom9ei06u2AXLMBkU6VpznZl+R4qIgnUfByt3Ix5b3h4Cl6gzXMAB1hJrrrCkq+WvWb3Fy0vmk/DUbJEz8i8mQPff2gsHBE1nMPvHVAMw1GMk9ImB4PxucVek4ZbUzVqxZXphaAgUXFK2FSFU+Q+q1SPvHbUsjtIyL+cLA6H/6ybFF9Ffp27Y14AHPw29+243/SpMisbGcj2KD+evBwIDAQAB"
	privateKey := "MIIEpAIBAAKCAQEA29Ozn6SdfeqlCmsnNQ/rk6gvZr1jjkhDzOrYJXfOq/MGZOla3ymRSN8+6P0UV0o6QQDINNU/5ZtrJPtJzMOg7bQt8Z1ATidk/OhUjC9sTyBgwXxwjaikB73n9o/fgcgesz1ofkjDpZJeV76Cn5j4JENzaQP3xjvade3B307ReJOzEKVHjyyqfJ25DcGHqYEtdnbuYt9ijfaCfB0oEXJuUxUsIQQALzbfgTCVjdXKHjqCrRFKs1et8k12d0m2xoPTWf2YX72oFbRWqEZD13VhzL/Q2hU32/ENNXbHHrTZtTb3yFHw5Uj67pUmPzNhZFiHek0BkPVvKD89uErvIFpyWQIDAQABAoIBACoRo6iDmlhElX0e8IvpFg5V+2xQBkNudPs8Xk0dVoH1ql2Zgvh+Pf2SK7nu5Puniupxud7SiL3qNmEHbiIvthaHittYWrwaMetskvGZCcNC0QF2TRvvECUjJMc81WtC3w0yTVMNndOL5V4paVodri9ScT3BsqNPRQmYjKetr8zBLII6cpjRbT339RD7Z1FrNSKKQQWPGYH6JAd2sJR+hLHHylSvtpD2IVD3RJgVw11Ge9CKiDmCD6BkHIL5zPe3YErEWAscWsBnjF9sBdTLbkzlF83AIsDipS89dNF/WCCKejCj6A7Rl5kddu5jPSLfwCF6z0YJp7DCAzVZJXg1pgECgYEA9en/Krg9iHU5kP2d58OooHMz73srmtLhJiMINw7098r+1or309MqeSmlB4HGFUdl1D/zd5eQm3VJlixXSmlo5Ew+0xlalL2n+9fNJnbQjEieUFHKdPaBMRW7FQ7pD0UmaTWot1HxZeAt3rkNA0e6oxYPPtOjpswEP/nnvkS17mECgYEA5NfM6xdvB5aITe1uHsIn8KAoBVFuUZu9GgTZgKQgR/JW1QZ5Ba/DCRSK/gmQWbOQHhi6hyLqcUUNaQbuA/1T7ByNCCdk9Mr/zmehyh+BdYj12SQ8i3Af2H45l4xX3JDAXa9l+ba+HuaJ5W4FJcYiynGnG+uNgxNTsI+OckTxVvkCgYEAt0jWZDK5uhEU/NnqbSlJb30twlpdH6H5KYGGx/Kf5mgoFCOznu+Ogovlcnjo+EckwFOB1SrkHtoGJKWb0dxKz418bb5B4waQQ4aOYxK/US92v4qWiSKJG9qEe6eHUVhKzrOtsiSi9TlnNs9ZwY4erxrr9fmryc/Zgw1yCkAQEUECgYEA464hXzUtbmtCqeW0Tj315t4xczkVfXRprF1u2SJyS6K86a1K83FvprUdpKp3SAfzNz57NsByaMe/E+OlI6sDuEKfvqETPMpLwFwzCBpYf0wI7kWzRzgDNy4+tp0XPYd3HL7Jwq0iczQDtpTD4lVDgA+bp5ewb9zmwx/RJbeaNmECgYAaooellJxLAXT0HFsKWmIMn/ckQu+wXcg+sqFOkhKCOMXNnGiADWT5LcyEVcejWQAru3QsvTzPdundFbQR+0hSWDFiPITUopINLATDW0zvSUcbL27Z30HZUXd+LASmU/EC0I6rNbsTZ/nyImpfOul7v2T/BS+eVRjevDrFMGtoeQ=="
	// 初始化支付宝客户端
	// appId：应用ID
	// privateKey：应用私钥，支持PKCS1和PKCS8
	// isProd：是否是正式环境
	client := alipay.NewClient("2021002164657197", privateKey, false)
	//配置公共参数
	client.SetCharset("utf-8").
		SetSignType(alipay.RSA2).
		SetPrivateKeyType(alipay.PKCS1).
		SetNotifyUrl("https://fpv-app-api.bluetrans.cn")

	// 请求参数
	body := make(gopay.BodyMap)
	body.Set("subject", "测试APP支付")
	body.Set("out_trade_no", "202111111111")
	body.Set("total_amount", "0.10")
	// 手机APP支付参数请求
	payParam, err := client.TradeAppPay(body)
	if err != nil {
		xlog.Error("err:", err)
		return
	}
	xlog.Debug("payParam:", payParam)
}
