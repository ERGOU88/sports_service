package cpost

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mposting"
	"sports_service/server/models/muser"
	redismq "sports_service/server/redismq/event"
)

type PostModule struct {
	context      *gin.Context
	engine       *xorm.Session
	post         *mposting.PostingModel
	attention    *mattention.AttentionModel
	user         *muser.UserModel
}

func New(c *gin.Context) PostModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return PostModule{
		context: c,
		post: mposting.NewPostingModel(socket),
		attention: mattention.NewAttentionModel(socket),
		user: muser.NewUserModel(socket),
		engine: socket,
	}
}

// 帖子审核
func (svc *PostModule) AudiPost(param *mposting.AudiPostParam) int {
	post, err := svc.post.GetPostById(param.Id)
	if post == nil || err != nil {
		return errdef.POST_NOT_FOUND
	}

	status := fmt.Sprint(post.Status)
	// 帖子已删除
	if status == consts.POST_DELETE_STATUS {
		svc.engine.Rollback()
		return errdef.POST_ALREADY_DELETE
	}

	// 帖子已通过审核 只能执行逻辑删除
	if status == consts.POST_AUDIT_SUCCESS && fmt.Sprint(param.Status) != consts.POST_DELETE_STATUS {
		svc.engine.Rollback()
		return errdef.POST_ALREADY_PASS
	}

	// 通过 / 不通过 / 执行删除操作 且 视频状态为审核通过 则只能逻辑删除/不通过 直接更新视频状态
	if fmt.Sprint(param.Status) == consts.POST_AUDIT_SUCCESS || fmt.Sprint(param.Status) == consts.POST_AUDIT_FAILURE ||
		(fmt.Sprint(param.Status) == consts.POST_DELETE_STATUS && status == consts.POST_AUDIT_SUCCESS) {
		post.Status = param.Status
		// 更新视频状态
		if err := svc.post.UpdateStatusByPost(); err != nil {
			svc.engine.Rollback()
			return errdef.POST_EDIT_STATUS_FAIL
		}

		// 如果是审核通过
		if fmt.Sprint(param.Status) == consts.POST_AUDIT_SUCCESS {
			// 获取发布者用户信息
			user := svc.user.FindUserByUserid(post.UserId)
			if user != nil {
				// 获取发布者粉丝们的userId
				userIds := svc.attention.GetFansList(user.UserId)
				for _, userId := range userIds {
					// 给发布者的粉丝 发送 发布新帖子推送
					//event.PushEventMsg(config.Global.AmqpDsn, userId, user.NickName, video.Cover, "", consts.FOCUS_USER_PUBLISH_VIDEO_MSG)
					redismq.PushEventMsg(redismq.NewEvent(userId, fmt.Sprint(post.Id), user.NickName, "", "", consts.FOCUS_USER_PUBLISH_POST_MSG))
				}
			}
		}

		svc.engine.Commit()
		return errdef.SUCCESS
	}

	// 如果执行删除操作 且 视频状态未审核通过 删除相关所有数据
	if fmt.Sprint(param.Status) == consts.VIDEO_DELETE_STATUS && status != consts.VIDEO_AUDIT_SUCCESS {
		// 物理删除发布的帖子、帖子所属话题、帖子统计数据
		if err := svc.post.DelPublishPostById(param.Id); err != nil {
			svc.engine.Rollback()
			return errdef.POST_DELETE_PUBLISH_FAIL
		}

		// 删除帖子所属话题
		if err := svc.post.DelPostTopics(param.Id); err != nil {
			svc.engine.Rollback()
			return errdef.POST_DELETE_TOPIC_FAIL
		}

		// 删除帖子统计数据
		if err := svc.post.DelPostStatistic(param.Id); err != nil {
			svc.engine.Rollback()
			return errdef.POST_DELETE_STATISTIC_FAIL
		}
	}

	svc.engine.Commit()

	return errdef.SUCCESS
}




