package cuser

import (
  "github.com/gin-gonic/gin"
  "net/http/httptest"
  "sports_service/server/models/muser"
  "testing"
  "sports_service/server/dao"
  "fmt"
)

func init() {
  dao.Engine = dao.InitXorm("root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4", []string{"root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4"})
}

func TestGetUserInfoByUserid(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  syscode, list := svc.GetUserInfoByUserid("202009101933004667")
  t.Logf("syscode:%d, list:%v\n", syscode, list)
}

func BenchmarkGetUserInfoByUserid(b *testing.B) {
  for i := 0; i < b.N; i++{
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    svc.GetUserInfoByUserid("202009101933004667")
  }
}

func TestEditUserInfo(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  params := &muser.EditUserInfoParams{
    Avatar:    1,
    NickName:  "陈二go",
    Born:      "1993-06-20",
    Gender:    1,
    CountryId: 1,
    Signature: "菩提本无树，明镜亦非台",
  }

  code := svc.EditUserInfo("202009101933004667", params)
  t.Logf("syscode:%d", code)
}

func BenchmarkEditUserInfo(b *testing.B) {
  num := 0
  for i := 0; i < b.N; i++ {
    num++
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    params := &muser.EditUserInfoParams{
      Avatar:    1,
      NickName:  fmt.Sprintf("陈二g.ou%d", num),
      Born:      "1993-06-20",
      Gender:    1,
      CountryId: 1,
      Signature: "菩提本无树，明镜亦非台",
    }

    code := svc.EditUserInfo("202009101933004667", params)
    b.Logf("syscode:%d\n", code)
  }
}
