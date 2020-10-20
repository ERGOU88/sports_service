package cadmin

import (
  "github.com/gin-gonic/gin"
  "github.com/go-xorm/xorm"
  "sports_service/server/dao"
  "sports_service/server/global/backend/errdef"
  "sports_service/server/models/madmin"
  "sports_service/server/util"
  "strings"
  "time"
  "fmt"
)

type AdminModule struct {
  context     *gin.Context
  engine      *xorm.Session
  admin       *madmin.AdminModel
}

func New(c *gin.Context) AdminModule {
  socket := dao.Engine.Context(c)
  defer socket.Close()
  return AdminModule{
    context: c,
    admin: madmin.NewAdminModel(socket),
    engine: socket,
  }
}

// 添加后台管理员 todo:测试使用
func (svc *AdminModule) AddAdminUser(params *madmin.AdminRegOrLoginParams) int {
  admin := svc.admin.FindAdminUserByName(params.UserName)
  if admin != nil {
    return errdef.ADMIN_NOT_EXISTS
  }

  // 去掉空格 换行
  name := strings.Trim(params.UserName, " ")
  name = strings.Replace(name, "\n", "", -1)
  now := time.Now().Unix()
  // +盐
  svc.admin.User.Salt = util.GenSecret(util.CHAR_MODE, 8)
  svc.admin.User.Username = name
  svc.admin.User.Password = util.Md5String(fmt.Sprintf("%s%s", params.Password, svc.admin.User.Salt))
  // 0为主账号
  svc.admin.User.SubAccount = 0
  svc.admin.User.CreateAt = int(now)
  svc.admin.User.UpdateAt = int(now)
  if err := svc.admin.AddAdminUser(); err != nil {
    return errdef.ADMIN_ADD_FAIL
  }

  return errdef.SUCCESS
}

// 管理员登陆 todo:rbac
func (svc *AdminModule) AdminLogin(params *madmin.AdminRegOrLoginParams) int {
  admin := svc.admin.FindAdminUserByName(params.UserName)
  if admin == nil {
    return errdef.ADMIN_NOT_EXISTS
  }

  pwd := util.Md5String(fmt.Sprintf("%s%s", params.Password, admin.Salt))
  if pwd != admin.Password {
    return errdef.ADMIN_PASSWORD_NOT_MATCH
  }

  return errdef.SUCCESS
}
