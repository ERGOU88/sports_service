package appointment

import (
	"github.com/gin-gonic/gin"
	"sports_service/server/middleware/sign"
)

// 预约模块路由
func Router(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	appointment := api.Group("/appointment")
	appointment.Use(sign.CheckSign())
	{
		appointment.GET("/date", AppointmentDate)

	}
}
