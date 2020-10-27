package cattention

import (
  "github.com/gin-gonic/gin"
  "net/http/httptest"
  "sports_service/server/dao"
  "testing"
)

func init() {
  dao.Engine = dao.InitXorm("root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4", []string{"root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4"})
}

// 添加关注
func TestAddAttention(t *testing.T) {
 c, _ := gin.CreateTestContext(httptest.NewRecorder())
 svc := New(c)
 syscode := svc.AddAttention("202009101933004667", "202010101545291936")
 t.Logf("syscode:%d", syscode)
}

func BenchmarkAddAttention(b *testing.B) {
 for i := 0; i < b.N; i++{
   c, _ := gin.CreateTestContext(httptest.NewRecorder())
   svc := New(c)
   svc.AddAttention("202009101933004667", "202010101545291936")
 }
}

// 关注的用户列表
func TestGetAttentionUserList(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  list := svc.GetAttentionUserList("202009101933004667", 1, 10)
  t.Logf("\n attention list len :%d", len(list))
}

func BenchmarkAttentionUserList(b *testing.B) {
  for i := 0; i < b.N; i++{
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    list := svc.GetAttentionUserList("202009101933004667", 1, 10)
    b.Logf("\n attention list len :%d", len(list))
  }
}

// 粉丝列表
func TestGetFansList(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  list := svc.GetFansList("202010101545291936", 1, 10)
  t.Logf("\n fans list len :%d", len(list))
}

func BenchmarkFansList(b *testing.B) {
  for i := 0; i < b.N; i++ {
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    list := svc.GetFansList("202010101545291936", 1, 10)
    b.Logf("\n fans list len :%d", len(list))
  }
}

// 取消关注
func TestCancelAttention(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  svc.CancelAttention("202009101933004667", "202010101545291936")
}

