package cuser

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	"sports_service/server/tools/filter"
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

// 通过userId获取用户信息
func (svc *UserModule) GetUserInfoByUserid(userId string) (int, *muser.UserInfoResp) {
	info := svc.user.FindUserByUserid(userId)
	if info == nil {
		log.Log.Errorf("user_trace: user not found, uid:%s", userId)
		return errdef.USER_NOT_EXISTS, nil
	}

	resp := &muser.UserInfoResp{
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
	nameLen := svc.GetStrLen([]rune(nickName))
	if nameLen < consts.MIN_NAME_LEN || nameLen > consts.MAX_NAME_LEN {
		log.Log.Errorf("user_trace: invalid nickname len, length:%d", nameLen)
		return errdef.USER_INVALID_NAME_LEN
	}

	signLen := svc.GetStrLen([]rune(params.Signature))
	if signLen > consts.MAX_SIGNATURE_LEN {
		log.Log.Errorf("user_trace: invalid signature len, length:%d", signLen)
		return errdef.USER_INVALID_SIGN_LEN
	}

	mfilter := filter.NewFilterModel()
	// 校验昵称是否存在敏感词
	pass, err := mfilter.ValidateText(nickName)
	if err != nil || !pass {
		log.Log.Errorf("user_trace: validate nickname err: %s，pass: %v", err, pass)
		return errdef.USER_INVALID_NAME
	}

	// 校验签名是否存在敏感词
	pass, err = mfilter.ValidateText(params.Signature)
	if !pass || err != nil {
		log.Log.Errorf("user_trace: validate signature err: %s，pass: %v", err, pass)
		return errdef.USER_INVALID_SIGNATURE
	}

	// 昵称是否重复
	if isRepeat := svc.user.IsRepeatOfNickName(nickName); isRepeat {
		log.Log.Errorf("user_trace: nick name has exists, nickname:%s", nickName)
		return errdef.USER_NICK_NAME_EXISTS
	}

	// 查看系统头像是否存在
	avatarInfo := svc.GetDefaultAvatarById(params.AvatarId)
	if avatarInfo == nil {
		log.Log.Errorf("user_trace: nick name not exists, avatar id:%d", params.AvatarId)
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

// 获取世界信息（暂时只有国家）
func (svc *UserModule) GetWorldInfo() []*models.WorldInfo {
	return svc.user.GetWorldInfo()
}

// 通过id获取世界信息（暂时只有国家）
func (svc *UserModule) GetWorldInfoById(id int32) *models.WorldInfo {
	return svc.user.GetWorldInfoById(id)
}

// 获取系统默认头像列表
func (svc *UserModule) GetDefaultAvatarList() []*models.DefaultAvatar {
	return svc.user.GetSystemAvatarList()
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



