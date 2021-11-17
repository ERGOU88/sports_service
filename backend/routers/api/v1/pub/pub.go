package pub

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mposting"
	"sports_service/server/models/mvideo"
	"strconv"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/backend/controller/cpub"
	"sports_service/server/global/backend/log"
)

// 发布视频
func PubVideo(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	params := new(mvideo.VideoPublishParams)
	if err := c.BindJSON(params); err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cpub.New(c)
	syscode := svc.RecordPubVideoInfo(userId.(string), params)
	reply.Response(http.StatusOK, syscode)
}

// 发布帖子
func PubPost(c *gin.Context) {
	reply := errdef.New(c)
	userId, ok := c.Get(consts.USER_ID)
	if !ok {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	params := new(mposting.PostPublishParam)
	if err := c.BindJSON(params); err != nil {
		log.Log.Errorf("post_trace: post publish params err:%s, params:%+v", err, params)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	if params.SectionId <= 0 {
		// 默认为x友讨论区
		params.SectionId = 1
	}

	svc := cpub.New(c)
	code := svc.PublishPosting(userId.(string), params)
	reply.Response(http.StatusOK, code)
}

// 获取上传签名（腾讯云）
func UploadSign(c *gin.Context) {
	reply := errdef.New(c)
	biteRate, err := strconv.Atoi(c.Query("bite_rate"))
	if err != nil {
		biteRate = 0
	}

	svc := cpub.New(c)
	syscode, sign, taskId := svc.GetUploadSign(c.Query("user_id"), int64(biteRate))
	reply.Data["sign"] = sign
	reply.Data["task_id"] = taskId

	reply.Response(http.StatusOK, syscode)
}

func PubInformation(c *gin.Context) {
	reply := errdef.New(c)
	param := &models.Information{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cpub.New(c)
	reply.Response(http.StatusOK, svc.PubInformation(param))
}

func SectionInfo(c *gin.Context) {
	reply := errdef.New(c)
	svc := cpub.New(c)
	sectionType := c.DefaultQuery("section_type", "0")
	code, list := svc.GetHomepageSectionInfo(sectionType)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}
