package tencentCloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/parnurzeal/gorequest"
	"sports_service/server/util"
	"time"
	"fmt"
)

type TencentCloud struct {
	secretId  string
	secretKey string
	apiDomain string
	Client    *gorequest.SuperAgent
}

func New(secretId, secretKey, apiDomain string) (client *TencentCloud) {
	client = &TencentCloud{
		secretId: secretId,
		secretKey: secretKey,
		apiDomain: apiDomain,
		Client: gorequest.New(),
	}

	return client
}

type SourceContext struct {
	UserId    string   `json:"user_id"`   // 用户id
	TaskId    int64    `json:"task_id"`   // 任务id
}

// 生成上传签名
func (tc *TencentCloud) GenerateSign(userId string, taskId int64) string {
	timestamp := time.Now().Unix()
	expireTime := timestamp + ONE_DAY
	random := util.GetXID()
	sourceContext := &SourceContext{
		UserId: userId,
		TaskId:	taskId,
	}

	context, _ := util.JsonFast.Marshal(sourceContext)
	original := fmt.Sprintf("secretId=%s&currentTimeStamp=%d&sourceContext=%s&expireTime=%d&random=%d", tc.secretId, timestamp, string(context), expireTime, random)
	signature := tc.generateHmacSHA1(original)
	signature = append(signature, []byte(original)...)
	signatureB64 := base64.StdEncoding.EncodeToString(signature)
	fmt.Println(signatureB64)
	return signatureB64
}

func (tc *TencentCloud) generateHmacSHA1(original string) []byte {
	mac := hmac.New(sha1.New, []byte(tc.secretKey))
	sha1.New()
	mac.Write([]byte(original))
	return mac.Sum(nil)
}
