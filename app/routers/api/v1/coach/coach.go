package coach

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/coach"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mcoach"
	"sports_service/server/util"
)

func CoachList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := coach.New(c)
	syscode, list := svc.GetCoachList(page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}

func CoachDetail(c *gin.Context) {
	reply := errdef.New(c)
	coachId := c.Query("coach_id")
	if coachId == "" {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := coach.New(c)
	code, detail := svc.GetCoachDetail(coachId)
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, code)
}

// 私教评价列表
func CoachEvaluate(c *gin.Context) {
	reply := errdef.New(c)
	coachId := c.Query("coach_id")
	if coachId == "" {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := coach.New(c)
	code, list := svc.GetEvaluateList(coachId, page, size)
	reply.Data["list"] = list
	total, score := svc.GetCoachScore(coachId)
	reply.Data["score"] = score
	reply.Data["total_num"] = total
	reply.Response(http.StatusOK, code)
}


func CoachEvaluateConf(c *gin.Context) {
	reply := errdef.New(c)
	svc := coach.New(c)
	code, list := svc.GetEvaluateConfig()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

// 发布评价 [私教]
func PubEvaluate(c *gin.Context) {
	reply := errdef.New(c)
	param := &mcoach.PubEvaluateParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("coach_trace: invalid param, param:%+v, err:%s", param, err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	userId, _ := c.Get(consts.USER_ID)
	svc := coach.New(c)
	code := svc.PubEvaluate(userId.(string), param)
	reply.Response(http.StatusOK, code)
}
