package cappointment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"testing"
	"time"
	"github.com/go-sql-driver/mysql"
)

func init() {
	dao.VenueEngine = dao.InitXorm("root:bluetrans888@tcp(192.168.5.12:3306)/venue?charset=utf8mb4", []string{"root:bluetrans888@tcp(192.168.5.12:3306)/venue?charset=utf8mb4"})
	dao.AppEngine = dao.InitXorm("root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4", []string{"root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4"})
	dao.InitRedis("192.168.5.12:6378", "")
}

//func TestAddAttention(t *testing.T) {
//	c, _ := gin.CreateTestContext(httptest.NewRecorder())
//	code, list := GetAppointmentDate(NewVenue(c))
//	t.Logf("code:%d, list:%+v", code, list)
//}

// 添加库存
func TestAddStock(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	svc := NewVenue(c)
	now := int(time.Now().Unix())
	data := make([]*models.VenueAppointmentStock, 1)
	info := &models.VenueAppointmentStock{
		Date: "2021-08-10",
		TimeNode: "18:00-19:00",
		QuotaNum: 5,
		PurchasedNum: 2,
		AppointmentType: 0,
		VenueId: 1,
		CreateAt: now,
		UpdateAt: now,
	}

	//
	affected, err := svc.appointment.AddStockInfo(info)
	if err != nil && affected == 0 {
		log.Log.Errorf("")
	}

	var myerr *mysql.MySQLError
	t.Logf("err:%s, affected:%d, now:%d, b:%v, %#v, num:%d", err, affected, now, errors.As(err, &myerr), err, myerr.Number)
}


func TestUpdateStock(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	svc := NewVenue(c)
	affected, err := svc.appointment.UpdateVenueStockInfo("18:00-19:00", "2021-08-10", 1, 1,
		0, 1)
	var myerr *mysql.MySQLError
	t.Logf("affected:%d, err:%s, bool:%v", affected, err, errors.As(err, &myerr))
}
