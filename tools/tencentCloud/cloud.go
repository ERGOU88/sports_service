package tencentCloud

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
