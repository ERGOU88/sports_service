package tencentCloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms/v20200713"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	vodSdk "github.com/tencentyun/vod-go-sdk"
	"sports_service/server/app/config"
	"sports_service/server/util"
	"time"
	"strings"
)

// 透传数据
type SourceContext struct {
	UserId    string   `json:"user_id"`   // 用户id
	TaskId    int64    `json:"task_id"`   // 任务id
	Tm        int64    `json:"tm"`        // 任务开始时间
	Mode      string   `json:"mode"`      // 环境
}

// 生成上传签名 todo: 任务流模版名  procedure
func (tc *TencentCloud) GenerateSign(userId, procedure string, taskId int64) string {
	timestamp := time.Now().Unix()
	expireTime := timestamp + ONE_DAY
	random := util.GetXID()
	sourceContext := &SourceContext{
		UserId: userId,
		TaskId:	taskId,
		Tm: time.Now().Unix(),
		Mode: config.Global.Mode,
	}

	context, _ := util.JsonFast.Marshal(sourceContext)
	original := fmt.Sprintf("secretId=%s&currentTimeStamp=%d&procedure=%s&sourceContext=%s&expireTime=%d&random=%d", tc.secretId, timestamp, procedure, string(context), expireTime, random)
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
func (tc *TencentCloud) PullEvents() (*vod.PullEventsResponse, error){
	credential := common.NewCredential(
		tc.secretId,
		tc.secretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 30
	cpf.HttpProfile.Endpoint = tc.apiDomain
	client, _ := vod.NewClient(credential, "ap-shanghai", cpf)
	req := vod.NewPullEventsRequest()
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.PullEvents(req)
	// 处理异常
	//fmt.Println(err)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		//fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	// 非SDK异常，直接失败。
	if err != nil {
		//fmt.Printf("Request Pull Event error: %s", err)
		return nil, err
	}

	// 打印返回的json字符串
	//fmt.Printf("%s", response.ToJsonString())

	return response, nil
}

// 确认事件通知
// EventHandles 事件句柄，即 拉取事件通知 接口输出参数中的 EventSet. EventHandle 字段。
// 数组长度限制：16。
func (tc *TencentCloud) ConfirmEvents(events []string) error {
	credential := common.NewCredential(
		tc.secretId,
		tc.secretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 30
	cpf.HttpProfile.Endpoint = tc.apiDomain
	client, _ := vod.NewClient(credential, "ap-shanghai", cpf)
	req := vod.NewConfirmEventsRequest()
	ps := common.StringPtrs(events)
	req.EventHandles = ps
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.ConfirmEvents(req)
	// 处理异常
	//fmt.Println(err)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		//fmt.Printf("An API error has returned: %s", err)
		return err
	}
	// 非SDK异常，直接失败。
	if err != nil {
		//fmt.Printf("Request Confirm Event error: %s", err)
		return err
	}

	// 打印返回的json字符串
	fmt.Printf("%s", response.ToJsonString())
	return nil
}

// 文本内容检测
func (tc *TencentCloud) TextModeration(content string) (bool, string, error) {
	if content == "" {
		return true, content, nil
	}

	credential := common.NewCredential(
		tc.secretId,
		tc.secretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 30
	cpf.HttpProfile.Endpoint = tc.apiDomain
	client, _ := tms.NewClient(credential, "ap-shanghai", cpf)
	req := tms.NewTextModerationRequest()
	req.Content = common.StringPtr(base64.StdEncoding.EncodeToString([]byte(content)))
	// 通过client对象调用想要访问的接口，需要传入请求对象
	response, err := client.TextModeration(req)
	// 处理异常
	//fmt.Println(err)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		//fmt.Printf("An API error has returned: %s", err)
		return false, content, err
	}
	// 非SDK异常，直接失败。
	if err != nil {
		//fmt.Printf("Request Pull Event error: %s", err)
		return false, content, err
	}

	// 打印返回的json字符串
	//fmt.Printf("%s", response.ToJsonString())
	// Label Normal：正常，Polity：涉政，Porn：色情，Illegal：违法，Abuse：谩骂，Terror：暴恐，Ad：广告，Custom：自定义关键词
	// Suggestion Block：建议打击，Review：建议复审，Normal：建议通过。
	if *response.Response.Suggestion == "Block" {
		//fmt.Printf("Content Not Pass, Label:%s, Suggestion:%s, Content:%s",
			//*response.Response.Label, *response.Response.Suggestion, content)
		if len(response.Response.Keywords) > 0 {
			for _, str := range response.Response.Keywords {
				content = strings.Replace(content, *str, "***", -1)
			}
		}

		return false, content, nil
	}

	return true, content, nil
}

// taskId 任务id
// userId 用户id
// path: 文件路径
// region: 区域 例如 ap-shanghai
// procedure: 任务流模版名称
func (tc *TencentCloud) Upload(taskId int64, userId, token, path, region, procedure string) (*vodSdk.VodUploadResponse, error) {
	vodClient := &vodSdk.VodUploadClient{}
	vodClient.SecretId = tc.secretId
	vodClient.SecretKey = tc.secretKey
	vodClient.Token = token
	vodClient.Timeout = 30

	req := vodSdk.NewVodUploadRequest()
	req.MediaFilePath = common.StringPtr(path)
	req.Procedure = common.StringPtr(procedure)
	context, _ := util.JsonFast.Marshal(&SourceContext{
		TaskId: taskId,
		UserId: userId,
	})
	req.SourceContext = common.StringPtr(string(context))

	rsp, err := vodClient.Upload(region, req)
	if err != nil {
		//fmt.Printf("Request Upload error: %s", err)
		return nil, err
	}

	//fmt.Println(*rsp.Response.FileId)
	//fmt.Println(*rsp.Response.MediaUrl)
	//fmt.Println(*rsp.Response.CoverUrl)

	return rsp, nil
}
