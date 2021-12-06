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

// 禁用管理员
func (svc *AdminModule) ForbidAdminUser(username string, status int) int {
  admin := svc.admin.FindAdminUserByName(username)
  if admin == nil {
    return errdef.ADMIN_NOT_EXISTS
  }

  admin.Status = status
  if _, err := svc.admin.UpdateAdminUser(admin); err != nil {
    return errdef.ERROR
  }

  return errdef.SUCCESS
}

// 添加菜单详情
func (svc *AdminModule) AddMenuDetail(menu *models.SystemMenu) int {
  if _, err := svc.admin.AddMenu(menu); err != nil {
    return errdef.ERROR
  }

  return errdef.SUCCESS
}

// 更新菜单详情
func (svc *AdminModule) UpdateMenuDetail(menu *models.SystemMenu) int {
  if _, err := svc.admin.UpdateMenu(menu); err != nil {
    return errdef.ERROR
  }

  return errdef.SUCCESS
}

func (svc *AdminModule) GetAdminDetail(username string) *models.SystemUser {
  return svc.admin.FindAdminUserByName(username)
}

// 获取管理员列表
func (svc *AdminModule) GetAdminList(page, size int) (int, []*models.SystemUser) {
  offset := (page - 1) * size
  list, err := svc.admin.GetAdminUserList(offset, size)
  if err != nil {
    return errdef.ERROR, nil
  }

  if len(list) == 0 {
    return errdef.SUCCESS, []*models.SystemUser{}
  }

  return errdef.SUCCESS, list
}

func (svc *AdminModule) GetMenuDetail(menuId string) *models.SystemMenu {
  ok, err := svc.admin.GetMenu(menuId)
  if !ok || err != nil {
    return nil
  }

  return svc.admin.Menu
}

// 获取菜单列表
func (svc *AdminModule) GetMenuList(page, size int) (int, []*models.SystemMenu) {
  offset := (page - 1) * size
  list, err := svc.admin.GetMenuList(offset, size)
  if err != nil {
    return errdef.ERROR, nil
  }

  if len(list) == 0 {
    return errdef.SUCCESS, []*models.SystemMenu{}
  }

  return errdef.SUCCESS, list
}

func (svc *AdminModule) GetRoleMenuList(roleId string) (int, []*models.SystemRoleMenu) {
  list, err := svc.admin.GetRoleMenu(roleId)
  if err != nil {
    return errdef.ERROR, nil
  }

  if len(list) == 0 {
    return errdef.SUCCESS, []*models.SystemRoleMenu{}
  }

  //res := make([]*models.SystemMenu, 0)
  //for _, item := range list {
  //  ok, err := svc.admin.GetMenu(fmt.Sprint(item.MenuId))
  //  if !ok || err != nil {
  //    continue
  //  }
  //
  //  res = append(res, svc.admin.Menu)
  //}

  return errdef.SUCCESS, list
}

// 管理员登陆 todo:rbac
func (svc *AdminModule) AdminLogin(params *madmin.AdminRegOrLoginParams) (int, string, []*models.SystemRoleMenu) {
  admin := svc.admin.FindAdminUserByName(params.UserName)
  if admin == nil {
    return errdef.ADMIN_NOT_EXISTS, "", nil
  }

  if admin.Status == 1 {
    return errdef.ADMIN_STATUS_FORBID, "", nil
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
  if err != nil {
    return errdef.ERROR, "", nil
  }

  if len(menus) == 0 {
    return errdef.SUCCESS, token, []*models.SystemRoleMenu{}
  }
  

  //res := make([]*models.SystemMenu, 0)
  //for _, item := range menus {
  //  ok, err := svc.admin.GetMenu(fmt.Sprint(item.MenuId))
  //  if !ok || err != nil {
  //    continue
  //  }
  //
  //  res = append(res, svc.admin.Menu)
  //}

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

func (svc *AdminModule) AddRoleMenuList(param *madmin.AddRoleMenuParam) int {
  list := make([]*models.SystemRoleMenu, 0)
  now := time.Now()
  for _, item := range param.Menus {
    exists, err := svc.admin.HasExistsRoleMenu(item.RoleId, item.MenuId)
    if exists || err != nil {
      continue
    }
    
    info := &models.SystemRoleMenu{
      RoleId: item.RoleId,
      MenuId: item.MenuId,
      RoleName: item.RoleName,
      CreateAt: now,
      UpdateAt: now,
    }
    
    list = append(list, info)
  }
  
  
  if _, err := svc.admin.AddRoleMenuList(list); err != nil {
    log.Log.Errorf("admin_trace: add role menu list fail, err:%s", err)
    return errdef.ERROR
  }

  return errdef.SUCCESS
}

func (svc *AdminModule) UpdateRoleMenuList(param *madmin.AddRoleMenuParam) int {
  for _, item := range param.Menus {
    if _, err := svc.admin.UpdateRoleMenu(item); err != nil {
      log.Log.Errorf("admin_trace: update role menu list fail, err:%s", err)
      return errdef.ERROR
    }
  }

  return errdef.SUCCESS
}

// 获取角色列表
func (svc *AdminModule) GetRoleList(page, size int) (int, []*models.SystemRole) {
  offset := (page - 1) * size
  roleList, err := svc.admin.GetRoleList(offset, size)
  if err != nil {
    return errdef.ERROR, nil
  }

  if len(roleList) == 0 {
    return errdef.SUCCESS, []*models.SystemRole{}
  }

  return errdef.SUCCESS, roleList
}

// 添加角色
func (svc *AdminModule) AddRole(role *models.SystemRole) int {
  if _, err := svc.admin.AddRole(role); err != nil {
    return errdef.ERROR
  }

  return errdef.SUCCESS
}
