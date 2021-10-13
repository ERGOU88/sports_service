package contest

import "github.com/gin-gonic/gin"

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	contest := api.Group("/contest")
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
		// 设置赛事积分排行
		contest.POST("/set/integral/ranking", SetIntegralRanking)
		// 编辑赛事积分排行
		contest.POST("/edit/integral/ranking", EditIntegralRanking)
		// 赛事积分排行列表
		contest.GET("/integral/ranking/list", IntegralRankingList)
	}
}
