package cappointment

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"sports_service/server/dao"
	"testing"
)

func init() {
	dao.AppEngine = dao.InitXorm("root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4", []string{"root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4"})
	dao.InitRedis("192.168.5.12:6378", "")
}

func TestAddAttention(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	code, list := GetAppointmentDate(NewVenue(c))
	t.Logf("code:%d, list:%+v", code, list)
}

