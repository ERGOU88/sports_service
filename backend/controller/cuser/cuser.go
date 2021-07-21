package cuser

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/log"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mbarrage"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mcomment"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"strconv"
	"time"
	"fmt"
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
	barrage     *mbarrage.BarrageModel
}

func New(c *gin.Context) UserModule {
	socket := dao.Engine.NewSession()
	defer socket.Close()
	return UserModule{
		context: c,
		user: muser.NewUserModel(socket),
		attention: mattention.NewAttentionModel(socket),
		like: mlike.NewLikeModel(socket),
		collect: mcollect.NewCollectModel(socket),
		video: mvideo.NewVideoModel(socket),
		comment: mcomment.NewCommentModel(socket),
		barrage: mbarrage.NewBarrageModel(socket),
		engine: socket,
	}
}

// 后台获取用户列表 todo:增加排序
func (svc *UserModule) GetUserListBySort(queryId, sortType, condition string, page, size int) ([]*muser.UserInfo, int64) {
	var (
		total int64
		userId, mobileNum string
	)
	if queryId != "" {
		if _, err := strconv.Atoi(queryId); err != nil {
			return []*muser.UserInfo{}, total
		}

		// 通过uid查询用户是否存在
		user := svc.user.FindUserByUserid(queryId)
		if user != nil {
			userId = user.UserId
			total = 1
		}

		// 通过手机号查询用户是否存在
		user = svc.user.FindUserByPhone(queryId)
		if user != nil {
			mobileNum = fmt.Sprint(user.MobileNum)
			total =1
		}

		// 都不存在
		if userId == "" && mobileNum == ""  {
			return []*muser.UserInfo{}, total
		}

	} else {
		total = svc.GetUserTotalCount()
	}

	offset := (page - 1) * size
	list := svc.user.GetUserListBySort(userId, mobileNum, sortType, condition, offset, size)
	if len(list) == 0 {
		return []*muser.UserInfo{}, total
	}

	return list, total
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
			MobileNum: int64(info.MobileNum),
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
			LastLoginTime: info.LastLoginTime,
			Platform: info.DeviceType,
			UserType: int32(info.UserType),
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
			TotalBarrage: svc.barrage.GetUserTotalVideoBarrage(info.UserId),
			// 用户总评论数
			TotalComment: svc.comment.GetUserTotalComments(info.UserId),
			// 用户总浏览数
			TotalBrowse: svc.video.GetUserTotalBrowse(info.UserId),
		}

		if country := svc.GetWorldInfoById(int32(info.Country)); country != nil {
			resp.CountryCn = country.Name
		}

		res[index] = resp
	}

	return res
}

// 获取用户总数
func (svc *UserModule) GetUserTotalCount() int64 {
	return svc.user.GetUserTotalCount()
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

// 通过id获取世界信息（暂时只有国家）
func (svc *UserModule) GetWorldInfoById(id int32) *models.WorldMap {
	return svc.user.GetWorldInfoById(id)
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
