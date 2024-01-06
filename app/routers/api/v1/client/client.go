package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/app/config"
	"sports_service/app/controller/cuser"
	"sports_service/app/controller/cvideo"
	_ "sports_service/app/routers/api/v1/swag"
	"sports_service/global/app/errdef"
	"sports_service/util"
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
// @Success 200 {string} json "{"code":200,"data":{},"msg":"success","tm":"1588888888"}"
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

	type H5Info struct {
		PrivacyTreaty string `json:"privacy_treaty"` // 隐私协议
		UserTreaty    string `json:"user_treaty"`    // 用户协议
		CommentReport string `json:"comment_report"` // 举报评论
		CommonProblem string `json:"common_problem"` // 常见问题
		AboutFpv      string `json:"about_fpv"`      // 关于fpv
		Feedback      string `json:"feedback"`       // 问题反馈
		NoticeDetail  string `json:"notice_detail"`  // 消息详情
	}

	h5Info := &H5Info{
		PrivacyTreaty: fmt.Sprintf("%s%s", config.Global.StaticDomain, "/privacyagreement"),
		UserTreaty:    fmt.Sprintf("%s%s", config.Global.StaticDomain, "/useragreement"),
		CommentReport: fmt.Sprintf("%s%s", config.Global.StaticDomain, "/commentreport"),
		Feedback:      fmt.Sprintf("%s%s", config.Global.StaticDomain, "/problemfeedback"),
		AboutFpv:      fmt.Sprintf("%s%s", config.Global.StaticDomain, "/about"),
		CommonProblem: fmt.Sprintf("%s%s", config.Global.StaticDomain, "/problemcommon"),
		NoticeDetail:  fmt.Sprintf("%s%s", config.Global.StaticDomain, "/noticedetail"),
	}

	reply.Data["h5_info"] = h5Info
	reply.Response(http.StatusOK, errdef.SUCCESS)
}
