package tencentCloud

import (
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"time"
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
