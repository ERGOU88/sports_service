package main

import (
  "flag"
  "fmt"
  "github.com/gin-gonic/gin"
  "sports_service/server/app/config"
  "sports_service/server/app/routers"
  "sports_service/server/dao"
  "sports_service/server/global/app/log"
  "sports_service/server/global/consts"
  "sports_service/server/job"
  "sports_service/server/log/zap"
  "sports_service/server/models/pprof"
  "sports_service/server/tools/nsq"
  "sports_service/server/util"
)

var (
	configFile = flag.String("c", "config/app_dev.yaml", "-c 配置文件")
)

// 配置
func setupConfig() error {
	flag.Parse()
	if err := config.Global.Load(*configFile); err != nil {
		fmt.Printf("Load config error %v\n", err)
		return err
	}

	if config.Global.Debug {
		fmt.Println(fmt.Sprintf("conf is %+v\n", config.Global))
	}

	return nil
}

// 日志
func setupLogger() {
	// 初始化日志
	log.Log = zap.InitZapLog(config.Global.Log.Path, config.Global.Log.ShowColor, config.Global.Log.Level)
	log.Log.Debug("setup log success")
}

// 初始化mysql
func setupMysql() {
	dao.Engine = dao.InitXorm(config.Global.Mysql.Main.Master, config.Global.Mysql.Main.Slave)
}

// 初始化redis
func setupRedis() {
	rdshost := fmt.Sprintf("%s:%d", config.Global.Redis.Main.Master.Ip, config.Global.Redis.Main.Master.Port)
	dao.InitRedis(rdshost, config.Global.RedisPassword)
}

// 性能监控
func setupPprof() {
	pprof.Start(config.Global.PprofAddr)
}

// 设置模式
func setupRunMode() {
	gin.SetMode(gin.DebugMode)
	if config.Global.Mode == string(consts.ModeTest) {
		gin.SetMode(gin.TestMode)
	}

	if config.Global.Mode == string(consts.ModeProd) {
		gin.SetMode(gin.ReleaseMode)
	}
}

// snow id
func setupSnowId() {
	util.InitSnowId()
}

// 任务列表
func setupJob() {
	/*----检测banner(是否上架/是否过期)任务----*/
	go job.CheckBanners()
	/*----检测banner(是否上架/是否过期)任务----*/

	/*----主动拉取腾讯云回调事件任务----*/
	go job.PullEventsJob()
	/*----主动拉取腾讯云回调事件任务----*/
}

// 初始化nsq（生产者）
func setupNsq() {
	nsq.ConnectNsqProduct(config.Global.NsqAddr)
}

func init() {
	// 配置
	if err := setupConfig(); err != nil {
		panic(err)
	}

	// 日志
	setupLogger()
	// mysql
	setupMysql()
	// redis
	setupRedis()
	// 性能监控
	setupPprof()
	// snow id
	setupSnowId()
	// 设置运行模式
	setupRunMode()
	// 任务
	setupJob()
	// 初始化nsq
	setupNsq()
}

// @title 电竞社区平台（应用服）
// @version 1.0
// @host fpv-web.youzu.com
// @schemes https
// @description ### 一、公共参数说明（此栏参数均为Headers请求头传递）
// @description | 参数名 | 说明 | 示例 |
// @description | ------ | :----- | :----- |
// @description | AppId | AppId(区分Android,iOS,web) | 5EesXF1i |
// @description | Secret | 服务端下发的secret 通过/api/v1/client/init接口获取 调用该接口时无需传 且 不参与签名 | DnaukFwVILpcewX6 |
// @description | Timestamp | 请求时间戳 单位：秒 | 1588888888 |
// @description | Sign | 签名 | 签名 md5签名32位值 |
// @description | Version | 当前版本 | 1.0.1 |
// @description ### 二、请求体说明（此栏参数均为POST JSON传递，不可用form-data提交）
// @description
// @description     {
// @description         'mobileNum': '13177656222',
// @description         'platform': 0
// @description     }
// @description ### 三、接口签名生成方式
// @description 签名加密示例:
// @description params = 请求的url路径(不包含域名与参数) + & + Header头参数以`&`拼接（无需按照字典序，具体看以下栗子） + & + appKey
// @description appKey由服务端下发 并进行保存
// @description sign = md5(params) 取md5 32位小写
// @description 如：md5(/api/v1/user/mobile/login&AppId=5EesXF1i&Timestamp=1588888888&Version=1.0.1&Secret=DnaukFwVILpcewX6&RfhHecN9zsNcy19Y)
// @description appKey为RfhHecN9zsNcy19Y
// @description ### 四、API错误码文档
// @description [点击查看](/api/v1/doc)
// @description ### 五、HTTP状态码说明
// @description | 状态码 | 说明 |
// @description | ------ | :----- |
// @description | 200 | 操作成功 |
// @description | 400 | 参数错误 |
// @description | 500 | 内部错误 |
func main() {
	// 启动服务
	engine := gin.New()
  engine.Static("/static", "./static")
	routers.InitRouters(engine)
	if err := engine.Run(config.Global.PublicAddr); err != nil {
		fmt.Printf("engine.Run err:%v", err)
		return
	}

}
