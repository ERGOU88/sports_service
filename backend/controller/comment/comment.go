package comment

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mvideo"
	"fmt"
	"strings"
)

type CommentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	comment     *mcomment.CommentModel
	video       *mvideo.VideoModel
}

func New(c *gin.Context) CommentModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return CommentModule{
		context: c,
		comment: mcomment.NewCommentModel(socket),
		video: mvideo.NewVideoModel(socket),
		engine: socket,
	}
}

// 获取后台视频评论列表
func (svc *CommentModule) GetVideoComments(sortType string, page, size int) []*mcomment.VideoCommentInfo {
	offset := (page - 1) * size
	list := svc.comment.GetVideoCommentsBySort(sortType, offset, size)
	if len(list) == 0 {
		return  []*mcomment.VideoCommentInfo{}
	}

	for _, comment := range list {
		video := svc.video.FindVideoById(fmt.Sprint(comment.VideoId))
		if video != nil {
			comment.Title = video.Title
			comment.VideoDuration = video.VideoDuration
			comment.VideoAddr = video.VideoAddr
			comment.Cover = video.Cover
			comment.Describe = video.Describe
			comment.VideoHeight = video.VideoHeight
			comment.VideoWidth = video.VideoWidth
		}
	}

	return list
}

// 删除视频评论（物理删除）
func (svc *CommentModule) DelVideoComments(param *mcomment.DelCommentParam) int {
	comment := svc.comment.GetVideoCommentById(param.CommentId)
	if comment == nil {
		return errdef.COMMENT_NOT_EXISTS
	}

	commentIds := svc.comment.GetVideoReplyIdsById(param.CommentId)
	// 当前评论 及 回复 一并删除
	commentIds = append(commentIds, param.CommentId)
	// 删除视频评论及当前评论下的回复
	if err := svc.comment.DelVideoComments(strings.Join(commentIds,",")); err != nil {
		return errdef.COMMENT_DELETE_FAIL
	}

	return errdef.SUCCESS
}
