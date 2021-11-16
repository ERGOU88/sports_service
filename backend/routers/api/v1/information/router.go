package information

import (
	"github.com/gin-gonic/gin"
)

func Router(engine *gin.Engine) {
	api := engine.Group("/backend/v1")
	information := api.Group("/information")
	//information.Use(jwt.JwtAuth())
	{
		information.GET("/list", InformationList)
		information.DELETE("/delete", DeleteInformation)
	}
}
