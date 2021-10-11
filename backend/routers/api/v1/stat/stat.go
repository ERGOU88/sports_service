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
	code, detail := svc.GetHomePageInfo()
	reply.Data["detail"] = detail
	reply.Response(http.StatusOK, code)

}
