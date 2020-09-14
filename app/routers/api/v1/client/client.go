package client

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/global/app/errdef"
	"sports_service/server/util"
)

func InitInfo(c *gin.Context) {
	// 生成secret
	secret := util.GenSecret(util.MIX_MODE, 16)
	reply := errdef.New(c)
	reply.Data["secret"] = secret
	reply.Response(http.StatusOK, errdef.SUCCESS)
}
