package cvideo

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"sports_service/dao"
	"sports_service/util"
	"testing"
)

func init() {
	dao.AppEngine = dao.InitXorm("root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4", []string{"root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4"})
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

func TestIntMp(t *testing.T) {
	mp := util.NewIntMap(4)
	mp.Insert(100, 0)
	mp.Insert(2, 1)
	mp.Insert(3, 2)
	mp.Insert(9, 3)
	key, v, b := mp.GetByOrderIndex(mp.Size() - 1)
	t.Logf("mp:%+v, key:%d, v:%v, size:%d", mp, key, v, mp.Size())
	if b {
		switch v {
		case 0:
			t.Logf("点赞数：%d", key)
		case 1:
			t.Logf("收藏数：%d", key)
		case 2:
			t.Logf("弹幕数：%d", key)
		case 3:
			t.Logf("评论数：%d", key)
		}
	}
}

// 转化为中文展示
func TestTransferChinese(t *testing.T) {
	chinese := util.TransferChinese(20)
	t.Logf("chinese:%s", chinese)
}

// 视频详情页 标签推荐视频
func TestGetDetailRecommend(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	svc := New(c)
	list := svc.GetDetailRecommend("202009101933004667", "59", 1, 10)
	for _, val := range list {
		t.Logf("video:%+v\n", val)
	}
}
