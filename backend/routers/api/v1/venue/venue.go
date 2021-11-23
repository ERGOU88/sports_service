package venue

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/backend/controller/cvenue"
	"sports_service/server/models"
	"sports_service/server/models/morder"
	"sports_service/server/models/mvenue"
)

func VenueList(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvenue.New(c)
	code, list := svc.GetVenueList()
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func VenueDetail(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvenue.New(c)
	info, err := svc.GetVenueInfo(c.Query("id"))
	if err != nil {
		reply.Response(http.StatusOK, errdef.ERROR)
		return
	}

	reply.Data["detail"] = info
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

func EditVenue(c *gin.Context) {
	reply := errdef.New(c)
	param := &models.VenueInfo{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvenue.New(c)
	reply.Response(http.StatusOK, svc.EditVenueInfo(param))
}

func UpdateRefundRate(c *gin.Context) {
	reply := errdef.New(c)
	param := &morder.RefundRateParam{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvenue.New(c)
	reply.Response(http.StatusOK, svc.UpdateRefundRate(param))
}

func RefundRules(c *gin.Context) {
	reply := errdef.New(c)
	svc := cvenue.New(c)
	code, rules := svc.GetRefundRules()
	reply.Data["list"] = rules
	reply.Response(http.StatusOK, code)
}

func AddVenue(c *gin.Context) {
	reply := errdef.New(c)
	param := &models.VenueInfo{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvenue.New(c)
	reply.Data["venue_id"] = param.Id
	reply.Response(http.StatusOK, svc.AddVenueInfo(param))

}

func AddMark(c *gin.Context) {
	reply := errdef.New(c)
	param := &mvenue.AddMarkParam{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cvenue.New(c)
	reply.Response(http.StatusOK, svc.AddMark(param))
}

func DelMark(c *gin.Context) {
	reply := errdef.New(c)
	param := &mvenue.DelMarkParam{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cvenue.New(c)
	reply.Response(http.StatusOK, svc.DelMark(param))
}
