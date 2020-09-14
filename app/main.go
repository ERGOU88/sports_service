package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/log/zap"
	"sports_service/server/app/config"
	"sports_service/server/app/routers"
	"sports_service/server/models/pprof"
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
	dao.InitRedis(rdshost, "")
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
}

// @title FPV电竞APP（应用服）
// @version 1.0
// @description ### 一、公共参数说明（此栏参数均为Headers请求头传递）
// @description | 参数名 | 说明 | 示例 |
// @description | ------ | :----- | :----- |
// @description | User-Agent | 用户代理 | android |
// @description | Version | 当前版本 | 1.0.1 |
// @description ### 二、请求体说明（此栏参数均为POST JSON传递，不可用form-data提交）
// @description
// @description     {
// @description         'mobileNum': '13177656222',
// @description         'platform': 0
// @description     }
// @description ### 三、API错误码文档
// @description [点击查看](/api/v1/doc)
// @description ### 四、HTTP状态码说明
// @description | 状态码 | 说明 |
// @description | ------ | :----- |
// @description | 200 | 操作成功 |
// @description | 400 | 参数错误 |
// @description | 500 | 内部错误 |
func main() {
	// 启动服务
	engine := gin.New()
	routers.InitRouters(engine)
	if err := engine.Run(config.Global.PublicAddr); err != nil {
		fmt.Printf("engine.Run err:%v", err)
		return
	}
}
