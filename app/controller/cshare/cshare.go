package cshare

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/models/mcommunity"
	"sports_service/server/models/mposting"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
)

type ShareModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	posting     *mposting.PostingModel
	video       *mvideo.VideoModel
	community   *mcommunity.CommunityModel
}

func New(c *gin.Context) ShareModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return ShareModule{
		context: c,
		user: muser.NewUserModel(socket),
		posting: mposting.NewPostingModel(socket),
		video: mvideo.NewVideoModel(socket),
		community: mcommunity.NewCommunityModel(socket),
		engine: socket,
	}
}

// 分享/转发数据
func (svc *ShareModule) ShareData(userId string) int {
	return errdef.SUCCESS
}
