package stat

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/backend/controller/cstat"
	"sports_service/server/global/backend/errdef"
	"net/http"
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

}
