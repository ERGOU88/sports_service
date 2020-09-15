package client

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cuser"
	"sports_service/server/global/app/errdef"
	"sports_service/server/util"
)

// @Summary 初始化接口 (ok)
// @Tags 通用接口接口
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId" default(android)
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Success 200 {string} json "{"code":200,"data":{"secret": "DnaukFwVILpcewX6"},"msg":"success","tm":"1588888888"}"
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/client/init [post]
func InitInfo(c *gin.Context) {
	// 生成secret
	secret := util.GenSecret(util.MIX_MODE, 16)
	reply := errdef.New(c)
	svc := cuser.New(c)
	// 系统头像配置列表
	avatarList := svc.GetDefaultAvatarList()
	// 世界信息（暂时仅有国家信息）
	worldList := svc.GetWorldInfo()
	reply.Data["secret"] = secret
	reply.Data["avatarList"] = avatarList
	reply.Data["worldList"] = worldList
	reply.Response(http.StatusOK, errdef.SUCCESS)
}