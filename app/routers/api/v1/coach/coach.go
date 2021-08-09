package coach

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/coach"
	"sports_service/server/global/app/errdef"
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
