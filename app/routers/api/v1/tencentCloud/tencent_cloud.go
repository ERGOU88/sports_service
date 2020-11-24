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
  client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
  info, err := client.GetCosTempAccess("ap-shanghai")
  if err != nil {
    log.Log.Errorf("cloud_trace: get cos temp access err:%s", err)
    reply.Response(http.StatusOK, errdef.CLOUD_COS_ACCESS_FAIL)
    return
  }

  reply.Data["access_info"] = info
  reply.Response(http.StatusOK, errdef.SUCCESS)
}
