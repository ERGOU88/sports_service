package appointment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cappointment"
	"sports_service/server/global/app/errdef"
)

func AppointmentDate(c *gin.Context) {
	reply := errdef.New(c)
	syscode, list := cappointment.GetAppointmentDate(cappointment.NewVenue(c))
	reply.Data["list"] = list
	reply.Response(http.StatusOK, syscode)
}
