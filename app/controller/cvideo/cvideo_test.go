package cvideo

import (
  "github.com/gin-gonic/gin"
  "net/http/httptest"
  "testing"
  "sports_service/server/dao"
)

func init() {
  dao.Engine = dao.InitXorm("root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4", []string{"root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4"})
}

func TestUserBrowseVideosRecord(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  list := svc.UserBrowseVideosRecord("202009101933004667", 1, 10)
  t.Logf("list:%v\n", list)
}

func BenchmarkBrowseVideosRecord(b *testing.B) {
  for i := 0; i < b.N; i++ {
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    svc.UserBrowseVideosRecord("202009101933004667", 1, 10)
  }
}

func TestGetUserPublishList(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  list := svc.GetUserPublishList("202010101545291936", "-1", "-1", 1, 10)
  t.Logf("list:%v\n", list)
}

func BenchmarkGetUserPublishList(b *testing.B) {
  for i := 0; i < b.N; i++ {
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    svc.GetUserPublishList("202010101545291936", "-1", "-1", 1, 10)
  }
}
