package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/backend/controller/cuser"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	"sports_service/server/util"
)

// 获取用户列表
func UserList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	condition := c.Query("condition")
	sortType := c.Query("sort_type")
	queryId := c.Query("query_id")

	svc := cuser.New(c)
	//list := svc.GetUserList(page, size)
	list, total := svc.GetUserListBySort(queryId, sortType, condition, page, size)
	reply.Data["list"] = list
	reply.Data["total"] = total
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 封禁用户
func ForbidUser(c *gin.Context) {
	reply := errdef.New(c)
	param := new(muser.ForbidUserParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	syscode := svc.ForbidUser(param.Id)
	reply.Response(http.StatusOK, syscode)
}

// 解封用户
func UnForbidUser(c *gin.Context) {
	reply := errdef.New(c)
	param := new(muser.UnForbidUserParam)
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
		return
	}

	svc := cuser.New(c)
	syscode := svc.UnForbidUser(param.Id)
	reply.Response(http.StatusOK, syscode)
}

func OfficialUserList(c *gin.Context) {
	reply := errdef.New(c)
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))
	svc := cuser.New(c)
	code, list := svc.GetOfficialUsers(page, size)
	reply.Data["list"] = list
	reply.Data["total"] = svc.GetOfficialUserTotal()
	reply.Response(http.StatusOK, code)
}

func AddOfficialUser(c *gin.Context) {
	reply := errdef.New(c)
	param := &models.User{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}

	param.RegIp = c.ClientIP()
	svc := cuser.New(c)
	reply.Response(http.StatusOK, svc.AddOfficialUser(param))
}

func EditOfficialUser(c *gin.Context) {
	reply := errdef.New(c)
	param := &models.User{}
	if err := c.BindJSON(param); err != nil {
		reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
		return
	}
	
	svc := cuser.New(c)
	reply.Response(http.StatusOK, svc.AddOfficialUser(param))
}
