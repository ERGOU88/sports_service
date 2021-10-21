package cadmin

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/go-xorm/xorm"
  "sports_service/server/dao"
  "sports_service/server/global/backend/errdef"
  "sports_service/server/global/backend/log"
  "sports_service/server/global/consts"
  "sports_service/server/middleware/jwt"
  "sports_service/server/models"
  "sports_service/server/models/madmin"
  "sports_service/server/models/mldap"
  "sports_service/server/util"
  "strings"
  "time"
)

type AdminModule struct {
  context     *gin.Context
  engine      *xorm.Session
  admin       *madmin.AdminModel
  ldap        *mldap.LdapService
}

func New(c *gin.Context) AdminModule {
  socket := dao.AppEngine.NewSession()
  defer socket.Close()
  return AdminModule{
    context: c,
    admin: madmin.NewAdminModel(socket),
    ldap: mldap.NewAdModel(),
    engine: socket,
  }
}

// 添加后台管理员 todo:测试使用
func (svc *AdminModule) AddAdminUser(params *models.SystemUser) int {
  if params.Username == "" || params.Password == "" {
    return errdef.INVALID_PARAMS
  }

  admin := svc.admin.FindAdminUserByName(params.Username)
  if admin != nil {
    return errdef.ADMIN_HAS_EXISTS
  }

  // 去掉空格 换行
  name := strings.Trim(params.Username, " ")
  name = strings.Replace(name, "\n", "", -1)
  // +盐
  params.Salt = util.GenSecret(util.CHAR_MODE, 8)
  params.Username = name
  params.Password = util.Md5String(fmt.Sprintf("%s%s", params.Password, params.Salt))
  now := time.Now()
  params.CreateAt = now
  params.UpdateAt = now
  if err := svc.admin.AddAdminUser(params); err != nil {
    log.Log.Errorf("user_trace: add admin fail, err:%s", err)
    return errdef.ADMIN_ADD_FAIL
  }

  return errdef.SUCCESS
}

// 更新后台管理员
func (svc *AdminModule) UpdateAdminUser(params *models.SystemUser) int {
  if params.Username == "" || params.Password == "" {
    return errdef.INVALID_PARAMS
  }

  admin := svc.admin.FindAdminUserByName(params.Username)
  if admin == nil {
    return errdef.ADMIN_NOT_EXISTS
  }

  // 去掉空格 换行
  name := strings.Trim(params.Username, " ")
  name = strings.Replace(name, "\n", "", -1)
  // +盐
  params.Salt = util.GenSecret(util.CHAR_MODE, 8)
  params.Username = name
  params.Password = util.Md5String(fmt.Sprintf("%s%s", params.Password, params.Salt))
  if _, err := svc.admin.UpdateAdminUser(params); err != nil {
    log.Log.Errorf("user_trace: update admin fail, err:%s", err)
    return errdef.ADMIN_UPDATE_FAIL
  }

  return errdef.SUCCESS
}

// 管理员登陆 todo:rbac
func (svc *AdminModule) AdminLogin(params *madmin.AdminRegOrLoginParams) (int, string, []*models.SystemRoleMenu) {
  admin := svc.admin.FindAdminUserByName(params.UserName)
  if admin == nil {
    return errdef.ADMIN_NOT_EXISTS, "", nil
  }

  pwd := util.Md5String(fmt.Sprintf("%s%s", params.Password, admin.Salt))
  if pwd != admin.Password {
    return errdef.ADMIN_PASSWORD_NOT_MATCH, "", nil
  }

  ok, err := svc.admin.GetRole(fmt.Sprint(admin.RoleId))
  if !ok || err != nil {
    return errdef.UNAUTHORIZED, "", nil
  }

  jwtInfo :=  make([]jwt.JwtInfo, 0)
  jwtInfo = append(jwtInfo, jwt.JwtInfo{Key: consts.USER_NAME, Val: params.UserName},
  jwt.JwtInfo{Key: consts.EXPIRE, Val: time.Now().Add(time.Hour * 1).Unix()},
  jwt.JwtInfo{Key: consts.ROLE_ID, Val: fmt.Sprint(svc.admin.Role.RoleId)},
  jwt.JwtInfo{Key: consts.ROLE_KEY, Val: svc.admin.Role.RoleKey},
  jwt.JwtInfo{Key: consts.ROLE_NAME, Val: svc.admin.Role.RoleName},
  jwt.JwtInfo{Key: consts.IDENTIFY, Val: fmt.Sprint(svc.admin.User.UserId)})

  token, err := jwt.GenerateJwt(svc.context, jwtInfo)
  if err != nil {
    return errdef.ERROR, "", nil
  }

  menus, err := svc.admin.GetRoleMenu(fmt.Sprint(admin.RoleId))
  if err != nil || len(menus) == 0 {
    return errdef.ERROR, "", nil
  }

  return errdef.SUCCESS, token, menus
}

// 域用户登录
func (svc *AdminModule) AdUserLogin(params *madmin.AdminRegOrLoginParams) int {
  if err := svc.ldap.CheckLogin(params.UserName, params.Password); err != nil {
    log.Log.Errorf("user_trace: check login err:%s", err)
    return errdef.ADMIN_PASSWORD_NOT_MATCH
  }

  return errdef.SUCCESS
}
