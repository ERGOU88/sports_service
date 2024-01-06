package contest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/app/controller/contest"
	"sports_service/global/app/errdef"
	"sports_service/util"
)

func BannerList(c *gin.Context) {
	reply := errdef.New(c)
	svc := contest.New(c)
	list := svc.GetBanner()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

func LiveList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	ts := c.Query("ts")
	pullType := c.DefaultQuery("pull_type", "1")
	svc := contest.New(c)
	code, list, pullUpTm, pullDownTm := svc.GetLiveList("", pullType, ts, page, size)
	reply.Data["list"] = list
	reply.Data["pull_up_tm"] = pullUpTm
	reply.Data["pull_down_tm"] = pullDownTm
	reply.Response(http.StatusOK, code)
}

// 首页推荐默认取2条
func RecommendLive(c *gin.Context) {
	reply := errdef.New(c)
	svc := contest.New(c)
	code, list, _, _ := svc.GetLiveList("1", "", "", 1, 2)
	reply.Data["list"] = list
	reply.Data["count"] = svc.GetLiveCount()
	reply.Response(http.StatusOK, code)
}

func ScheduleInfo(c *gin.Context) {
	reply := errdef.New(c)
	svc := contest.New(c)
	code, detail := svc.GetScheduleInfo()
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, code)
}

// 晋级信息
func PromotionInfo(c *gin.Context) {
	reply := errdef.New(c)
	contestId := c.Query("contest_id")
	scheduleId := c.Query("schedule_id")
	svc := contest.New(c)
	code, list := svc.GetPromotionInfo(contestId, scheduleId)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func IntegralRanking(c *gin.Context) {
	reply := errdef.New(c)
	contestId := c.Query("contest_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := contest.New(c)
	code, list := svc.GetIntegralRanking(contestId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func GetContestSection(c *gin.Context) {
	reply := errdef.New(c)
	svc := contest.New(c)
	code, sectionId := svc.GetContestSection()
	reply.Data["section_id"] = sectionId
	reply.Response(http.StatusOK, code)
}

func LiveScheduleData(c *gin.Context) {
	reply := errdef.New(c)
	svc := contest.New(c)
	liveId := c.Query("live_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	code, list := svc.GetLiveScheduleData(liveId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}
