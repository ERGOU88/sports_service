package cuser

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/log"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"time"
)

type UserModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	attention   *mattention.AttentionModel
	like        *mlike.LikeModel
	collect     *mcollect.CollectModel
	video       *mvideo.VideoModel
	comment     *mcomment.CommentModel
}

func New(c *gin.Context) UserModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return UserModule{
		context: c,
		user: muser.NewUserModel(socket),
		attention: mattention.NewAttentionModel(socket),
		like: mlike.NewLikeModel(socket),
		collect: mcollect.NewCollectModel(socket),
		video: mvideo.NewVideoModel(socket),
		comment: mcomment.NewCommentModel(socket),
		engine: socket,
	}
}

// 后台获取用户列表
func (svc *UserModule) GetUserList(page, size int) []*muser.UserInfo {
	offset := (page - 1) * size
	list := svc.user.GetUserList(offset, size)
	if len(list) == 0 {
		return []*muser.UserInfo{}
	}

	res := make([]*muser.UserInfo, len(list))
	for index, info := range list {
		resp := &muser.UserInfo{
			UserId: info.UserId,
			Avatar: info.Avatar,
			MobileNum: int32(info.MobileNum),
			NickName: info.NickName,
			Gender: int32(info.Gender),
			Signature: info.Signature,
			Status: int32(info.Status),
			IsAnchor: int32(info.IsAnchor),
			BackgroundImg: info.BackgroundImg,
			Born: info.Born,
			Age: info.Age,
			Country: int32(info.Country),
			RegIp: info.RegIp,
			LastLoginTm: info.LastLoginTime,
			Platform: info.DeviceType,
			Id: info.Id,
			// 被点赞总数
			TotalBeLiked: svc.like.GetUserTotalBeLiked(info.UserId),
			// 用户关注总数
			TotalAttention: svc.attention.GetTotalAttention(info.UserId),
			// 用户粉丝总数
			TotalFans: svc.attention.GetTotalFans(info.UserId),
			// 用户总收藏（包含视频 和 后续的帖子）
			TotalCollect: svc.collect.GetUserTotalCollect(info.UserId),
			// 用户点赞的总数
			TotalLikes: svc.like.GetUserTotalLikes(info.UserId),
			// 用户发布的视频总数（已审核）
			TotalPublish: svc.video.GetTotalPublish(info.UserId),
			// todo: 弹幕
			TotalBarrage: 0,
			// 用户总评论数
			TotalComment: svc.comment.GetUserTotalComments(info.UserId),
			// 用户总浏览数
			TotalBrowse: svc.video.GetUserTotalBrowse(info.UserId),
		}

		res[index] = resp
	}

	return res
}

// 封禁用户
func (svc *UserModule) ForbidUser(id string) int {
	svc.user.User.UpdateAt = int(time.Now().Unix())
	svc.user.User.Status = consts.USER_FORBID
	if err := svc.user.UpdateUserStatus(id); err != nil {
		log.Log.Errorf("user_trace: forbid user err:%s", err)
		return errdef.USER_FORBID_FAIL
	}

	return errdef.SUCCESS
}

// 解封用户
func (svc *UserModule) UnForbidUser(id string) int {
	svc.user.User.UpdateAt = int(time.Now().Unix())
	svc.user.User.Status = consts.USER_NORMAL
	if err := svc.user.UpdateUserStatus(id); err != nil {
		log.Log.Errorf("user_trace: un forbid user err:%s", err)
		return errdef.USER_UNFORBID_FAIL
	}

	return errdef.SUCCESS
}
