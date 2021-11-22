package tencentCloud

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"net/http"
	"net/url"
	"sports_service/server/util"
	"time"
	"errors"
	"fmt"
)

// 获取腾讯对象存储临时通行证
func (tc *TencentCloud) GetCosTempAccess(region string) (map[string]interface{}, error) {
	credential := sts.NewClient(
		tc.secretId,
		tc.secretKey,
		nil,
	)

	opt := &sts.CredentialOptions{
		DurationSeconds: int64(2 * time.Hour.Seconds()),
		Region:          region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
						// 分片上传
						"name/cos:InitiateMultipartUpload",
						"name/cos:ListMultipartUploads",
						"name/cos:ListParts",
						"name/cos:UploadPart",
						"name/cos:CompleteMultipartUpload",
					},
					Effect: "allow",
					Resource: []string{
						//这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						"qcs::cos:" + region + ":uid/" + APPID + ":" + BUCKET + "/" + CRPATH + "/*",
						//"allowPrefixes:fpv/images/*:fpv/videos/*",
					},
				},
			},
		},
	}

	res, err := credential.GetCredential(opt)
	if err != nil {
		return nil, err
	}

	credentials := map[string]interface{}{
		"tmp_secret_id":  res.Credentials.TmpSecretID,
		"tmp_secret_key": res.Credentials.TmpSecretKey,
		"session_token":  res.Credentials.SessionToken,
	}

	resp := map[string]interface{}{
		"credentials":  credentials,
		"expired_time": res.ExpiredTime,
		"start_time":   res.StartTime,
		"dir":          CRPATH,
	}

	return resp, nil

}

// 图片审核
func (tc *TencentCloud) RecognitionImage(baseUrl, path string) (*cos.ImageRecognitionResult, error) {
	u, _ := url.Parse(baseUrl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID: tc.secretId,
			SecretKey: tc.secretKey,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
	opt := &cos.ImageRecognitionOptions{
		// 审核类型 porn（涉黄识别）、terrorist（涉暴恐识别）、politics（涉政识别）、ads（广告识别）四种
		DetectType: "porn,terrorist,politics",
	}

	res, resp, err := c.CI.ImageRecognition(context.Background(), path, opt)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		if resp.StatusCode != 200 {
			return nil, errors.New("http status not 200")
		}
	}

	//if cos.IsNotFoundError(err) {
	//	fmt.Println("WARN: Resource is not existed")
	//} else if e, ok := cos.IsCOSError(err); ok {
	//	fmt.Printf("ERROR: Code: %v\n", e.Code)
	//	fmt.Printf("ERROR: Message: %v\n", e.Message)
	//	fmt.Printf("ERROR: Resource: %v\n", e.Resource)
	//	fmt.Printf("ERROR: RequestId: %v\n", e.RequestID)
	//	// ERROR
	//} else {
	//	fmt.Printf("ERROR: %v\n", err)
	//	// ERROR
	//}

	return res, nil
}

const (
	CDN_SECRET = "DjL77HnpevmDlNrR2ACvjn60N1"
	CDN_HOST   = "https://resource-1253904687.file.myqcloud.com"
)
func (tc *TencentCloud) GenCdnUrl(baseUrl string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}
	
	now := time.Now().Unix()
	sign := util.Md5String(fmt.Sprintf("%s%s%d", CDN_SECRET, u.Path, now))
	return fmt.Sprintf("%s%s?sign=%s&t=%d", CDN_HOST, u.Path, sign, now), nil
}
