package contest

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/jwt"
)

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	contest := api.Group("/contest")
	contest.Use(jwt.JwtAuth())
	{
		// 添加选手信息
		contest.POST("/add/player", AddPlayer)
		// 编辑选手信息
		contest.POST("/edit/player", EditPlayer)
		// 选手列表
		contest.GET("/player/list", PlayerList)
		// 添加赛程分组
		contest.POST("/add/group", AddContestGroup)
		// 编辑赛程分组
		contest.POST("/edit/group", EditContestGroup)
		// 赛程分组列表
		contest.GET("/group/list", ContestGroupList)
		// 获取赛程信息
		contest.GET("/schedule", ContestSchedule)
		// 添加赛程详情
		contest.POST("/add/schedule/detail", AddContestScheduleDetail)
		// 赛程详情列表
		contest.GET("/schedule/detail/list", ContestScheduleDetailList)
		// 删除赛程详情数据
		contest.DELETE("/del/schedule/detail", DelScheduleDetail)
		// 设置赛事积分排行
		contest.POST("/set/integral/ranking", SetIntegralRanking)
		// 编辑赛事积分排行
		contest.POST("/edit/integral/ranking", EditIntegralRanking)
		// 赛事积分排行列表
		contest.GET("/integral/ranking/list", IntegralRankingList)
		// 添加赛事直播
		contest.POST("/add/live", AddContestLive)
		// 更新赛事直播
		contest.POST("/update/live", UpdateContestLive)
		// 删除赛事直播
		contest.DELETE("/del/live", DelContestLive)
		// 赛事直播列表
		contest.GET("/live/list", ContestLiveList)

	}
}
