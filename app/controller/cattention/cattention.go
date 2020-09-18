package cattention

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mattention"
	"sports_service/server/models/muser"
	"strings"
	"time"
)

type AttentionModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	attention   *mattention.AttentionModel
}

func New(c *gin.Context) AttentionModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return AttentionModule{
		context: c,
		user: muser.NewUserModel(socket),
		attention: mattention.NewAttentionModel(socket),
		engine:  socket,
	}
}

// 添加关注 attentionUid 关注的用户id userId 被关注的用户id
func (svc *AttentionModule) AddAttention(attentionUid, userId string) int {
	if attentionUid == userId {
		log.Log.Errorf("attention_trace: don't focus yourself, attentionUid:%s, userId", attentionUid, userId)
		return errdef.ATTENTION_YOURSELF_FAIL
	}
	// 关注的用户是否存在
	if attentionUser := svc.user.FindUserByUserid(attentionUid); attentionUser == nil {
		log.Log.Errorf("attention_trace: user not found, attentionUid:%s", attentionUid)
		return errdef.USER_NOT_EXISTS
	}

	// 被关注的用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("attention_trace: user not found, userId:%s", userId)
		return errdef.ATTENTION_USER_NOT_EXISTS
	}

	info := svc.attention.GetAttentionInfo(attentionUid, userId)
	// 是否已关注
	if info != nil && info.Status == consts.ALREADY_ATTENTION {
		log.Log.Errorf("attention_trace: already attention, attentionUid:%s, userId:%s", attentionUid, userId)
		return errdef.ATTENTION_ALREADY_EXISTS
	}

	// 未关注 添加/更新
	// 记录存在 且 状态为未关注 之前取消了关注 更新状态为关注
	if info != nil && info.Status == consts.NO_ATTENTIO {
		info.Status = consts.ALREADY_ATTENTION
		info.CreateAt = int(time.Now().Unix())
		if err := svc.attention.UpdateAttentionStatus(); err != nil {
			log.Log.Errorf("attention_trace: update attention status err:%s", err)
			return errdef.ATTENTION_USER_FAIL
		}
	}

	// 添加关注记录
	if err := svc.attention.AddAttention(attentionUid, userId, consts.ALREADY_ATTENTION); err != nil {
		log.Log.Errorf("attention_trace: add attention err:%s", err)
		return errdef.ATTENTION_USER_FAIL
	}

	return errdef.SUCCESS
}

// 取消关注
func (svc *AttentionModule) CancelAttention(attentionUid, userId string) int {
	// 关注的用户是否存在
	if attentionUser := svc.user.FindUserByUserid(attentionUid); attentionUser == nil {
		log.Log.Errorf("attention_trace: user not found, attentionUid:%s", attentionUid)
		return errdef.USER_NOT_EXISTS
	}

	// 被关注的用户是否存在
	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("attention_trace: user not found, userId:%s", userId)
		return errdef.ATTENTION_USER_NOT_EXISTS
	}

	info := svc.attention.GetAttentionInfo(attentionUid, userId)
	// 是否已关注
	// 记录不存在 未关注
	if info == nil {
		log.Log.Errorf("attention_trace: record not found, no attention, attentionUid:%s, userId:%s", attentionUid, userId)
		return errdef.ATTENTION_RECORD_NOT_EXISTS
	}

	// 状态已是未关注 提示重复取消
	if info.Status != consts.ALREADY_ATTENTION {
		log.Log.Errorf("attention_trace: already cancel attention, attentionUid:%s, userId:%s", attentionUid, userId)
		return errdef.ATTENTION_REPEAT_CANCEL
	}

	// 更新状态为取消关注
	info.Status = consts.NO_ATTENTIO
	info.CreateAt = int(time.Now().Unix())
	if err := svc.attention.UpdateAttentionStatus(); err != nil {
		log.Log.Errorf("attention_trace: update attention status err:%s", err)
		return errdef.ATTENTION_CANCEL_FAIL
	}

	return errdef.SUCCESS
}

// 获取关注的用户列表
func (svc *AttentionModule) GetAttentionUserList(attentionUid string, page, size int) []*muser.UserInfoResp {
	userIds := svc.attention.GetAttentionList(attentionUid)
	if len(userIds) == 0 {
		log.Log.Errorf("attention_trace: not following any users")
		return nil
	}

	offset := (page - 1) * size
	uids := strings.Join(userIds, ",")
	userList := svc.user.FindUserByUserids(uids, offset, size)
	if len(userList) == 0 {
		log.Log.Errorf("attention_trace: not found user list info, len:%d, uids:%s", len(userList), uids)
		return nil
	}

	resp := make([]*muser.UserInfoResp, len(userList))
	for index, user := range userList {
		info := &muser.UserInfoResp{
			NickName:  user.NickName,
			UserId: user.UserId,
			Avatar: user.Avatar,
			MobileNum: int32(user.MobileNum),
			Gender: int32(user.Gender),
			Signature: user.Signature,
			Status: int32(user.Status),
			IsAnchor: int32(user.IsAnchor),
			BackgroundImg: user.BackgroundImg,
			Born: user.Born,
			Age: user.Age,
			UserType: user.UserType,
			Country: int32(user.Country),
		}

		resp[index] = info
	}

	return resp
}

// 获取粉丝列表
func (svc *AttentionModule) GetFansList(userId string, page, size int) []*muser.UserInfoResp {
	userIds := svc.attention.GetFansList(userId)
	if len(userIds) == 0 {
		log.Log.Errorf("attention_trace: not has any fans")
		return nil
	}

	offset := (page - 1) * size
	uids := strings.Join(userIds, ",")
	userList := svc.user.FindUserByUserids(uids, offset, size)
	if len(userList) == 0 {
		log.Log.Errorf("attention_trace: not found user list info, len:%d, uids:%s", len(userList), uids)
		return nil
	}

	// 重新组装数据
	resp := make([]*muser.UserInfoResp, len(userList))
	for index, user := range userList {
		info := &muser.UserInfoResp{
			NickName:  user.NickName,
			UserId: user.UserId,
			Avatar: user.Avatar,
			MobileNum: int32(user.MobileNum),
			Gender: int32(user.Gender),
			Signature: user.Signature,
			Status: int32(user.Status),
			IsAnchor: int32(user.IsAnchor),
			BackgroundImg: user.BackgroundImg,
			Born: user.Born,
			Age: user.Age,
			UserType: user.UserType,
			Country: int32(user.Country),
		}

		resp[index] = info
	}

	return resp
}
