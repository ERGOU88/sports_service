package clike

import (
  "github.com/gin-gonic/gin"
  "net/http/httptest"
  "sports_service/server/dao"
  "testing"
)

func init() {
  dao.AppEngine = dao.InitXorm("root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4", []string{"root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4"})
}

// 视频点赞
func TestGiveLikeForVideo(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  syscode := svc.GiveLikeForVideo("202009101933004667", 59)
  t.Logf("\n syscode:%v", syscode)
}

func BenchmarkLikeForVideo(b *testing.B) {
  for i := 0; i < b.N; i++ {
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    svc.GiveLikeForVideo("202009101933004667", 59)
  }
}

func TestGetUserLikeVideos(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  list := svc.GetUserLikeVideos("202009101933004667", 1, 10)
  t.Logf("\n list len: %d", len(list))
}

func BenchmarkUserLikeVideos(b *testing.B) {
  for i := 0; i < b.N; i ++ {
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    list := svc.GetUserLikeVideos("202009101933004667", 1, 10)
    b.Logf("\n list len: %d", len(list))
  }
}

// 视频取消点赞
func TestCancelLikeForVideo(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  syscode := svc.CancelLikeForVideo("202009101933004667", 59)
  t.Logf("\n syscode:%d", syscode)
}

func BenchmarkCancelLikeForVideo(b *testing.B) {
  for i := 0; i < b.N; i++ {
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    svc.CancelLikeForVideo("202009101933004667", 59)
  }
}

// 给评论点赞
func TestGiveLikeForComment(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  syscode := svc.GiveLikeForVideoComment("202009101933004667", 43)
  t.Logf("\n syscode:%v", syscode)
}

func BenchmarkGiveLikeForComment(b *testing.B) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  syscode := svc.GiveLikeForVideoComment("202009101933004667", 43)
  b.Logf("\n syscode:%d", syscode)
}

// 取消评论点赞
func TestCancelLikeForComment(t *testing.T) {
 c, _ := gin.CreateTestContext(httptest.NewRecorder())
 svc := New(c)
 svc.CancelLikeForVideoComment("202009101933004667", 43)
}

func BenchmarkCancelLikeForComment(b *testing.B) {
 c, _ := gin.CreateTestContext(httptest.NewRecorder())
 svc := New(c)
 svc.CancelLikeForVideoComment("202009101933004667", 43)
}






