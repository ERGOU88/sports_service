package main

import (
  "flag"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/go-xorm/xorm"
  "io/ioutil"
  "net/http/httptest"
  "regexp"
  "sports_service/server/dao"
  "sports_service/server/global/app/errdef"
  "sports_service/server/global/consts"
  "sports_service/server/models"
  "sports_service/server/models/mvideo"
  "sports_service/server/app/controller/cvideo"
  cloud "sports_service/server/tools/tencentCloud"
  "sports_service/server/util"
  "strings"
  "time"
)

var (
  server = flag.String("svr", "", "-svr 指定服务器(本地(local)/测试服(test)/qa服(qa)/自定义(custom))")
  masterDb = flag.String("m", "", "-m 主数据库地址")
  spiderDb = flag.String("s", "", "-s 爬虫数据库地址")
  rdshost = flag.String("r", "", "-r redis地址")
  pwd = flag.String("p", "", "-p redis密码")
  uid = flag.String("u", "", "-u 指定用户id")
  exec = flag.String("e", "", "-e upload 执行upload")
)

func main() {
  flag.Parse()
  now := time.Now().Format("2006-01-02 15:04:05")
  fmt.Printf("当前时间点:%s\n", now)
  switch *server {
  case "local":
    dao.Engine = dao.InitXorm("root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4", []string{"root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4"})
    dao.InitRedis("192.168.5.12:6378", "")
    engine2 = dao.InitXorm("root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4", []string{"root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4"})
  // 测试服
  case "test":
    dao.Engine = dao.InitXorm(fmt.Sprintf("%s", *masterDb), nil)
    dao.InitRedis(*rdshost, *pwd)
    engine2 = dao.InitXorm(fmt.Sprintf("%s", *spiderDb), nil)
  // qa服
  case "qa":
    if *rdshost == "" {
      dao.InitRedis("192.168.5.12:6378", "")
    } else {
      dao.InitRedis(*rdshost, *pwd)
    }

    if *masterDb != "" {
      dao.Engine = dao.InitXorm(fmt.Sprintf("%s", *masterDb), nil)
    } else {
      dao.Engine = dao.InitXorm("root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4", []string{"root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4"})
    }

    if *spiderDb != "" {
      engine2 = dao.InitXorm(fmt.Sprintf("%s", *spiderDb), nil)
    } else {
      engine2 = dao.InitXorm("root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4", []string{"root:bluetrans888@tcp(192.168.5.12:3306)/sports_service?charset=utf8mb4"})
    }

  // 自定义
  case "custom":
    dao.Engine = dao.InitXorm(fmt.Sprintf("%s", *masterDb), nil)
    dao.InitRedis(*rdshost, *pwd)
    engine2 = dao.InitXorm(fmt.Sprintf("%s", *spiderDb), nil)
  default:
    fmt.Printf("unsupport svr flag:%s\nUSAGE: -svr=local | test | online | custom\n", *server)
    return
  }

  if *exec != "upload" {
    fmt.Printf("unsupport exec flag:%s\nUSAGE: -e=upload\n", *exec)
    return
  }


  var uids []string
  if *uid == "" {
    uids = GetUserIds()
    if len(uids) == 0 {
      fmt.Printf("unsupport uid:%s\nUSAGE: -u=[user_id]\n", *uid)
      return
    }
  }

  listFile := ReadDirInfo()
  if len(listFile) == 0 {
    fmt.Println("mp4 file not found")
    return
  }

  //spiderInfo := GetSpiderInfo()
  //if len(spiderInfo) == 0 {
  //  fmt.Println("spider info not found")
  //  return
  //}

  for _, file := range listFile {
    //randNum := util.GenerateRandnum(0, len(spiderInfo) - 1)
    randIndex := util.GenerateRandnum(0, len(uids) - 1)
    if *uid == "" {
      *uid = uids[randIndex]
    }

    c, _ := gin.CreateTestContext(httptest.NewRecorder())
    svc := cvideo.New(c)
    syscode, _, taskId := svc.GetUploadSign(*uid)
    if syscode != errdef.SUCCESS {
      fmt.Printf("\nsyscode:%d", syscode)
      continue
    }

    client := cloud.New(consts.TX_CLOUD_SECRET_ID, consts.TX_CLOUD_SECRET_KEY, consts.VOD_API_DOMAIN)
    resp, err := client.Upload(taskId, *uid, "", fmt.Sprintf("./%s", file),
      "ap-shanghai", consts.VOD_PROCEDURE_NAME)
    if err != nil {
      fmt.Printf("\nupload file err:%s", err)
      continue
    }

    params := new(mvideo.VideoPublishParams)

    info := GetSpiderInfoByFileName(file)
    fmt.Printf("info:%v", info)
    if info != nil {
      params.Title = TrimHtml(info.Title)
      fmt.Println("\ntitle:", params.Title)
      params.Describe = TrimHtml(info.Description)
      //params.Cover = info.Pic
    }

    params.FileId = *resp.Response.FileId
    params.VideoAddr = *resp.Response.MediaUrl
    if *resp.Response.CoverUrl != "" {
      params.Cover = *resp.Response.CoverUrl
    }



    params.VideoLabels = "1,2"
    params.TaskId = taskId
    fmt.Printf("\n taskId:%d\n", taskId)
    if syscode := svc.RecordPubVideoInfo(*uid, params); syscode != errdef.SUCCESS {
      fmt.Printf("\nuser publish video err:%s", err)
      continue
    }

    fmt.Println("\nupload success")
  }

  return
}

type Bili struct {
  Title        string   `json:"title"`
  Description  string   `json:"description"`
  Pic          string   `json:"pic"`
}

var engine2  *xorm.EngineGroup
func GetSpiderInfo() []*Bili {
  var list []*Bili
  if err := engine2.Find(&list); err != nil {
    fmt.Printf("get spider info err:%s", err)
    return nil
  }

  return list
}

func TrimHtml(src string) string {
  // 将HTML标签全转换成小写
  re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
  src = re.ReplaceAllStringFunc(src, strings.ToLower)

  // 去除STYLE
  re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
  src = re.ReplaceAllString(src, "")

  // 去除SCRIPT
  re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
  src = re.ReplaceAllString(src, "")

  // 去除所有尖括号内的HTML代码，并换成换行符
  re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
  src = re.ReplaceAllString(src, "\n")

  // 去除连续的换行符
  re, _ = regexp.Compile("\\s{2,}")
  src = re.ReplaceAllString(src, "\n")

  return strings.TrimSpace(src)
}

// 通过文件名获取爬虫数据
func GetSpiderInfoByFileName(fileName string) *Bili {
  info := new(Bili)
  sql := "SELECT * FROM bili WHERE video_url like '%" + fileName + "' LIMIT 1"
  ok, err := engine2.SQL(sql).Get(info)
  if !ok || err != nil {
    return nil
  }

  return info
}

func ReadDirInfo() []string {
  dir, err := ioutil.ReadDir("./")
  if err != nil {
    return []string{}
  }

  var listFile []string
  for _, file := range dir{
    ok := strings.HasSuffix(file.Name(), ".mp4")
    if !ok {
      continue
    }


    listFile = append(listFile, file.Name())
  }

  return listFile

}

// 获取用户id列表
func GetUserIds() []string {
  var uids []string
  if err := dao.Engine.Table(&models.User{}).Cols("user_id").Find(&uids); err != nil {
    return nil
  }

  return uids
}

