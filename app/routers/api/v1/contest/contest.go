package contest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/contest"
	"sports_service/server/global/app/errdef"
)

func BannerList(c *gin.Context) {
	reply := errdef.New(c)
	svc := contest.New(c)
	list := svc.GetBanner()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}
