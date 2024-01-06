package venue

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/backend/controller/cvenue"
	"sports_service/global/backend/errdef"
	"sports_service/models"
	"sports_service/models/mvenue"
	"sports_service/util"
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
	param := &models.VenueRefundRules{}
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
	code := svc.AddVenueInfo(param)
	reply.Data["venue_id"] = param.Id
	reply.Response(http.StatusOK, code)

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

func MarkList(c *gin.Context) {
	reply := errdef.New(c)
	venueId := c.Query("venue_id")
	svc := cvenue.New(c)
	code, list := svc.MarkList(venueId)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}

func AddStoreManager(c *gin.Context) {
	reply := errdef.New(c)
	param := &mvenue.VenueAdminParam{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvenue.New(c)
	reply.Response(http.StatusOK, svc.AddStoreManager(param))
}

func EditStoreManager(c *gin.Context) {
	reply := errdef.New(c)
	param := &models.VenueAdministrator{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	svc := cvenue.New(c)
	reply.Response(http.StatusOK, svc.EditStoreManage(param))
}

func StoreManagerList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := cvenue.New(c)
	code, list := svc.StoreManageList(page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, code)
}
