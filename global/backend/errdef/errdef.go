package errdef

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Gin struct {
	Context *gin.Context
	Data    map[string]interface{}
}

func New(c *gin.Context) *Gin {
	return &Gin{
		Context: c,
		Data:    make(map[string]interface{}, 0),
	}
}

func (g *Gin) Response(httpCode, errCode int) {
	var data interface{}
	data = g.Data
	g.Context.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  GetMsg(errCode),
		"data": data,
		"tm":   time.Now().Unix(),
	})

	return
}
