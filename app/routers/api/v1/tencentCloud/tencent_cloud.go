package tencentCloud

import (
  "github.com/gin-gonic/gin"
  "net/http"
  "sports_service/server/global/app/log"
  "sports_service/server/global/consts"
  "sports_service/server/global/app/errdef"
  "sports_service/server/tools/tencentCloud"
)

// 获取腾讯cos临时通行证
func CosTempAccess(c *gin.Context) {
  reply := errdef.New(c)
  client := tencentCloud.New(consts.TX_CLOUD_COS_SECRET_ID, consts.TX_CLOUD_COS_SECRET_KEY, consts.TMS_API_DOMAIN)
  uploadType := c.DefaultQuery("upload_type", "private")
  info, err := client.GetCosTempAccess("ap-shanghai", uploadType)
  if err != nil {
    log.Log.Errorf("cloud_trace: get cos temp access err:%s", err)
    reply.Response(http.StatusOK, errdef.CLOUD_COS_ACCESS_FAIL)
    return
  }

  reply.Data["access_info"] = info
  reply.Response(http.StatusOK, errdef.SUCCESS)
}

type ValidateParam struct {
  Text     string     `json:"text"`
}
// 校验文本
func ValidateText(c *gin.Context) {
  reply := errdef.New(c)
  param := new(ValidateParam)
  if err := c.BindJSON(param); err != nil {
    log.Log.Errorf("cloud_trace: invalid param, err:%s, param:%+v", err, param)
    reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
    return
  }

  client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
  isPass, _,  err := client.TextModeration(param.Text)
  if !isPass || err != nil {
  	log.Log.Errorf("cloud_trace: invalid param err: %s，pass: %v", err, isPass)
  	reply.Response(http.StatusOK, errdef.COMMENT_INVALID_CONTENT)
  	return
  }

  reply.Response(http.StatusOK, errdef.ERROR)

}
