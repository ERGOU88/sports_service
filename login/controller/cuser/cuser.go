package cuser

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/login/errdef"
	"sports_service/server/global/login/log"
	"sports_service/server/models"
	"sports_service/server/models/muser"
	"sports_service/server/util"
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

// 手机一键登陆/注册
func (svc *UserModule) MobileLoginOrReg(param *muser.LoginParams) (int, string, *models.User) {
	if b := svc.user.CheckCellPhoneNumber(param.MobileNum); !b {
		log.Log.Errorf("user_trace: invalid mobile num %v", param.MobileNum)
		return errdef.INVALID_MOBILE_NUM, "", nil
	}

	// 根据手机号查询用户 不存在 注册用户 用户存在 为登陆
	if user := svc.user.FindUserByPhone(param.MobileNum); user == nil {
		// 登陆
		reg := muser.NewMobileRegister()
		if err := reg.Register(svc.user, param); err != nil {
			log.Log.Errorf("user_trace: register err:%s", err)
			return errdef.USER_REGISTER_FAIL, "", nil
		}

		// 添加用户
		if err := svc.user.AddUser(); err != nil {
			log.Log.Errorf("reg_trace: add user info err:%s", err)
			return errdef.USER_ADD_INFO_FAIL, "", nil
		}

		// 生成token
		token := svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("user_trace: save user token err:%s", err)
		}

		return errdef.SUCCESS, token, svc.user.User

	}

	// 用户已注册过, 则直接从redis中获取token并返回
	token, err := svc.user.GetUserToken(svc.user.User.UserId)
	if err != nil && err == redis.ErrNil {
		// redis 没有，重新生成token
		token = svc.user.GenUserToken(svc.user.User.UserId, util.Md5String(svc.user.User.Password))
		// 重新保存到redis
		if err := svc.user.SaveUserToken(svc.user.User.UserId, token); err != nil {
			log.Log.Errorf("user_trace: save user token err:%s", err)
		}
	}

	return errdef.SUCCESS, token, svc.user.User
}


