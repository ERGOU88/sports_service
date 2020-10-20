package admin

import (
	"github.com/gin-gonic/gin"
  "net/http"
  "sports_service/server/backend/controller/cadmin"
  "sports_service/server/global/backend/errdef"
  "sports_service/server/models/madmin"
)

// 注册后台用户
func RegAdminUser(c *gin.Context) {
  reply := errdef.New(c)
  params := new(madmin.AdminRegOrLoginParams)
  if err := c.BindJSON(params); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  syscode := svc.AddAdminUser(params)
  reply.Response(http.StatusOK, syscode)
}

// 后台管理员登陆
func LoginByPassword(c *gin.Context) {
  reply := errdef.New(c)
  params := new(madmin.AdminRegOrLoginParams)
  if err := c.BindJSON(params); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  syscode := svc.AdminLogin(params)
  reply.Response(http.StatusOK, syscode)
}
