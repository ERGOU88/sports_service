package cposting

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"sports_service/server/dao"
	"sports_service/server/models/mposting"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	dao.AppEngine = dao.InitXorm("root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4", []string{"root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4"})
	dao.InitRedis("192.168.5.12:6378", "")
}

func TestSanitizeHtml(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	svc := New(c)
	content := svc.SanitizeHtml(`<body><p>hello~ world</p><script language="JavaScript"> <!-- while (true) {window.open("URI");} --> </script> <p>这是一条测试记录</p><a href='www.baidu.com'>@百度</a><br/></body>`)
	t.Logf("\n content:%s", content)
}

func TestPublishPosting(t *testing.T) {
	Convey("发布帖子", t, func() {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		svc := New(c)
		code := svc.PublishPosting("13918242", &mposting.PostPublishParam{
			Title:  "测试",
			Describe: "测试描述",
			ImagesAddr: []string{"https://fpv-1251316249.cos.ap-shanghai.myqcloud.com/fpv/1607511255127.png","https://fpv-1251316249.cos.ap-shanghai.myqcloud.com/fpv/1607511255127.png"},
			SectionId: 1,
			TopicIds: []string{"1","2"},
			AtInfo: []string{"10240133"},
		})

		So(code,  ShouldEqual, 200)
	})
}
