package cuser

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/global/rdskey"
	"sports_service/server/models"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mconfigure"
	"sports_service/server/models/mlike"
	"sports_service/server/models/mnotify"
	"sports_service/server/models/morder"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvenue"
	"sports_service/server/models/mvideo"
	"sports_service/server/models/sms"
	"sports_service/server/tools/im"
	"sports_service/server/tools/tencentCloud"
	"sports_service/server/util"
	"strings"
	"time"
	"unicode"
	"fmt"
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
	configure   *mconfigure.ConfigModel
	venue       *mvenue.VenueModel
	order       *morder.OrderModel
}

func New(c *gin.Context) UserModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()

	venueSocket := dao.VenueEngine.NewSession()
	defer venueSocket.Close()
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
		configure: mconfigure.NewConfigModel(socket),
		venue: mvenue.NewVenueModel(venueSocket),
		order: morder.NewOrderModel(venueSocket),
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

	// 查看国家是否存在
	countryInfo := svc.GetWorldInfoById(int32(info.Country))
	if countryInfo != nil {
		resp.CountryName = countryInfo.Name
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
	if isRepeat := svc.user.IsRepeatOfNickName(nickName, userId); isRepeat {
		log.Log.Errorf("user_trace: nick name has exists, nickname:%s", nickName)
		return errdef.USER_NICK_NAME_EXISTS
	}

	//if params.Avatar != 0 {
	//  // 查看系统头像是否存在
	//  avatarInfo := svc.GetDefaultAvatarById(params.Avatar)
	//  if avatarInfo == nil {
	//    log.Log.Errorf("user_trace: avatar not exists, avatar id:%d", params.Avatar)
	//    return errdef.USER_AVATAR_NOT_EXISTS
	//  }
	//
	//}

	if params.Avatar != "" {
		info.Avatar = params.Avatar
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

	if param.Describe != "" {
		client := tencentCloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.TMS_API_DOMAIN)
		// 检测反馈的内容
		isPass, err := client.TextModeration(param.Describe)
		if !isPass {
			log.Log.Errorf("user_trace: validate feedback describe err: %s，pass: %v", err, isPass)
			return errdef.USER_INVALID_FEEDBACK
		}
	}

	if err := svc.user.RecordUserFeedback(userId, param, int(time.Now().Unix())); err != nil {
		log.Log.Errorf("user_trace: record user feedback err:%s", err)
		return errdef.USER_FEEDBACK_FAIL
	}

	return errdef.SUCCESS
}

// 获取个人空间用户信息 toUserId 被查看人
func (svc *UserModule) GetUserZoneInfo(userId, toUserId string) (int, *muser.UserInfoResp, *muser.UserZoneInfoResp) {
	info := svc.user.FindUserByUserid(toUserId)
	if info == nil {
		log.Log.Errorf("user_trace: user not found, toUserid:%s", toUserId)
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

	// 查看国家是否存在
	countryInfo := svc.GetWorldInfoById(int32(info.Country))
	if countryInfo != nil {
		resp.CountryName = countryInfo.Name
	}

	if userId != toUserId {
		// 当前用户 是否关注 被查看人
		attentionInfo := svc.attention.GetAttentionInfo(userId, toUserId)
		if attentionInfo != nil {
			resp.IsAttention = int32(attentionInfo.Status)
		}
	}

	zoneRes := &muser.UserZoneInfoResp{
		// 被点赞总数
		TotalBeLiked: svc.like.GetUserTotalBeLiked(toUserId),
		// 用户关注总数
		TotalAttention: svc.attention.GetTotalAttention(toUserId),
		// 用户粉丝总数
		TotalFans: svc.attention.GetTotalFans(toUserId),
		// 用户总收藏（包含视频 和 后续的帖子）
		TotalCollect: svc.collect.GetUserTotalCollect(toUserId),
		// 用户点赞的总数
		TotalLikes: svc.like.GetUserTotalLikes(toUserId),
		// 用户发布的视频总数（已审核）
		TotalPublish: svc.video.GetTotalPublish(toUserId),
	}

	return errdef.SUCCESS, resp, zoneRes
}

// 绑定设备token
func (svc *UserModule) BindDeviceToken(userId string, param *muser.BindDeviceTokenParam) int {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("user_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS
	}

	// 一致 则不做操作
	if strings.Compare(user.DeviceToken, param.DeviceToken) == 0 {
		return errdef.SUCCESS
	}

	svc.user.SetDeviceToken(param.DeviceToken)
	svc.user.SetDeviceType(param.Platform)
	// 条件
	condition := fmt.Sprintf("id=%d", user.Id)
	// 字段
	cols := "device_token, device_type"
	if _, err := svc.user.UpdateUserInfos(condition, cols); err != nil {
		log.Log.Errorf("user_trace: bind device token fail, userId:%s", userId)
		return errdef.USER_BIND_DEVICE_TOKEN
	}

	return errdef.SUCCESS
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

// 获取卡包信息
func (svc *UserModule) GetKabawInfo(userId string) (int, *muser.UserKabawInfo) {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("user_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS, nil
	}

	ok, err := svc.venue.GetVenueInfoById("1")
	if !ok || err != nil {
		log.Log.Error("user_trace: get venue info fail, venueId:", 1)
		return errdef.VENUE_NOT_EXISTS, nil
	}

	// 暂时只有一个场馆
	ok, err = svc.venue.GetVenueVipInfo(userId, 1)
	if err != nil {
		log.Log.Errorf("user_trace: get venue vip info fail, userId:%s, err:%s", userId, err)
		return errdef.VENUE_VIP_INFO_FAIL, nil
	}

	kabaw := &muser.UserKabawInfo{}
	// 默认非会员
	kabaw.StartTm = 0
	kabaw.EndTm = 0
	kabaw.QrCodeInfo = fmt.Sprintf("U%d", util.GetSnowId())
	kabaw.RemainDuration = 0
	kabaw.Tips = "对准闸机扫描口 高度5cm刷码入场"
	kabaw.VenueName = svc.venue.Venue.VenueName
	kabaw.IsVip = false
	// 存在会员信息
	if ok {
		// 表示已注册为会员
		if svc.venue.Vip.StartTm > 0 {
			ok, err := svc.venue.GetVenueProductByType(svc.venue.Vip.VipType)
			if err != nil {
				log.Log.Errorf("user_trace: get product by type fail,productType:%d, err:%s", svc.venue.Vip.VipType, err)
			}

			kabaw.VipName = fmt.Sprintf("%s%s", svc.venue.Venue.VenueName, "会员")
			if ok {
				kabaw.VipImage = svc.venue.Product.Image
				kabaw.VipName = svc.venue.Product.ProductName
			}

			kabaw.IsVip = true

			kabaw.StartTm = svc.venue.Vip.StartTm
			kabaw.EndTm = svc.venue.Vip.EndTm
			kabaw.RemainDuration = svc.venue.Vip.Duration
			// 查看会员是否过期 已过期
			if svc.venue.Vip.EndTm < time.Now().Unix() {
				// 会员已过期
				kabaw.HasExpire = true
			}
		}
	}

	if err := svc.order.SaveQrCodeInfo(kabaw.QrCodeInfo, userId, rdskey.KEY_EXPIRE_MIN * 60); err != nil {
		log.Log.Errorf("user_trace: save qrcode kabaw info fail, userId:%s, err:%s", userId, err)
		return errdef.VENUE_VIP_INFO_FAIL, nil
	}

	return errdef.SUCCESS, kabaw
}

// 更新游客的腾讯im签名
func (svc *UserModule) UpdateTencentImSignByGuest() (int, string) {
	code, info := svc.GetTencentImSignByGuest(2)
	return code, info.Sign
}

// 更新用户/游客的腾讯im签名
func (svc *UserModule) UpdateTencentImSign(userId string) (int, string) {
	if userId == "" {
		// 更新游客im签名
		return svc.UpdateTencentImSignByGuest()
	}

	// 更新用户im签名
	return svc.UpdateTencentImSignByUser(userId)
}

// 更新用户im签名
func (svc *UserModule) UpdateTencentImSignByUser(userId string) (int, string) {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("user_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS, ""
	}

	sig, err := im.GenSig(userId, 3600 * 24 * 90)
	if err != nil {
		return errdef.USER_GEN_IM_SIGN_FAIL, ""
	}

	svc.user.User.TxAccid = userId
	svc.user.User.TxToken = sig
	// 条件
	condition := fmt.Sprintf("id=%d", user.Id)
	// 字段
	cols := "tx_accid, tx_token"
	if _, err := svc.user.UpdateUserInfos(condition, cols); err != nil {
		log.Log.Errorf("user_trace: update user im sign fail, userId:%s, err:%s", userId, err)
		return errdef.USER_UPDATE_IM_SIGN_FAIL, ""
	}

	return errdef.SUCCESS, sig
}

// 腾讯im 添加游客
func (svc *UserModule) AddGuestByTencentIm() (int, *muser.TencentImUser) {
	userId := fmt.Sprint(util.GetSnowId())
	avatar := consts.DEFAULT_AVATAR
	nickName := fmt.Sprintf("游客%d", util.GenerateRandnum(100000, 999999))

	sign, err := im.Im.AddUser(userId, nickName, avatar)
	if err != nil {
		log.Log.Errorf("user_trace: register im user fail, err:%s", err)
		return errdef.USER_ADD_GUEST_FAIL, nil
	}

	info := &muser.TencentImUser{
		UserId: userId,
		Avatar: avatar,
		NickName: nickName,
		Sign: sign,
	}

	svc.SaveGuestInfo(info)

	return errdef.SUCCESS, info
}

// 保存游客信息 [腾讯im]
func (svc *UserModule) SaveGuestInfo(info *muser.TencentImUser) {
	str, _ := util.JsonFast.MarshalToString(info)
	// 游客相关信息 保存到redis 便于复用
	if err := svc.SaveGuestInfoByTencentIm(str); err != nil {
		log.Log.Errorf("user_trace: save guest sign by tencent im fail, err:%s", err)
	}
}

// 保存游客签名等[腾讯im]
func (svc *UserModule) SaveGuestInfoByTencentIm(info string) error {
	rds := dao.NewRedisDao()
	return rds.Set(rdskey.USER_TENCENT_IM_GUEST_SIGN, info)
}

// 获取游客签名[腾讯im]
func (svc *UserModule) GetGuestSignByTencentIm() (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.USER_TENCENT_IM_GUEST_SIGN)
}

// 获取腾讯im签名 [游客]
// action 1 获取 2 更新
func (svc *UserModule) GetTencentImSignByGuest(action int) (int, *muser.TencentImUser) {
	// 查看redis 游客信息是否存在
	info, err := svc.GetGuestSignByTencentIm()
	if err != nil && err != redis.ErrNil {
		log.Log.Errorf("user_trace: get guest sign fail, err:%s", err)
		return errdef.USER_GET_GUEST_SIGN_FAIL, nil
	}

	if info == "" || err == redis.ErrNil {
		// 未获取到 添加
		return svc.AddGuestByTencentIm()
	}

	imUser := &muser.TencentImUser{}
	if err = util.JsonFast.UnmarshalFromString(info, imUser); err != nil {
		log.Log.Errorf("user_trace: unmarshal im user fail, err:%s", err)
		return errdef.USER_GET_GUEST_SIGN_FAIL, nil
	}

	// action == 2 签名过期 重新生成
	if action == 2 {
		sig, err := im.GenSig(imUser.UserId, 3600 * 24 * 90)
		if err != nil {
			return errdef.USER_GEN_IM_SIGN_FAIL, nil
		}

		imUser.Sign = sig
		svc.SaveGuestInfo(imUser)
	}

	return errdef.SUCCESS, imUser
}

// 获取腾讯im签名
func (svc *UserModule) GetTencentImSign(userId string) (int, *muser.TencentImUser) {
	if userId == "" {
		// 获取游客签名
		return svc.GetTencentImSignByGuest(1)
	}

	// 获取用户签名
	return svc.GetTencentImSignByUser(userId)
}

// 获取腾讯im签名[已注册用户]
func (svc *UserModule) GetTencentImSignByUser(userId string) (int, *muser.TencentImUser) {
	user := svc.user.FindUserByUserid(userId)
	if user == nil {
		log.Log.Errorf("user_trace: user not found, userId:%s", userId)
		return errdef.USER_NOT_EXISTS, nil
	}

	// 已导入
	if user.TxAccid != "" && user.TxToken != "" {
		info := &muser.TencentImUser{
			NickName: user.NickName,
			Avatar: user.Avatar,
			UserId: user.TxAccid,
			Sign: user.TxToken,
		}

		return errdef.SUCCESS, info
	}

	sign, err := im.Im.AddUser(user.UserId, user.NickName, user.Avatar)
	if err != nil {
		log.Log.Errorf("user_trace: register im user fail, userId:%s, err:%s", userId, err)
		return errdef.USER_ADD_GUEST_FAIL, nil
	}

	svc.user.User.TxAccid = userId
	svc.user.User.TxToken = sign
	// 条件
	condition := fmt.Sprintf("id=%d", user.Id)
	// 字段
	cols := "tx_accid, tx_token"
	if _, err := svc.user.UpdateUserInfos(condition, cols); err != nil {
		log.Log.Errorf("user_trace: update user im sign fail, userId:%s, err:%s", userId, err)
		return errdef.USER_UPDATE_IM_SIGN_FAIL, nil
	}

	info := &muser.TencentImUser{
		NickName: user.NickName,
		Avatar: user.Avatar,
		UserId: user.TxAccid,
		Sign: user.TxToken,
	}

	return errdef.SUCCESS, info
}
