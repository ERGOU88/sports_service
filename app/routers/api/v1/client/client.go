package client

import (
	"github.com/gin-gonic/gin"
	"net/http"
  "sports_service/server/app/config"
	_ "sports_service/server/app/routers/api/v1/swag"
  "sports_service/server/app/controller/cuser"
  "sports_service/server/app/controller/cvideo"
	"sports_service/server/global/app/errdef"
	"sports_service/server/util"
	"fmt"
)

// @Summary 初始化接口 (ok)
// @Tags 通用接口
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId" default(android)
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	    header    string 	true  "版本" default(1.0.0)
// @Success 200 {object}  swag.ClientInitSwag
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/client/init [get]
func InitInfo(c *gin.Context) {
	// 生成secret
	secret := util.GenSecret(util.MIX_MODE, 16)
	reply := errdef.New(c)
	svc := cuser.New(c)
	// 系统头像配置列表
	avatarList := svc.GetDefaultAvatarList()
	// 世界配置信息（暂时仅有国家信息）
	worldList := svc.GetWorldInfo()

	svc2 := cvideo.New(c)
	// 视频标签配置
  labelList := svc2.GetVideoLabelList()
	reply.Data["secret"] = secret
	reply.Data["avatar_list"] = avatarList
  reply.Data["world_list"] = worldList
  reply.Data["label_list"] = labelList
  // 登陆协议
  reply.Data["login_treaty"] = fmt.Sprintf("%s%s", config.Global.StaticDomain, "/static/template/login_treaty.html")
  // 上传协议
  reply.Data["upload_treaty"] = fmt.Sprintf("%s%s", config.Global.StaticDomain, "/static/template/upload_treaty.html")
  // 常见问题
  reply.Data["faq_h5"] = fmt.Sprintf("%s%s", config.Global.StaticDomain, "/static/template/faq.html")
  // 关于
  reply.Data["about"] = fmt.Sprintf("%s%s", config.Global.StaticDomain, "/static/template/about_content.html")
  reply.Response(http.StatusOK, errdef.SUCCESS)
}
