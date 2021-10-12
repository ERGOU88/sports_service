package finance

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sports_service/server/global/backend/errdef"
)

func OrderList(c *gin.Context) {
	reply := errdef.New(c)
	reply.Response(http.StatusOK, errdef.SUCCESS)
}
