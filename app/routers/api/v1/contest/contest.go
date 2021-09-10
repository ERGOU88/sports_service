package contest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/contest"
	"sports_service/server/global/app/errdef"
	"sports_service/server/util"
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
	svc := contest.New(c)
	code, list := svc.GetLiveList("", page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

// 首页推荐默认取两条
func RecommendLive(c *gin.Context) {
	reply := errdef.New(c)
	svc := contest.New(c)
	code, list := svc.GetLiveList("1",1, 2)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func ScheduleInfo(c *gin.Context) {
	reply := errdef.New(c)
	svc := contest.New(c)
	code, detail := svc.GetScheduleInfo()
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, code)
}
