package appointment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cappointment"
	"sports_service/server/global/app/errdef"
	"strconv"
)

// 预约日期信息
func AppointmentDate(c *gin.Context) {
	reply := errdef.New(c)
	var i cappointment.IAppointment
	queryType := c.DefaultQuery("query_type", "0")

	switch queryType {
	case "0":
		i = cappointment.NewVenue(c)
	case "1":
		i = cappointment.NewCoach(c)
	default:
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
	}

	syscode, list := cappointment.GetAppointmentDate(i)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}

// 预约时间选项
func AppointmentOptions(c *gin.Context) {
	reply := errdef.New(c)
	var i cappointment.IAppointment
	queryType := c.DefaultQuery("query_type", "0")
	week, err := strconv.Atoi(c.Query("week"))
	if err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	relatedId, err := strconv.Atoi(c.Query("related_id"))
	if err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	if week > 6 || week < 0 {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	switch queryType {
	case "0":
		svc := cappointment.NewVenue(c)
		svc.RelatedId = relatedId
		svc.WeekNum = week
		svc.AppointmentType = 0
		i = svc
	case "1":
		svc := cappointment.NewCoach(c)
		svc.RelatedId = relatedId
		svc.WeekNum = week
		svc.AppointmentType = 1
		i = svc

	default:
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	syscode, list := cappointment.GetAppointmentOptions(i)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}
