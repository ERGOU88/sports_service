package comment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/log"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mvideo"
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

// 获取视频评论总数
func (svc *CommentModule) GetCommentTotal() int64 {
  return svc.comment.GetCommentTotal()
}

// 删除视频评论（物理删除）
func (svc *CommentModule) DelVideoComments(param *mcomment.DelCommentParam) int {
	comment := svc.comment.GetVideoCommentById(param.CommentId)
	if comment == nil {
		return errdef.COMMENT_NOT_EXISTS
	}

	commentIds := svc.comment.GetVideoReplyIdsById(param.CommentId)
	ids := make([]string, 0)
	// 递归查询
	svc.recursionComments(&ids, &commentIds)

	// 当前评论 及 回复 一并删除
	ids = append(ids, param.CommentId)

	log.Log.Errorf("++++++++commentIds:%v", strings.Join(ids, ","))

	// 删除视频评论及当前评论下的回复
	if err := svc.comment.DelVideoComments(strings.Join(ids,",")); err != nil {
		return errdef.COMMENT_DELETE_FAIL
	}

	return errdef.SUCCESS
}

// 递归查询当前要删除的评论下的所有回复
func (svc *CommentModule) recursionComments(ids *[]string, commentIds *[]string) {
	*ids = append(*ids, *commentIds...)
	if len(*commentIds) > 0 {
		for _, commentId := range *commentIds {
			replyIds := svc.comment.GetVideoReplyIdsById(commentId)
			svc.recursionComments(ids, &replyIds)
		}
	}
}

