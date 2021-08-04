package cappointment

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"sports_service/server/dao"
	"testing"
)

func init() {
	dao.Engine = dao.InitXorm("root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4", []string{"root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4"})
}

func TestAddAttention(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	code, list := GetAppointmentDate(NewCoach(c))
	t.Logf("code:%d, list:%+v", code, list)
}

