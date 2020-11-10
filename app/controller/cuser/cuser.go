package cuser

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mlike"
	"sports_service/server/models/mnotify"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"sports_service/server/models/sms"
	"sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"strings"
	"time"
	"unicode"
)

type UserModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	social      *muser.SocialModel
	notify      *mnotify.NotifyModel
	attention   *mattention.AttentionModel
	like        *mlike.LikeModel
	collect     *mcollect.CollectModel
	video       *mvideo.VideoModel
	sms         *sms.SmsModel
}

func New(c *gin.Context) UserModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return UserModule{
		context: c,
		user: muser.NewUserModel(socket),
		social: muser.NewSocialPlatform(socket),
		notify: mnotify.NewNotifyModel(socket),
		attention: mattention.NewAttentionModel(socket),
		like: mlike.NewLikeModel(socket),
		collect: mcollect.NewCollectModel(socket),
		video: mvideo.NewVideoModel(socket),
		sms: sms.NewSmsModel(),
		engine: socket,
	}
}

// 通过userId获取用户信息
func (svc *UserModule) GetUserInfoByUserid(userId string) (int, *muser.UserInfoResp) {
	info := svc.user.FindUserByUserid(userId)
	if info == nil {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId)
		return errdef.USER_NOT_EXISTS, nil
	}

	// 重新组装返回数据
	resp := &muser.UserInfoResp{
		UserId: info.UserId,
		Avatar: info.Avatar,
		MobileNum: info.MobileNum,
		NickName: info.NickName,
		Gender: int32(info.Gender),
		Signature: info.Signature,
		Status: int32(info.Status),
		IsAnchor: int32(info.IsAnchor),
		BackgroundImg: info.BackgroundImg,
		Born: info.Born,
		Age: info.Age,
		UserType: info.UserType,
		Country: int32(info.Country),
	}

	return errdef.SUCCESS, resp
}

// 修改用户信息 todo: 脏词过滤接第三方？
func (svc *UserModule) EditUserInfo(userId string, params *muser.EditUserInfoParams) int {
	info := svc.user.FindUserByUserid(userId)
	if info == nil {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	// 去掉空格 换行
	nickName := strings.Trim(params.NickName, " ")
	nickName = strings.Replace(nickName, "\n", "", -1)
	nameLen := util.GetStrLen([]rune(nickName))
	if nameLen < consts.MIN_NAME_LEN || nameLen > consts.MAX_NAME_LEN {
		log.Log.Errorf("user_trace: invalid nickname len, length:%d", nameLen)
		return errdef.USER_INVALID_NAME_LEN
	}

	signLen := util.GetStrLen([]rune(params.Signature))
	if signLen > consts.MAX_SIGNATURE_LEN {
		log.Log.Errorf("user_trace: invalid signature len, length:%d", signLen)
		return errdef.USER_INVALID_SIGN_LEN
	}

	//mfilter := filter.NewFilterModel()
	// 校验昵称是否存在敏感词
	//pass, err := mfilter.ValidateText(nickName)
	//if err != nil || !pass {
	//	log.Log.Errorf("user_trace: validate nickname err: %s，pass: %v", err, pass)
	//	return errdef.USER_INVALID_NAME
	//}

	// 校验签名是否存在敏感词
	//pass, err = mfilter.ValidateText(params.Signature)
	//if !pass || err != nil {
	//	log.Log.Errorf("user_trace: validate signature err: %s，pass: %v", err, pass)
	//	return errdef.USER_INVALID_SIGNATURE
	//}

	client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测昵称
	isPass, err := client.TextModeration(nickName)
	if !isPass {
		log.Log.Errorf("user_trace: validate nick name err: %s，pass: %v", err, isPass)
		return errdef.USER_INVALID_NAME
	}

	// 检测签名
	isPass, err = client.TextModeration(params.Signature)
	if !isPass {
		log.Log.Errorf("user_trace: validate signature err: %s，pass: %v", err, isPass)
		return errdef.USER_INVALID_SIGNATURE
	}

	// 昵称是否重复
	if isRepeat := svc.user.IsRepeatOfNickName(nickName); isRepeat {
		log.Log.Errorf("user_trace: nick name has exists, nickname:%s", nickName)
		return errdef.USER_NICK_NAME_EXISTS
	}

	// 查看系统头像是否存在
	avatarInfo := svc.GetDefaultAvatarById(params.Avatar)
	if avatarInfo == nil {
		log.Log.Errorf("user_trace: nick name not exists, avatar id:%d", params.Avatar)
		return errdef.USER_AVATAR_NOT_EXISTS
	}

	// 查看国家是否存在
	countryInfo := svc.GetWorldInfoById(params.CountryId)
	if countryInfo == nil {
		log.Log.Errorf("user_trace: country info not exists, country id:%d", params.CountryId)
		return errdef.USER_COUNTRY_NOT_EXISTS
	}

	// 通过出生年月日 获取 年龄
	info.Age = util.GetAge(util.GetTimeFromStrDate(params.Born))
	info.Born = params.Born
	info.Gender = int(params.Gender)
	info.NickName = nickName
	info.Signature = params.Signature
	info.Country = int(params.CountryId)
	info.Avatar = avatarInfo.Avatar
	info.UpdateAt = int(time.Now().Unix())
	if err := svc.user.UpdateUserInfo(); err != nil {
		log.Log.Errorf("user_trace: update user info err:%s", err)
		return errdef.USER_UPDATE_INFO_FAIL
	}

	return errdef.SUCCESS
}

// 记录用户反馈
func (svc *UserModule) RecordUserFeedback(userId string, param *muser.FeedbackParam) int {
	info := svc.user.FindUserByUserid(userId)
	if info == nil {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
	// 检测反馈的内容
	isPass, err := client.TextModeration(param.Describe)
	if !isPass {
		log.Log.Errorf("user_trace: validate feedback describe err: %s，pass: %v", err, isPass)
		return errdef.USER_INVALID_FEEDBACK
	}

	if err := svc.user.RecordUserFeedback(userId, param, int(time.Now().Unix())); err != nil {
		log.Log.Errorf("user_trace: record user feedback err:%s", err)
		return errdef.USER_FEEDBACK_FAIL
	}

	return errdef.SUCCESS
}

// 获取个人空间用户信息
func (svc *UserModule) GetUserZoneInfo(userId string) (int, *muser.UserInfoResp, *muser.UserZoneInfoResp) {
	info := svc.user.FindUserByUserid(userId)
	if info == nil {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId)
		return errdef.USER_NOT_EXISTS, nil, nil
	}

	// 重新组装返回数据
	resp := &muser.UserInfoResp{
		UserId: info.UserId,
		Avatar: info.Avatar,
		MobileNum: info.MobileNum,
		NickName: info.NickName,
		Gender: int32(info.Gender),
		Signature: info.Signature,
		Status: int32(info.Status),
		IsAnchor: int32(info.IsAnchor),
		BackgroundImg: info.BackgroundImg,
		Born: info.Born,
		Age: info.Age,
		UserType: info.UserType,
		Country: int32(info.Country),
	}

	zoneRes := &muser.UserZoneInfoResp{
		// 被点赞总数
		TotalBeLiked: svc.like.GetUserTotalBeLiked(userId),
		// 用户关注总数
		TotalAttention: svc.attention.GetTotalAttention(userId),
		// 用户粉丝总数
		TotalFans: svc.attention.GetTotalFans(userId),
		// 用户总收藏（包含视频 和 后续的帖子）
		TotalCollect: svc.collect.GetUserTotalCollect(userId),
		// 用户点赞的总数
		TotalLikes: svc.like.GetUserTotalLikes(userId),
		// 用户发布的视频总数（已审核）
		TotalPublish: svc.video.GetTotalPublish(userId),
	}

	return errdef.SUCCESS, resp, zoneRes
}

// 获取世界信息（暂时只有国家）
func (svc *UserModule) GetWorldInfo() []*models.WorldMap {
	list := svc.user.GetWorldInfo()
	if len(list) == 0 {
	  return []*models.WorldMap{}
  }

  return list
}

// 通过id获取世界信息（暂时只有国家）
func (svc *UserModule) GetWorldInfoById(id int32) *models.WorldMap {
	return svc.user.GetWorldInfoById(id)
}

// 获取系统默认头像列表
func (svc *UserModule) GetDefaultAvatarList() []*models.DefaultAvatar {
	list := svc.user.GetSystemAvatarList()
	if len(list) == 0 {
	  return []*models.DefaultAvatar{}
  }

  return list
}

// 通过id获取系统默认头像
func (svc *UserModule) GetDefaultAvatarById(id int32) *models.DefaultAvatar {
	return svc.user.GetSystemAvatarById(id)
}

// 获取昵称/签名 长度（产品需求：1个汉字=2个字符 昵称最多15个汉字（30个字符）最少1个字符 签名最多70个汉字（140个字符））
func (svc *UserModule) GetStrLen(r []rune) int {
	if len(r) == 0 {
		return 0
	}

	var letterlen, wordlen int
	for _, v := range r {
		// 是否为汉字
		if unicode.Is(unicode.Han, v) {
			wordlen++
			continue
		}

		letterlen++
	}

	length := letterlen + wordlen * 2
	return length
}



