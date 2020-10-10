package tencentCloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"sports_service/server/util"
	"time"
	"fmt"
)

type TencentCloud struct {
	secretId  string
	secretKey string
	apiDomain string
}

func New(secretId, secretKey, apiDomain string) (client *TencentCloud) {
	client = &TencentCloud{
		secretId: secretId,
		secretKey: secretKey,
		apiDomain: apiDomain,
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

// 主动拉取事件通知
func (tc *TencentCloud) PullEvents() (*v20180717.PullEventsResponse, error){
	credential := common.NewCredential(
		tc.secretId,
		tc.secretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 30
	cpf.HttpProfile.Endpoint = tc.apiDomain
	client, _ := v20180717.NewClient(credential, "ap-shanghai", cpf)
	req := v20180717.NewPullEventsRequest()
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.PullEvents(req)
	// 处理异常
	fmt.Println(err)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	// 非SDK异常，直接失败。
	if err != nil {
		fmt.Printf("Request Pull Event error: %s", err)
		return nil, err
	}

	// 打印返回的json字符串
	fmt.Printf("%s", response.ToJsonString())

	return response, nil
}

// 确认事件通知
func (tc *TencentCloud) ConfirmEvents() error {
	credential := common.NewCredential(
		tc.secretId,
		tc.secretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 30
	cpf.HttpProfile.Endpoint = tc.apiDomain
	client, _ := v20180717.NewClient(credential, "ap-shanghai", cpf)
	req := v20180717.NewConfirmEventsRequest()
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.ConfirmEvents(req)
	// 处理异常
	fmt.Println(err)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return err
	}
	// 非SDK异常，直接失败。
	if err != nil {
		fmt.Printf("Request Confirm Event error: %s", err)
		return err
	}

	// 打印返回的json字符串
	fmt.Printf("%s", response.ToJsonString())
	return nil
}

