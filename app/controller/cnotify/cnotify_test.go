package cnotify

import (
  "github.com/gin-gonic/gin"
  "net/http/httptest"
  "sports_service/server/models/mnotify"

  //"sports_service/server/models/mnotify"
  "testing"
  "sports_service/server/dao"
)

func init() {
  dao.Engine = dao.InitXorm("root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4", []string{"root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4"})
  dao.InitRedis("127.0.0.1:6379", "")
}

// 保存用户通知设置
func TestSaveUserNotifySetting(t *testing.T) {
 param := &mnotify.NotifySettingParams{
   CommentPushSet: 1,
   ThumbUpPushSet: 1,
   AttentionPushSet: 0,
   SharePushSet: 1,
   SlotPushSet: 1,
 }

 c, _ := gin.CreateTestContext(httptest.NewRecorder())
 svc := New(c)
 syscode := svc.SaveUserNotifySetting("202009181548217779", param)
 t.Logf("\n syscode:%d", syscode)
}

func BenchmarkSaveUserNotifySetting(b *testing.B) {
 for i := 0; i < b.N; i++ {
   param := &mnotify.NotifySettingParams{
     CommentPushSet:   1,
     ThumbUpPushSet:   1,
     AttentionPushSet: 0,
     SharePushSet:     1,
     SlotPushSet:      1,
   }

   c, _ := gin.CreateTestContext(httptest.NewRecorder())
   svc := New(c)
   syscode := svc.SaveUserNotifySetting("202009181548217779", param)
   b.Logf("\n syscode:%d", syscode)
 }
}

// 用户被点赞的作品列表
func TestGetBeLikedList(t *testing.T) {
 c, _ := gin.CreateTestContext(httptest.NewRecorder())
 svc := New(c)
 list := svc.GetBeLikedList("202010101545291936", 1, 10)
 t.Logf("\n list:%+v", list)
}

func BenchmarkBeLikedList(b *testing.B) {
 for i := 0; i < b.N; i++ {
   c, _ := gin.CreateTestContext(httptest.NewRecorder())
   svc := New(c)
   list := svc.GetBeLikedList("202010101545291936", 1, 10)
   b.Logf("\n list len:%d", len(list))
 }
}

// 获取用户 @ 通知
func TestGetReceiveAtNotify(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  list := svc.GetReceiveAtNotify("202009181548217779", 1, 10)
  t.Logf("\n list len: %d", len(list))
}

func BenchmarkGetReceiveAtNotify(b *testing.B) {
  for i := 0; i < b.N; i++ {
    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := New(c)
    list := svc.GetReceiveAtNotify("202009181548217779", 1, 10)
    b.Logf("\n list len: %d", len(list))
  }
}

// 获取用户消息设置
func TestGetUserNotifySetting(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  info := svc.GetUserNotifySetting("202009181548217779")
  t.Logf("\n info:%+v", info)
}

// 获取未读消息数
func TestGetUnreadNum(t *testing.T) {
  c, _ := gin.CreateTestContext(httptest.NewRecorder())
  svc := New(c)
  info := svc.GetUnreadNum("202010101545291936")
  t.Logf("\n unread info:%+v", info)
}

