package cuser

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/models/muser"
)

type UserModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	social      *muser.SocialModel
}

func New(c *gin.Context) UserModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return UserModule{
		context: c,
		user: muser.NewUserModel(socket),
		social: muser.NewSocialPlatform(socket),
		engine:  socket,
	}
}



