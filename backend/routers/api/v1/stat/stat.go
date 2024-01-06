package stat

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/backend/controller/cstat"
	"sports_service/global/backend/errdef"
)

func HomePageInfo(c *gin.Context) {
	reply := errdef.New(c)
	svc := cstat.New(c)
	minDate := c.Query("min_date")
	maxDate := c.Query("max_date")
	code, detail := svc.GetHomePageInfo(minDate, maxDate)
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, code)
}

func EcologicalInfo(c *gin.Context) {
	reply := errdef.New(c)
	svc := cstat.New(c)
	minDate := c.Query("min_date")
	maxDate := c.Query("max_date")
	code, videoStat := svc.GetVideoSubareaStat()
	if code != errdef.SUCCESS {
		reply.Response(http.StatusOK, code)
		return
	}

	code, postStat := svc.GetPostSectionStat()
	if code != errdef.SUCCESS {
		reply.Response(http.StatusOK, code)
		return
	}

	code, pubVideoStat := svc.PublishVideoDaily(minDate, maxDate)
	if code != errdef.SUCCESS {
		reply.Response(http.StatusOK, code)
		return
	}

	code, pubPostStat := svc.PublishPostDaily(minDate, maxDate)
	if code != errdef.SUCCESS {
		reply.Response(http.StatusOK, code)
		return
	}

	code, dailyPostStat := svc.DailyTotalPost(minDate, maxDate)
	if code != errdef.SUCCESS {
		reply.Response(http.StatusOK, code)
		return
	}

	code, dailyVideoStat := svc.DailyTotalVideo(minDate, maxDate)
	if code != errdef.SUCCESS {
		reply.Response(http.StatusOK, code)
		return
	}

	reply.Data["daily_total_video"] = dailyVideoStat
	reply.Data["daily_total_post"] = dailyPostStat
	reply.Data["video_stat"] = videoStat
	reply.Data["post_stat"] = postStat
	reply.Data["pub_video_stat"] = pubVideoStat
	reply.Data["pub_post_stat"] = pubPostStat

	reply.Response(http.StatusOK, errdef.SUCCESS)
}
