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
	queryType, err := strconv.Atoi(c.DefaultQuery("query_type", "0"))
	if err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	factory := &cappointment.AppointmentFactory{}
	i = factory.Create(queryType, c)

	syscode, list := cappointment.GetAppointmentDate(i)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}

// 预约时间选项
func AppointmentTimeOptions(c *gin.Context) {
	reply := errdef.New(c)
	queryType, err := strconv.Atoi(c.DefaultQuery("query_type", "0"))
	if err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

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

	// 日期id
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}


	var i cappointment.IAppointment
	factory := &cappointment.AppointmentFactory{}
	i = factory.Create(queryType, c)
	i.SetWeek(week)
	i.SetAppointmentType(queryType)
	i.SetRelatedId(relatedId)
	i.SetDateId(id)
	syscode, list := cappointment.GetAppointmentTimeOptions(i)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}

// 预约选项
func AppointmentOptions(c *gin.Context) {
	reply := errdef.New(c)
	queryType, err := strconv.Atoi(c.DefaultQuery("query_type", "0"))
	if err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	relatedId, err := strconv.Atoi(c.Query("related_id"))
	if err != nil || relatedId <= 0 {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	var i cappointment.IAppointment
	factory := &cappointment.AppointmentFactory{}
	i = factory.Create(queryType, c)
	syscode, list := cappointment.GetOptions(i, int64(relatedId))
	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}
