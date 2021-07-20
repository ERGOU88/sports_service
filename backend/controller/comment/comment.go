package comment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models/mbarrage"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"strconv"
)

type CommentModule struct {
	context     *gin.Context
	engine      *xorm.Session
	comment     *mcomment.CommentModel
	video       *mvideo.VideoModel
	barrage     *mbarrage.BarrageModel
	user        *muser.UserModel
}

func New(c *gin.Context) CommentModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return CommentModule{
		context: c,
		comment: mcomment.NewCommentModel(socket),
		video: mvideo.NewVideoModel(socket),
		barrage: mbarrage.NewBarrageModel(socket),
		user: muser.NewUserModel(socket),
		engine: socket,
	}
}

// 获取后台视频评论列表
func (svc *CommentModule) GetVideoComments(queryId, sortType, condition string, page, size int) ([]*mcomment.VideoCommentInfo, int64) {
	var (
		total int64
		userId, videoId string
	)
	if queryId != "" {
		if _, err := strconv.Atoi(queryId); err != nil {
			return []*mcomment.VideoCommentInfo{}, total
		}

		// 查询用户是否存在
		user := svc.user.FindUserByUserid(queryId)
		if user != nil {
			userId = user.UserId
			total = svc.GetCommentTotalByUserId(userId)
		}

		// 查询视频是否存在
		video := svc.video.FindVideoById(queryId)
		if video != nil {
			video.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
			videoId = fmt.Sprint(video.VideoId)
			total = svc.GetCommentTotalByVideoId(videoId)
		}

		// 都不存在
		if user == nil && video == nil  {
			return []*mcomment.VideoCommentInfo{}, total
		}

	} else {
		total = svc.GetCommentTotal(consts.COMMENT_TYPE_VIDEO)
	}

	offset := (page - 1) * size
	list := svc.comment.GetVideoCommentsBySort(userId, videoId, sortType, condition, offset, size)
	if len(list) == 0 {
		return []*mcomment.VideoCommentInfo{}, total
	}

	for _, comment := range list {
		video := svc.video.FindVideoById(fmt.Sprint(comment.VideoId))
		if video != nil {
			comment.Title = video.Title
			comment.VideoDuration = video.VideoDuration
			comment.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
			comment.Cover = video.Cover
			comment.Describe = video.Describe
			comment.VideoHeight = video.VideoHeight
			comment.VideoWidth = video.VideoWidth
		}
	}

	return list, total
}

// 获取视频评论总数
func (svc *CommentModule) GetCommentTotal(commentType int) int64 {
	return svc.comment.GetVideoCommentTotal()
}

// 通过用户id获取评论总数
func (svc *CommentModule) GetCommentTotalByUserId(queryId string) int64 {
	return svc.comment.GetCommentTotalByUserId(queryId)
}

// 通过视频id获取评论总数
func (svc *CommentModule) GetCommentTotalByVideoId(videoId string) int64 {
	return svc.comment.GetCommentTotalByVideoId(videoId)
}

// 删除视频评论（物理删除）
func (svc *CommentModule) DelVideoComments(param *mcomment.DelCommentParam) int {
	comment := svc.comment.GetVideoCommentById(param.CommentId)
	if comment == nil {
		return errdef.COMMENT_NOT_EXISTS
	}

	// 0 逻辑删除
	comment.Status = 0
	condition := fmt.Sprintf("id=%d", comment.Id)
	cols := "status"
	affected, err := svc.comment.UpdateCommentInfo(condition, cols)
	if affected != 1 || err != nil {
		return errdef.COMMENT_DELETE_FAIL
	}


	//commentIds := svc.comment.GetVideoReplyIdsById(param.ComposeId)
	//ids := make([]string, 0)
	//// 递归查询
	//svc.recursionComments(&ids, &commentIds)
	//
	//// 当前评论 及 回复 一并删除
	//ids = append(ids, param.ComposeId)
	//
	//log.Log.Errorf("++++++++commentIds:%v", strings.Join(ids, ","))
	//
	//// 删除视频评论及当前评论下的回复
	//if err := svc.comment.DelVideoComments(strings.Join(ids,",")); err != nil {
	//	return errdef.COMMENT_DELETE_FAIL
	//}

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

// 获取视频弹幕列表
func (svc *CommentModule) GetVideoBarrageList(page, size int) []*mbarrage.VideoBarrageInfo {
	offset := (page - 1) * size
	list := svc.barrage.GetVideoBarrageList(offset, size)
	if len(list) == 0 {
		return []*mbarrage.VideoBarrageInfo{}
	}

	for _, v := range list {
		v.VideoAddr = svc.video.AntiStealingLink(v.VideoAddr)
	}

	return list
}

// 获取视频弹幕总数（管理后台）
func (svc *CommentModule) GetVideoBarrageTotal() int64 {
	return svc.barrage.GetVideoBarrageTotal()
}

// 删除视频弹幕
func (svc *CommentModule) DelVideoBarrage(param *mbarrage.DelBarrageParam) error {
	return svc.barrage.DelVideoBarrage(param.Id)
}
