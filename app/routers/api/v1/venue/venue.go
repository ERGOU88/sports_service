package venue

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/app/controller/cvenue"
	"sports_service/global/app/errdef"
	"sports_service/global/app/log"
	"sports_service/global/consts"
	"sports_service/models/mvenue"
	"sports_service/util"
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

// 购买次卡/月卡/年卡
func PurchaseVipCard(c *gin.Context) {
	reply := errdef.New(c)
	param := &mvenue.PurchaseVipCardParam{}
	if err := c.BindJSON(param); err != nil {
		log.Log.Errorf("venue_trace: invalid param, err:%s", err)
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	userId, _ := c.Get(consts.USER_ID)
	svc := cvenue.New(c)
	channel, _ := c.Get(consts.CHANNEL)
	param.UserId = userId.(string)
	param.ChannelId = channel.(int)
	code, rsp := svc.PurchaseVipCard(param)
	if code == errdef.SUCCESS {
		reply.Data["info"] = rsp
	}

	reply.Response(http.StatusOK, code)
}

// 进出场记录
func ActionRecord(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	userId, _ := c.Get(consts.USER_ID)
	svc := cvenue.New(c)
	code, list := svc.GetActionRecord(userId.(string), page, size)
	if code == errdef.SUCCESS {
		reply.Data["list"] = list
	}

	reply.Response(http.StatusOK, code)
}
