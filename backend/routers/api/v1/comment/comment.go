package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
)

type CommentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	comment     *mcomment.CommentModel
	collect     *mcollect.CollectModel
	user        *muser.UserModel
	video       *mvideo.VideoModel
	like        *mlike.LikeModel
	attention   *mattention.AttentionModel
}

func New(c *gin.Context) CommentModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return CommentModule{
		context: c,
		comment: mcomment.NewCommentModel(socket),
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		like: mlike.NewLikeModel(socket),
		attention: mattention.NewAttentionModel(socket),
		engine: socket,
	}
}

// 获取后台视频评论列表
func (svc *CommentModule) GetVideoComments() {

}
