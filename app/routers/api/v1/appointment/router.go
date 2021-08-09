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
		// 预约日期选项
		appointment.GET("/date", AppointmentDate)
		// 预约时间选项
		appointment.GET("/time/options", AppointmentTimeOptions)
		// 预约选项 [场馆、私课、大课选项]
		appointment.GET("/options", AppointmentOptions)
	}
}
