package comment

import (
  "github.com/gin-gonic/gin"
  "net/http/httptest"
  "sports_service/server/models/mcomment"
  "testing"
  "sports_service/server/dao"
)

func init() {
  dao.Engine = dao.InitXorm("root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4", []string{"root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4"})
}

// 发布评论
func TestPublishComment(t *testing.T) {
  params := &mcomment.PublishCommentParams{
    VideoId: 97,
    Content: "我是1级评论",
  }
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  syscode := svc.PublishComment("202009101933004667", params)
  t.Logf("syscode:%d\n", syscode)
}

func BenchmarkPublishComment(b *testing.B) {
  for i := 0; i < b.N; i++ {
    params := &mcomment.PublishCommentParams{
      VideoId: 97,
      Content: "1级评论",
    }
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    syscode := svc.PublishComment("202009101933004667", params)
    b.Logf("syscode:%d\n", syscode)
  }
}

// 发布回复
func TestPublishReply(t *testing.T) {
  params := &mcomment.ReplyCommentParams{
    VideoId: 97,
    Content: "评论回复no.1",
    ReplyId: "57",
  }

  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  syscode := svc.PublishReply("202009181548217779", params)
  t.Logf("syscode:%d\n", syscode)
}

func BenchmarkPublishReply(b *testing.B) {
  for i := 0; i < b.N; i++ {
    params := &mcomment.ReplyCommentParams{
      VideoId: 97,
      Content: "评论回复",
      ReplyId: "57",
    }
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    syscode := svc.PublishReply("202009181548217779", params)
    b.Logf("syscode:%d\n", syscode)
  }
}

// 获取视频评论列表
func TestGetVideoComments(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  list := svc.GetVideoComments("202009181548217779", "97", "0", 1, 10)
  for _, v := range list {
    t.Logf("comment:%+v\n", v)
  }
}

func BenchmarkGetVideoComments(b *testing.B) {
  for i := 0; i < b.N; i++ {
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    list := svc.GetVideoComments("202009181548217779", "97", "0", 1, 10)
    b.Logf("list:%+v\n", list)
  }
}
