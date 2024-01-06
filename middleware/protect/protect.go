package protect

import (
	"github.com/gin-gonic/gin"
	"runtime"
	"sports_service/global/app/log"
)

// 捕获错误
func ProtectRun(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				log.Log.Errorf("RecoverRuntimeErr:%v", err)
				panic(err)
			default:
				log.Log.Errorf("RecoverCustomErr:%v", err)
			}
		}
	}()

	c.Next()
}
