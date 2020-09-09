package main

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"sports_service/server/app/routers"
	"sports_service/server/app/config"
	"flag"
	"sports_service/server/global/app/log"
	"sports_service/server/log/zap"
)

var (
	configFile = flag.String("c", "config/app.yaml", "-c 配置文件")
)

func init() {
	flag.Parse()
	// todo 初始化配置、日志、mysql等
	if err := config.Global.Load(*configFile); err != nil {
		fmt.Printf("Load config error %v\n", err)
		return
	}

	if config.Global.Debug {
		fmt.Println(fmt.Sprintf("conf is %+v\n", config.Global))
	}

	// 初始化日志
	log.Log = zap.InitZapLog(config.Global.Log.Dir, config.Global.Log.ShowColor, config.Global.Log.Level)
	log.Log.Debug("success")

}

func main() {
	// 启动服务
	engine := gin.New()
	routers.InitRouter(engine)
	if err := engine.Run(config.Global.PublicAddr); err != nil {
		fmt.Printf("engine.Run err:%v", err)
		return
	}
}
