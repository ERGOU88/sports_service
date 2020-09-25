package configure

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/backend/controller/configure"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models/mbanner"
	"sports_service/server/util"
)

// 添加banner
func AddBanner(c *gin.Context) {
	reply := errdef.New(c)
	params := new(mbanner.AddBannerParams)
	if err := c.BindJSON(params); err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := configure.New(c)
	syscode := svc.AddBanner(params)
	reply.Response(http.StatusOK, syscode)
}

// 删除banner
func DelBanner(c *gin.Context) {
	reply := errdef.New(c)
	param := new(mbanner.DelBannerParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := configure.New(c)
	syscode := svc.DelBanner(param)
	reply.Response(http.StatusOK, syscode)
}

// 获取banner列表
func GetBanners(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := configure.New(c)
	list := svc.GetBannerList(page, size)

	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}


