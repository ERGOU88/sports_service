package search

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/app/controller/csearch"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/util"
	_ "sports_service/server/app/routers/api/v1/swag"
	_ "sports_service/server/models/muser"
	_ "sports_service/server/models/mvideo"
	_ "sports_service/server/models/mattention"
)

// @Summary 视频搜索[分页获取] (ok)
// @Tags 搜索模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Param   name	  	  query  	string 	true  "搜索的名称"
// @Param   sort     	  query  	string 	true  "排序条件  0 播放数 1 弹幕数 2 点赞数 默认 播放数"
// @Param   duration   	  query  	string 	true  "视频时长 0 表示没有限制 1 表示 1～5分钟  2：5～10分钟 3：10～30分钟 4：30分钟以上"
// @Param   publish_time  query  	string 	true  "视频发布时间 0 不限制 1 一天内 2 一周内 3 半年内"
// @Param   user_id	 	    query  	string 	true  "用户id"
// @Success 200 {array}  mvideo.VideoDetailInfo
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/search/videos [get]
// 视频搜索
func VideoSearch(c *gin.Context) {
	reply := errdef.New(c)
	//userId, _ := c.Get(consts.USER_ID)
	userId := c.Query("user_id")
	name := c.Query("name")
	sort := c.Query("sort")
	duration := c.Query("duration")
	publishTime := c.Query("publish_time")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := csearch.New(c)
	// 搜索视频
	list := svc.VideoSearch(userId, name, sort, duration, publishTime, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 用户搜索[分页获取] (ok)
// @Tags 搜索模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Param   name	  	  query  	string 	true  "搜索的名称/userId"
// @Param   user_id	 	  query  	string 	true  "用户id"
// @Success 200 {array}  muser.UserSearchResults
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/search/users [get]
// 用户搜索
func UserSearch(c *gin.Context) {
	reply := errdef.New(c)
	//userId, _ := c.Get(consts.USER_ID)
	userId := c.Query("user_id")
	name := c.Query("name")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := csearch.New(c)
	list := svc.UserSearch(userId, name, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 综合搜索[默认搜视频、用户 各取三条记录] (ok)
// @Tags 搜索模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   name	  	  query  	string 	true  "搜索的名称"
// @Param   user_id	 	  query  	string 	true  "用户id"
// @Success 200 {object}  swag.ColligateSearchSwag
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/search/colligate [get]
// 综合搜索（视频 + 用户）
func ColligateSearch(c *gin.Context) {
	reply := errdef.New(c)
	// userId, _ := c.Get(consts.USER_ID)
	userId := c.Query("user_id")
	name := c.Query("name")

	svc := csearch.New(c)
	// 综合搜索
	videoList, userList, recommendList := svc.ColligateSearch(userId, name)
	reply.Data["video_list"] = videoList
	reply.Data["user_list"] = userList
	reply.Data["recommend_list"] = recommendList
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 标签搜索视频[分页获取] (ok)
// @Tags 搜索模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Param   label_id	  query  	string 	true  "搜索的视频标签"
// @Param   user_id	 	  query  	string 	true  "用户id"
// @Success 200 {array}  mvideo.VideoDetailInfo
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/search/label [get]
// 标签搜索视频
func LabelSearch(c *gin.Context) {
	reply := errdef.New(c)
	//userId, _ := c.Get(consts.USER_ID)
	userId := c.Query("user_id")
	labelId := c.Query("label_id")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := csearch.New(c)
	// 标签搜索视频
	list := svc.LabelSearch(userId, labelId, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 热门搜索 (ok)
// @Tags 搜索模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	    header    string 	true  "版本" default(1.0.0)"
// @Param   user_id 	    query     string 	true  "用户id"
// @Success 200 {object}  []string
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/search/hot [get]
// 热门搜索
func HotSearch(c *gin.Context) {
	reply := errdef.New(c)
	svc := csearch.New(c)
  userId := c.Query("user_id")
	// 获取后台配置的热门搜索内容
	hotSearch := svc.GetHotSearch()
	reply.Data["list"] = hotSearch
	reply.Data["history_list"] = svc.GetHistorySearch(userId)
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 搜索关注的用户[分页获取] (ok)
// @Tags 搜索模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Param   name    	  query  	string 	true  "搜索的用户名/userid"
// @Success 200 {array}  mattention.SearchContactRes
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/search/attention [get]
// 搜索关注的用户
func AttentionSearch(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	name := c.Query("name")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := csearch.New(c)
	list := svc.SearchAttentionUser(userId.(string), name, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}

// @Summary 搜索粉丝[分页获取] (ok)
// @Tags 搜索模块
// @Version 1.0
// @Description
// @Accept json
// @Produce  json
// @Param   AppId         header    string 	true  "AppId"
// @Param   Secret        header    string 	true  "调用/api/v1/client/init接口 服务端下发的secret"
// @Param   Timestamp     header    string 	true  "请求时间戳 单位：秒"
// @Param   Sign          header    string 	true  "签名 md5签名32位值"
// @Param   Version 	  header    string 	true  "版本" default(1.0.0)
// @Param   page	  	  query  	string 	true  "页码 从1开始"
// @Param   size	  	  query  	string 	true  "每页展示多少 最多50条"
// @Param   name    	  query  	string 	true  "搜索的用户名/userid"
// @Success 200 {array}  mattention.SearchContactRes
// @Failure 500 {string} json "{"code":500,"data":{},"msg":"fail","tm":"1588888888"}"
// @Router /api/v1/search/fans [get]
// 搜索粉丝
func FansSearch(c *gin.Context) {
	reply := errdef.New(c)
	userId, _ := c.Get(consts.USER_ID)
	name := c.Query("name")
	page, size := util.PageInfo(c.Query("page"), c.Query("size"))

	svc := csearch.New(c)
	list := svc.SearchFans(userId.(string), name, page, size)
	reply.Data["list"] = list
	reply.Response(http.StatusOK, errdef.SUCCESS)
}
