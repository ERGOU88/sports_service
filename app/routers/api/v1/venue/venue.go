package venue

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/cvenue"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"strconv"
)

// 获取场馆信息
func VenueInfo(c *gin.Context) {
	reply := errdef.New(c)
	venueId, err := strconv.Atoi(c.DefaultQuery("venue_id", "1"))
	if err != nil {
		log.Log.Errorf("venue_trace: invalid param, venueId:%s", venueId)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cvenue.New(c)
	code, venue, products := svc.GetHomePageInfo(int64(venueId))
	reply.Data["venue"] = venue
	reply.Data["products"] = products
	reply.Response(http.StatusOK, code)
}
