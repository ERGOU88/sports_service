package cposting

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"sports_service/server/dao"
	"testing"
)

func init() {
	dao.Engine = dao.InitXorm("root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4", []string{"root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4"})
	dao.InitRedis("192.168.5.12:6378", "")
}

// 保存用户通知设置
func TestSanitizeHtml(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	svc := New(c)
	content, err := svc.SanitizeHtml("<body><p>这是一条测试记录</p><a href='www.baidu.com'>@百度</a><br/></body>")
	t.Logf("\n content:%s, err:%s", content, err)
}
