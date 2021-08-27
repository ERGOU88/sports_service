package muser

import (
  "fmt"
  "github.com/go-xorm/xorm"
  "regexp"
  "sports_service/server/dao"
  "sports_service/server/global/app/log"
  "sports_service/server/global/consts"
  "sports_service/server/global/rdskey"
  "sports_service/server/models"
  "sports_service/server/util"
  "strings"
  "time"
)

type UserModel struct {
	User    *models.User
	Engine  *xorm.Session
}

func NewUserModel(engine *xorm.Session) *UserModel {
	return &UserModel{
		User: new(models.User),
		Engine: engine,
	}
}

// 后台封禁用户请求参数
type ForbidUserParam struct {
	Id        string       `json:"id"`
}

// 后台解封用户请求参数
type UnForbidUserParam struct {
	Id        string      `json:"id"`
}

// 用户信息（后台）
type UserInfo struct {
	Id            int64  `json:"id" example:"1000000000000"`
	UserId        string `json:"user_id" example:"2009011314521111"`
	Avatar        string `json:"avatar" example:"头像地址"`
	MobileNum     int64  `json:"mobile_num" example:"13177656222"`
	NickName      string `json:"nick_name" example:"昵称 陈二狗"`
	Gender        int32  `json:"gender" example:"0"`
	Signature     string `json:"signature" example:"个性签名"`
	Status        int32  `json:"status" example:"0"`
	IsAnchor      int32  `json:"is_anchor" example:"0"`
	BackgroundImg string `json:"background_img" example:"背景图"`
	Born          string `json:"born" example:"出生日期"`
	Age           int    `json:"age" example:"27"`
	Country       int32  `json:"country" example:"0"`
	CountryCn     string `json:"country_cn" example:"中国"`
	RegIp         string `json:"reg_ip" example:"192.168.0.108"`
	LastLoginTime int    `json:"last_login_time" example:"1600000000"`
	Platform      int    `json:"platform" example:"0"`
	UserType      int32  `json:"user_type" example:"0"`

	TotalBeLiked     int64  `json:"total_beLiked" example:"100"`     // 被点赞数
	TotalFans        int64  `json:"total_fans" example:"100"`        // 粉丝数
	TotalAttention   int64  `json:"total_attention" example:"100"`   // 关注数
	TotalCollect     int64  `json:"total_collect" example:"100"`     // 收藏的作品数
	TotalPublish     int64  `json:"total_publish" example:"100"`     // 发布的作品数
	TotalLikes       int64  `json:"total_likes" example:"100"`       // 点赞的作品数
	TotalComment     int64  `json:"total_comment" example:"100"`     // 总评价数
	TotalBrowse      int64  `json:"total_browse" example:"100"`      // 总浏览数
	TotalBarrage     int64  `json:"total_barrage" example:"100"`     // 总弹幕数
}

// 用户简单信息返回
type UserInfoResp struct {
	UserId        string `json:"user_id" example:"2009011314521111"`
	Avatar        string `json:"avatar" example:"头像地址"`
	MobileNum     int64  `json:"mobile_num" example:"13177656222"`
	NickName      string `json:"nick_name" example:"昵称 陈二狗"`
	Gender        int32  `json:"gender" example:"0"`
	Signature     string `json:"signature" example:"个性签名"`
	Status        int32  `json:"status" example:"0"`
	IsAnchor      int32  `json:"is_anchor" example:"0"`
	BackgroundImg string `json:"background_img" example:"背景图"`
	Born          string `json:"born" example:"出生日期"`
	Age           int    `json:"age" example:"27"`
	UserType      int    `json:"user_type" example:"0"`
	Country       int32  `json:"country" example:"0"`
	CountryName   string `json:"country_name"`
	IsAttention   int32  `json:"is_attention" example:"0"`
	IsReplyFocus  int32  `json:"is_reply_focus"`               // 对方是否关注
}

// 个人空间用户信息
type UserZoneInfoResp struct {
	TotalBeLiked     int64  `json:"total_beLiked" example:"100"`     // 被点赞数
	TotalFans        int64  `json:"total_fans" example:"100"`        // 粉丝数
	TotalAttention   int64  `json:"total_attention" example:"100"`   // 关注数
	TotalCollect     int64  `json:"total_collect" example:"100"`     // 收藏的作品数
	TotalPublish     int64  `json:"total_publish" example:"100"`     // 发布的作品数
	TotalLikes       int64  `json:"total_likes" example:"100"`       // 点赞的作品数
}

// 用户搜索
type UserSearchResults struct {
	UserId        string `json:"user_id" example:"2009011314521111"`
	Avatar        string `json:"avatar" example:"头像地址"`
	NickName      string `json:"nick_name" example:"昵称 陈二狗"`
	Gender        int32  `json:"gender" example:"0"`
	Signature     string `json:"signature" example:"个性签名"`
	Status        int32  `json:"status" example:"0"`
	IsAnchor      int32  `json:"is_anchor" example:"0"`
	BackgroundImg string `json:"background_img" example:"背景图"`
	Born          string `json:"born" example:"出生日期"`
	Age           int    `json:"age" example:"27"`
	IsAttention   int32  `json:"is_attention"  example:"1"`
	WorksNum      int64  `json:"works_num" example:"100"`                          // 作品数
	FansNum       int64  `json:"fans_num"  example:"100"`                          // 粉丝数
}

// 登陆请求所需的参数
type LoginParams struct {
	Platform  int      `json:"platform" example:"0"`                                // 平台 0 android 1 iOS 2 web
	//Token     string   `binding:"required" json:"token" example:"客户端token"`
	OpToken   string   `binding:"required" json:"op_token" example:"客户端返回的运营商token"`
	Operator  string   `binding:"required" json:"operator" example:"客户端返回的运营商，CMCC:中国移动通信, CUCC:中国联通通讯, CTCC:中国电信"`
}

// 修改用户信息请求参数
type EditUserInfoParams struct {
	Avatar    string `json:"avatar" example:"1"`            // 头像地址（暂时仅支持更换系统默认头像）
	NickName  string `json:"nick_name" example:"陈二狗"`     // 昵称
	Born      string `json:"born" example:"1993-06-20"`     // 出生年月
	Gender    int32  `json:"gender" example:"1"`            // 性别 1 男 2 女
	CountryId int32  `json:"country_id" example:"1"`        // 国家id
	Signature string `json:"signature" example:"emmmmmmmm"` // 个性签名
}

// 用户卡包信息
type UserKabawInfo struct {
	StartTm         int64  `json:"start_tm"`       // 会员开始时间
	EndTm           int64  `json:"end_tm"`         // 会员结束时间
	QrCodeInfo      string `json:"qr_code_info"`   // 二维码信息
	VipName         string `json:"vip_name"`       // 会员名称
	RemainDuration  int64  `json:"remain_duration"`// 剩余时长
	Tips            string `json:"tips"`           // 提示
	IsVip           bool   `json:"is_vip"`         // 是否会员
	HasExpire       bool   `json:"has_expire"`     // 会员是否过期 true 过期
	VenueName       string `json:"venue_name"`     // 场馆名称
}

// 用户反馈请求参数
type FeedbackParam struct {
	Phone    string `binding:"required" json:"phone" example:"手机号"`
	Describe string `json:"describe" example:"问题描述"`                              // 反馈内容
	Problem  string `binding:"required" json:"problem" example:"遇到的问题"`           // 遇到的问题
	Pics     string `json:"pics" example:"图片列表"`                                  // 图片（多张逗号分隔）
}

// 个人空间 用户信息请求参数
type UserZoneInfoParam struct {
	UserId  string `json:"user_id"`            // 用户userid
}

// 绑定设备token
type BindDeviceTokenParam struct {
  DeviceToken       string     `json:"device_token" binding:"required"`    // 设备token（推送）
  Platform          int        `json:"platform"`                           // 平台 0 android 1 ios 2 web
}

var validPhone = regexp.MustCompile(`^1\d{10}$`)
// 检验手机号
func (m *UserModel) CheckCellPhoneNumber(mobileNum string) bool {
	return validPhone.MatchString(mobileNum)
}

// 手机号查询用户
func (m *UserModel) FindUserByPhone(mobileNum string) *models.User {
	m.User = new(models.User)
	ok, err := m.Engine.Where("mobile_num=?", mobileNum).Get(m.User)
	if !ok || err != nil {
		log.Log.Errorf("user_trace: find user by phone err:%s", err)
		return nil
	}

	return m.User
}

// userid查询用户
func (m *UserModel) FindUserByUserid(userId string) *models.User {
	m.User = new(models.User)
	ok, err := m.Engine.Where("user_id=?", userId).Get(m.User)
	if !ok || err != nil {
		log.Log.Errorf("user_trace: find user by userid err:%s, user:%v", err, m.User)
		return nil
	}

	return m.User
}

// userid列表查询用户
func (m *UserModel) FindUserByUserids(userIds string, offset, size int) []*models.User {
	var list []*models.User
	sql := fmt.Sprintf("SELECT * FROM user WHERE user_id in(%s) ORDER BY id DESC LIMIT ?, ?", userIds)
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		log.Log.Errorf("user_trace: get user list err:%s", err)
		return nil
	}

	return list
}

// 分页获取用户列表（后台）
func (m *UserModel) GetUserList(offset, size int) []*models.User {
	var list []*models.User
	if err := m.Engine.Desc("id").Limit(size, offset).Find(&list); err != nil {
		log.Log.Errorf("user_trace: get user list err:%s", err)
		return []*models.User{}
	}

	return list
}

// 根据用户id、手机号查询用户(后台用户列表排序 0 按注册时间倒序 1 关注数 2 粉丝数 3 发布数 4 浏览数 5 点赞数 6 收藏数 7 评论数 8 弹幕数)
func (m *UserModel) GetUserListBySort(userId, mobileNum, sortType, condition string, offset, size int) []*UserInfo {
  sql := "SELECT u.*, total_attention, total_fans, tu.total_be_liked, total_likes, total_collect, total_publish, " +
    "total_barrage, total_comment, total_browse from user as u " +
    "LEFT JOIN (select ua.attention_uid, ua.user_id, count(ua.id) as total_attention from user_attention as ua " +
    "where ua.status=1 group by ua.attention_uid) as ua on u.user_id=ua.attention_uid " +
    "LEFT JOIN (select ua.attention_uid, ua.user_id, count(ua.user_id) as total_fans from user_attention as ua " +
    "where ua.status=1 group by ua.user_id) as ua2 on u.user_id=ua2.user_id  " +
    "LEFT JOIN (select tu.to_user_id, tu.user_id, count(tu.id) as total_be_liked from thumbs_up as tu " +
    "where tu.status=1 group by tu.to_user_id) as tu on u.user_id=tu.to_user_id  " +
    "LEFT JOIN (select tu.to_user_id, tu.user_id, count(tu.id) as total_likes from thumbs_up as tu " +
    "where tu.status=1 group by tu.user_id) as tu2 on u.user_id=tu2.user_id  " +
    "LEFT JOIN (select cr.user_id, count(cr.id) as total_collect from collect_record as cr " +
    "where cr.status=1 group by cr.user_id) as cr on u.user_id=cr.user_id " +
    "LEFT JOIN (select v.user_id, count(v.video_id) as total_publish from videos as v " +
    " group by v.user_id) as v on u.user_id=v.user_id " +
    "LEFT JOIN (select vb.user_id, count(vb.id) as total_barrage from video_barrage as vb " +
    "group by vb.user_id) as vb on u.user_id=vb.user_id  " +
    "LEFT JOIN (select vc.user_id, count(vc.id) as total_comment from video_comment as vc " +
    "WHERE vc.status=1 group by vc.user_id) as vc on u.user_id=vc.user_id " +
    "LEFT JOIN (select ubr.user_id, count(ubr.id) as total_browse from user_browse_record as ubr " +
    "group by ubr.user_id) as ubr on u.user_id=ubr.user_id "

   if userId != "" {
     sql += fmt.Sprintf("WHERE u.user_id=%s ", userId)
   }

   if mobileNum != "" {
     sql += fmt.Sprintf("WHERE u.mobile_num=%s ", mobileNum)
   }

   switch condition {
   case consts.USER_SORT_BY_TIME:
     sql += "ORDER BY u.create_at "
   case consts.USER_SORT_BY_ATTENTION:
     sql += "ORDER BY total_attention "
   case consts.USER_SORT_BY_FANS:
     sql += "ORDER BY total_fans "
   case consts.USER_SORT_BY_PUBLISH:
     sql += "ORDER BY total_publish "
   case consts.USER_SORT_BY_BROWSE:
     sql += "ORDER BY total_browse "
   case consts.USER_SORT_BY_LIKE:
     sql += "ORDER BY total_likes "
   case consts.USER_SORT_BY_COLLECT:
     sql += "ORDER BY total_collect "
   case consts.USER_SORT_BY_COMMENT:
     sql += "ORDER BY total_comment "
   case consts.USER_SORT_BY_BARRAGE:
     sql += "ORDER BY total_barrage "
   default:
     sql += "ORDER BY u.create_at "
   }

   // 1 正序 默认倒序
   if sortType == "1" {
     sql += "ASC "
   } else {
     sql += "DESC "
   }

   sql += "LIMIT ?, ?"
   var list []*UserInfo
   if err := m.Engine.Table(&models.User{}).SQL(sql, offset, size).Find(&list); err != nil {
     log.Log.Errorf("user_trace: get user list by sort, err:%s", err)
     return []*UserInfo{}
   }

   return list

}

// 添加用户
func (m *UserModel) AddUser() error {
	if _, err := m.Engine.InsertOne(m.User); err != nil {
		log.Log.Errorf("user_trace: add user err:%s", err)
		return err
	}

	return nil
}

// 昵称是否重复
func (m *UserModel) IsRepeatOfNickName(nickName, userId string) bool {
	count, _ := m.Engine.Where("nick_name = ? and user_id != ?", nickName, userId).Count(&models.User{})
	if count > 0 {
		return true
	}

	return false
}

// 更新用户信息
func (m *UserModel) UpdateUserInfo() error {
	if _, err := m.Engine.Where("id=?", m.User.Id).
		Cols("avatar, nick_name, born, age, gender, country, signature, update_at").
		Update(m.User); err != nil {
			return err
	}

	return nil
}

// 更新用户信息
func (m *UserModel) UpdateUserInfos(condition, cols string) (int64, error) {
  return m.Engine.Where(condition).Cols(cols).Update(m.User)
}

// 更新用户状态
func (m *UserModel) UpdateUserStatus(id string) error {
	if _, err := m.Engine.Where("id=?", id).Cols("status, update_at").Update(m.User); err != nil {
		log.Log.Errorf("user_trace: update user status err:%s", err)
		return err
	}

	return nil
}

// 获取世界信息（暂时只有国家）
func (m *UserModel) GetWorldInfo() []*models.WorldMap {
	var list []*models.WorldMap
	if err := m.Engine.Where("status=0").Desc("sortorder").Find(&list); err != nil {
		return []*models.WorldMap{}
	}

	return list
}


// 通过id获取世界信息（暂时只有国家）
func (m *UserModel) GetWorldInfoById(id int32) *models.WorldMap {
	info := new(models.WorldMap)
	ok, err := m.Engine.Where("id=?", id).Get(info)
	if !ok || err != nil {
		return nil
	}

	return info
}

// 获取系统默认头像列表
func (m *UserModel) GetSystemAvatarList() []*models.DefaultAvatar {
	var list []*models.DefaultAvatar
	if err := m.Engine.Desc("sortorder").Find(&list); err != nil {
		return []*models.DefaultAvatar{}
	}

	return list
}

// 添加系统头像请求参数
type AddSystemAvatarParams struct {
	Avatar    string `json:"avatar"`
	Sortorder int    `json:"sortorder"`
}

// 添加系统头像(默认上架)
func (m *UserModel) AddSystemAvatar(params *AddSystemAvatarParams) error {
	now := int(time.Now().Unix())
	info := &models.DefaultAvatar{
		Avatar:  params.Avatar,
		CreateAt: now,
		UpdateAt: now,
		Sortorder: params.Sortorder,
		Status: 0,
	}

	if _, err := m.Engine.InsertOne(info); err != nil {
		return err
	}

	return nil
}

// 删除系统头像 请求参数
type DelSystemAvatarParam struct {
	Id      string    `json:"id"`
}
// 删除系统头像
func (m *UserModel) DelSystemAvatar(id string) error {
	if _, err := m.Engine.Where("id=?", id).Delete(&models.DefaultAvatar{}); err != nil {
		return err
	}

	return nil
}

// 记录用户反馈信息
func (m *UserModel) RecordUserFeedback(userId string, param *FeedbackParam, now int) error {
	feedback := new(models.Feedback)
	feedback.Describe = param.Describe
	feedback.Problem = param.Problem
	feedback.CreateAt = now
	feedback.UpdateAt = now
	feedback.UserId = userId
	feedback.Phone = param.Phone
	feedback.Pics = param.Pics

	if _, err := m.Engine.InsertOne(feedback); err != nil {
		return err
	}

	return nil
}

// 通过id获取系统默认头像
func (m *UserModel) GetSystemAvatarById(id int32) *models.DefaultAvatar {
	info := new(models.DefaultAvatar)
	ok, err := m.Engine.Where("id=?", id).Get(info)
	if !ok || err != nil {
		return nil
	}

	return info
}

// 搜索用户
func (m *UserModel) SearchUser(name string, offset, size int) []*UserSearchResults {
	sql := "SELECT * FROM user WHERE nick_name like '%" + name + "%' OR user_id like '%" + name + "%' ORDER BY `id` DESC LIMIT ?, ?"
	var list []*UserSearchResults
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		log.Log.Errorf("search_trace: search user err:%s", err)
		return nil
	}

	return list
}

// 获取用户总数
func (m *UserModel) GetUserTotalCount() int64 {
  count, err := m.Engine.Count(&models.User{})
  if err != nil {
    return 0
  }

  return count
}

// 设置注册ip
func (m *UserModel) SetRegisterIp(ip string) {
	m.User.RegIp = ip
}

// 设置uid
func (m *UserModel) SetUid(uid string) {
	m.User.UserId = uid
}

// 设置昵称
func (m *UserModel) SetNickName(name string) {
	m.User.NickName = name
}

// 设置手机号
func (m *UserModel) SetPhone(mobileNum int64) {
	m.User.MobileNum = mobileNum
}

// 设置用户头像
func (m *UserModel) SetAvatar(avatar string) {
	m.User.Avatar = avatar
}

// 设置用户类型
func (m *UserModel) SetUserDeviceType(utype int) {
	m.User.DeviceType = utype
}

// 设置设备token
func (m *UserModel) SetDeviceToken(deviceToken string) {
  m.User.DeviceToken = deviceToken
}

// 设备类型
func (m *UserModel) SetDeviceType(deviceType int) {
	m.User.DeviceType = deviceType
}

// 设置性别
func (m *UserModel) SetGender(gender int) {
	m.User.Gender = gender
}

// 设置用户状态
func (m *UserModel) SetStatus(status int) {
	m.User.Status = status
}

// 设置创建时间
func (m *UserModel) SetCreateAt(tm int64) {
	m.User.CreateAt= int(tm)
}

// 设置更新时间
func (m *UserModel) SetUpdateAt(tm int64) {
	m.User.UpdateAt = int(tm)
}

// 设置密码
func (m *UserModel) SetPassword(password string) {
	m.User.Password = password
}

// 设置用户类型 1 手机号 2 微信 3 QQ 4 微博
func (u *UserModel) SetUserType(userType int) {
	u.User.UserType = userType
}

// 设置登陆时间
func (u *UserModel) SetLastLoginTime(tm int64) {
	u.User.LastLoginTime = int(tm)
}

// 获取用户token
func (m *UserModel) GetUserToken(uid string) (string, error) {
	rds := dao.NewRedisDao()
	return rds.Get(rdskey.MakeKey(rdskey.USER_AUTH, uid))
}

// 保存token(1周过期)
func (m *UserModel) SaveUserToken(uid, token string) error {
	rds := dao.NewRedisDao()
	return rds.SETEX(rdskey.MakeKey(rdskey.USER_AUTH, uid), rdskey.KEY_EXPIRE_WEEK, token)
}

// 获取token
func (m *UserModel) GenUserToken(uid, password string) string {
	auth := uid + "_" + strings.ToLower(util.Md5String(uid+"|"+password))
	return auth
}

// 最终按加密规则token校验
func (m *UserModel) TokenValid(account, password, hashcode string) (b bool) {
	md := strings.ToLower(util.MD5(fmt.Sprintf("%s|%s", account, util.Md5String(password))))
	log.Log.Debugf("@@@@@hashcode:%s, md:%s", hashcode, md)
	return md == hashcode
}

