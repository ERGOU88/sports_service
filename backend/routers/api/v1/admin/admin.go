package admin

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "io"
  "net/http"
  "os"
  "path"
  "sports_service/server/backend/config"
  "sports_service/server/backend/controller/cadmin"
  "sports_service/server/global/backend/errdef"
  "sports_service/server/global/backend/log"
  "sports_service/server/global/consts"
  "sports_service/server/middleware/jwt"
  "sports_service/server/models"
  "sports_service/server/models/madmin"
  "sports_service/server/util"
  "time"
)

// 注册后台用户
func RegAdminUser(c *gin.Context) {
  reply := errdef.New(c)
  params := new(models.SystemUser)
  if err := c.BindJSON(params); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  syscode := svc.AddAdminUser(params)
  reply.Response(http.StatusOK, syscode)
}

func AdLogin(c *gin.Context) {
  reply := errdef.New(c)
  params := new(madmin.AdminRegOrLoginParams)
  if err := c.BindJSON(params); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  //syscode := svc.AdminLogin(params)
  syscode := svc.AdUserLogin(params)
  reply.Response(http.StatusOK, syscode)
}


const UPLOAD_PRE_FIX = "./upload"
func UploadFile(c *gin.Context) {
  reply := errdef.New(c)
  file, fh, err := c.Request.FormFile("file_up")
  if file == nil || err != nil {
    log.Log.Errorf("upload_trace: receive form file err:%s", err)
    reply.Response(http.StatusBadRequest, errdef.INVALID_PARAMS)
    return
  }

  // 回绕指针
  if _, err := file.Seek(0, 0); err != nil {
    log.Log.Errorf("upload_trace: file seek err:%s", err)
    reply.Response(http.StatusBadRequest, errdef.ERROR)
    return
  }

  shortPath := ""
  var pathName string

  ext := path.Ext(fh.Filename)

  dateDir := time.Now().Format("2006_01_02") + "/"
  // 如果文件夹不存在则创建
  if _, err := os.Stat(UPLOAD_PRE_FIX + "/" + dateDir); os.IsNotExist(err) {
    os.MkdirAll(UPLOAD_PRE_FIX + "/" + dateDir, 0666)
  }

  shortPath = fmt.Sprintf("%s%d%s", dateDir, util.GetSnowId(), ext)
  pathName = UPLOAD_PRE_FIX + "/" + shortPath
  f, err := os.OpenFile(pathName, os.O_WRONLY|os.O_CREATE, 0666)
  if err != nil {
    log.Log.Errorf("upload_trace: open file err:%s", err)
    reply.Response(http.StatusOK, errdef.ERROR)
    return
  }
  defer f.Close()

  if _, err = io.Copy(f, file); err != nil {
    log.Log.Errorf("upload_trace: io copy err:%s", err)
    reply.Response(http.StatusOK, errdef.ERROR)
    return
  }

  reply.Data["file_url"] = config.Global.FileAddr + shortPath

  reply.Response(http.StatusOK, errdef.SUCCESS)
}

// 后台管理员登陆
func AdminLogin(c *gin.Context) {
  reply := errdef.New(c)
  params := new(madmin.AdminRegOrLoginParams)
  if err := c.BindJSON(params); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  syscode, token, menus := svc.AdminLogin(params)
  if syscode != errdef.SUCCESS {
    reply.Response(http.StatusOK, syscode)
    return
  }

  reply.Data["menus"] = menus
  reply.Data["token"] = token
  reply.Response(http.StatusOK, syscode)
}

func AddAdmin(c *gin.Context) {
  reply := errdef.New(c)
  param := &models.SystemUser{}
  if err := c.BindJSON(param); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  param.CreateBy, _ = util.StringToInt(jwt.GetUserInfo(c, consts.IDENTIFY))
  svc := cadmin.New(c)
  reply.Response(http.StatusOK, svc.AddAdminUser(param))
}

func UpdateAdmin(c *gin.Context) {
  reply := errdef.New(c)
  param := &models.SystemUser{}
  if err := c.BindJSON(param); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  param.UpdateBy, _ = util.StringToInt(jwt.GetUserInfo(c, consts.IDENTIFY))
  svc := cadmin.New(c)
  reply.Response(http.StatusOK, svc.UpdateAdminUser(param))
}

func ForbidAdmin(c *gin.Context) {
  reply := errdef.New(c)
  param := &madmin.ForbidAdmin{}
  if err := c.BindJSON(param); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  reply.Response(http.StatusOK, svc.ForbidAdminUser(param.UserName, param.Status))
}

func GetRoleMenu(c *gin.Context) {
  reply := errdef.New(c)
  roleId := c.Query("role_id")
  svc := cadmin.New(c)
  code, list := svc.GetRoleMenuList(roleId)
  reply.Data["list"] = list
  reply.Response(http.StatusOK, code)
}

func AddRoleMenu(c *gin.Context) {
  reply := errdef.New(c)
  param := &madmin.AddRoleMenuParam{}
  if err := c.BindJSON(param); err != nil {
    log.Log.Errorf("admin_trace: invalid param, param:%+v, err:%s", param, err)
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  code := svc.AddRoleMenuList(param)
  reply.Response(http.StatusOK, code)
}

func UpdateRoleMenu(c *gin.Context) {
  reply := errdef.New(c)
  param := &madmin.AddRoleMenuParam{}
  if err := c.BindJSON(param); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  code := svc.UpdateRoleMenuList(param)
  reply.Response(http.StatusOK, code)
}

func AdminDetail(c *gin.Context) {
  reply := errdef.New(c)
  username := c.Query("username")

  svc := cadmin.New(c)
  admin := svc.GetAdminDetail(username)
  reply.Data["detail"] = admin
  reply.Response(http.StatusOK, errdef.SUCCESS)
}

func AdminList(c *gin.Context) {
  reply := errdef.New(c)
  page, size := util.PageInfo(c.Query("page"), c.Query("size"))
  svc := cadmin.New(c)
  code, list := svc.GetAdminList(page, size)
  reply.Data["list"] = list
  reply.Response(http.StatusOK, code)
}

func MenuDetail(c *gin.Context) {
  reply := errdef.New(c)
  menuId := c.Query("menu_id")
  svc := cadmin.New(c)
  menu := svc.GetMenuDetail(menuId)
  reply.Data["detail"] = menu
  reply.Response(http.StatusOK, errdef.SUCCESS)
}

func MenuList(c *gin.Context) {
  reply := errdef.New(c)
  page, size := util.PageInfo(c.Query("page"), c.Query("size"))
  svc := cadmin.New(c)
  code, list := svc.GetMenuList(page, size)
  reply.Data["list"] = list
  reply.Response(http.StatusOK, code)
}

func AddMenu(c *gin.Context) {
  reply := errdef.New(c)
  param := &models.SystemMenu{}
  if err := c.Bind(param); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  code := svc.AddMenuDetail(param)
  reply.Response(http.StatusOK, code)
}

func UpdateMenu(c *gin.Context) {
  reply := errdef.New(c)
  param := &models.SystemMenu{}
  if err := c.Bind(param); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  code := svc.UpdateMenuDetail(param)
  reply.Response(http.StatusOK, code)
}

func RoleList(c *gin.Context) {
  reply := errdef.New(c)
  page, size := util.PageInfo(c.Query("page"), c.Query("size"))
  svc := cadmin.New(c)
  code, list := svc.GetRoleList(page, size)
  reply.Data["list"] = list
  reply.Response(http.StatusOK, code)
}

func AddRole(c *gin.Context) {
  reply := errdef.New(c)
  param := &models.SystemRole{}
  if err := c.BindJSON(param); err != nil {
    reply.Response(http.StatusOK, errdef.INVALID_PARAMS)
    return
  }

  svc := cadmin.New(c)
  reply.Response(http.StatusOK, svc.AddRole(param))
}
