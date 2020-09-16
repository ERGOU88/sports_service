package cvideo

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/models/mvideo"
)

type VideoModule struct {
	context      *gin.Context
	engine       *xorm.Session
	video        *mvideo.VideoModel
}

func New(c *gin.Context) VideoModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return VideoModule{
		context: c,
		video: mvideo.NewVideoModel(socket),
		engine: socket,
	}
}

// 用户发布的视频列表
func (svc *VideoModule) UserPublishVideos() {
	return
}

// 用户点赞的视频列表
func (svc *VideoModule) UserZanVideos() {
	return
}