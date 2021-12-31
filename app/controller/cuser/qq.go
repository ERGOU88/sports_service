package cuser

import (
	"github.com/garyburd/redigo/redis"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	third "sports_service/server/tools/thirdLogin"
	"sports_service/server/util"
	"time"
)

func (svc *UserModule) QQLoginOrReg(params *muser.QQLoginParams) (int, string, *models.User) {
	qq := third.NewQQ()
	unionid, err := qq.GetQQUnionID(params.AccessToken)
	if err != nil || unionid == "" {
		log.Log.Errorf("qq_trace: get qq unionid err:%v, unionid:%s", err, unionid)
		return errdef.QQ_UNIONID_FAIL, "", nil
	}

	// 数据库获取社交平台帐号信息
	info := svc.social.GetSocialAccountByType(consts.TYPE_QQ, unionid)
	// 微信信息不存在 则注册
	if info == nil {
		client := util.GetClient(svc.context.Request.UserAgent())
		// 获取QQ用户信息
		qqinfo, err := qq.GetQQUserInfo(params.Openid, params.AccessToken, client, params.Platform)
		if qqinfo == nil || err != nil {
			log.Log.Error("qq_trace: get qq user info error")
			return errdef.QQ_USER_INFO_FAIL, "", nil
		}

		r := muser.NewQQRegister()
		gender := svc.QQInterChangeGender(qqinfo.Gender)
		// 注册
		if err := r.Register(svc.user, svc.social, svc.context, unionid, qqinfo.Figureurl, qqinfo.Nickname, "",
			consts.TYPE_QQ, gender); err != nil {
			log.Log.Errorf("qq_trace: register err:%s", err)
			return errdef.QQ_REGISTER_FAIL, "", nil
		}

		// 开启事务
		svc.engine.Begin()
		// 添加用户
		if err := svc.user.AddUser(); err != nil {
			log.Log.Errorf("qq_trace: add user err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_ADD_INFO_FAIL, "", nil
		}

		// 添加社交帐号（QQ）
		if err := svc.social.AddSocialAccountInfo(); err != nil {
			log.Log.Errorf("qq_trace: add qq account err:%s", err)
			svc.engine.Rollback()
			return errdef.QQ_ADD_ACCOUNT_FAIL, "", nil
		}

		// 添加通知初始化设置
		if err := svc.notify.AddUserNotifySetting(svc.user.User.UserId, int(time.Now().Unix())); err != nil {
			log.Log.Errorf("qq_trace: add user notify set err:%s", err)
			svc.engine.Rollback()
			return errdef.USER_ADD_NOTIFY_SET_FAIL, "", nil
		}

		// 提交事务
		if err := svc.engine.Commit(); err != nil {
			log.Log.Errorf("qq_trace: commit transaction err:%s", err)
			return errdef.ERROR, "", nil
		}

		// 生成token
		token := svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("qq_trace: save user token err:%s", err)
		}

		return errdef.SUCCESS, token, svc.user.User
	}

	// 通过uid查询用户信息
	user := svc.user.FindUserByUserid(info.UserId)
	if user == nil {
		log.Log.Error("qq_trace: find user by uid error")
		return errdef.USER_GET_INFO_FAIL, "", nil
	}

  // 登陆的时候 检查用户状态
  if !svc.CheckUserStatus(user.Status) {
    log.Log.Errorf("user_trace: forbid status, userId:%s", user.UserId)
    return errdef.USER_FORBID_STATUS, "", nil
  }

	// 用户已注册过, 则直接从redis中获取token并返回
	token, err := svc.user.GetUserToken(user.UserId)
	if err != nil && err == redis.ErrNil {
		// redis 没有，重新生成token
		token = svc.user.GenUserToken(user.UserId, util.Md5String(user.Password))
		// 重新保存到redis
		if err := svc.user.SaveUserToken(user.UserId, token); err != nil {
			log.Log.Errorf("qq_trace: save user token err:%s", err)
		}
	}

	return errdef.SUCCESS, token, svc.user.User
}

// QQ性别替换成app的性别标识
func (svc *UserModule) QQInterChangeGender(gender string) int {
	switch gender {
	case "男":
		return consts.BOY
	case "女":
		return consts.GIRL
	default:
		return consts.BOY_OR_GIRL
	}
}
