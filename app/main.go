package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"sports_service/server/app/config"
	"sports_service/server/app/routers"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/job"
	"sports_service/server/log/zap"
	"sports_service/server/models/mlabel"
	"sports_service/server/models/pprof"
	"sports_service/server/nsqlx"
	"sports_service/server/rabbitmq"
	"sports_service/server/redismq"
	"sports_service/server/tools/nsq"
	"sports_service/server/util"
	"syscall"
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
	dao.AppEngine = dao.InitXorm(config.Global.Mysql.Main.Master, config.Global.Mysql.Main.Slave)
	dao.VenueEngine = dao.InitXorm(config.Global.Mysql.Main.Master, config.Global.Mysql.Main.Slave)
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

	/*----检测管理后台定时推送 是否已发送----*/
	go job.CheckTimedNotify()
	/*----检测管理后台定时推送 是否已发送----*/
}

// 初始化nsq（生产者）
func setupNsqProduct() {
	nsq.ConnectNsqProduct(config.Global.NsqAddr)
}
// 初始化nsql消费者
func setupNsqConsumer() {
	go nsqlx.InitNsqConsumer()
}

// SIGHUP 终端控制进程结束(终端连接断开), 十进制：1
// SIGQUIT 用户发送QUIT字符(Ctrl+/)触发, 十进制值：3
// SIGTERM 结束程序(可以被捕获、阻塞或忽略), 十进制值：15
// SIGINT 用户发送INTR字符(Ctrl+C)触发, 十进制值：2
// SIGSTOP 停止进程(不能被捕获、阻塞或忽略), 十进制值：17,19,23
// register signals handler
func setupSignal() {
	log.Log.Info("优雅关闭web服务")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	go func() {
		select {
		case <-sigChan:
			// 停止消费
			//nsq.Stop()
			// 停止生产
			//nsq.NsqProducer.Stop()
			os.Exit(-1)
		}
	}()
}

// 初始化rabbitmq消费者
func setupRabbitmqConsumer() {
	go rabbitmq.InitRabbitmqConsumer()
}

// 初始化redis消息队列 [消费者]
func setupRedisMqConsumer() {
	go redismq.InitRedisMq()
}

func setupLabelList() {
	go mlabel.InitLabelList()
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
	// 初始化nsq生产者
	//setupNsqProduct()
	// 初始化nsq消费者
	//setupNsqConsumer()
	// 初始化rabbitmq消费者
	//setupRabbitmqConsumer()
	// 初始化redis消息队列 [消费者]
	setupRedisMqConsumer()
	// 初始化视频标签配置列表 [load到内存]
	setupLabelList()
	// register signals handler
	setupSignal()
	// 本地运行时 不执行定时任务
	if config.Global.Mode != string(consts.ModeLocal) {
		// 任务
		setupJob()
	}

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
// @description ### 六、推送通知文档
// @description 推送通知类型
// @description   ↓↓↓
// @description [点击查看](/api/v1/notify/doc)
// @description
// @description 透传（自定义推送）说明
// @description | 客户端 | 字段 | 类型 | 说明 | 示例 |
// @description | ------ | :----- | :----- | :----- | :----- |
// @description | android | custom | string | android透传数据（参考 自定义消息结构说明） | {'msg_id':'123456','data':{'unread_num':100},'msg_type':20000,'send_time':1607665370,'display':true}' |
// @description | iOS | extra | string | iOS透传数据（参考 自定义消息结构说明） | {'msg_id':'123456','data':{'unread_num':100},'msg_type':20000,'send_time':1607665370,'display':true} |
// @description
// @description 自定义消息结构说明
// @description | 字段 | 类型 | 说明 |
// @description | ------ | :----- | :----- |
// @description | msg_id | string | 消息id |
// @description | msg_type | int32 | 消息类型 |
// @description | send_time | int64 | 发送时间 |
// @description | display | bool | false不展示 true展示 |
// @description | data | map[string]interface{} | 业务场景透传数据[参考接口返回data] |
// @description
// @description 推送数据结构（iOS）
// @description
// @description {
// @description    'appkey':'5fc60a125a31bc1b28ee9cb7',
// @description    'timestamp':1607658719,
// @description    'type':'unicast',
// @description    'device_tokens':'3fa23e7d5afa9705f7a3d7161a64bb99531100bf818b2ea6d35a14d8dec6b6ce',
// @description    'payload':{
// @description    'aps':{
// @description       'alert':{
// @description       'title':'订单消息',
// @description       'subtitle':'订单消息',
// @description       'body':'订单快超时啦，赶紧付款～'
// @description      },
// @description      'sound':'default'
// @description    },
// @description    'extra':'{'msg_id':'123456','data':{'unread_num':100},'msg_send_time':1607658720,'display':false}'
// @description    },
// @description    'policy':{
// @description       'expire_time':'2020-12-15 15:52:00'
// @description    },
// @description    'production_mode':false,
// @description    'description':'风暴英雄'
// @description }
// @description
// @description 推送数据结构（android）
// @description
// @description {
// @description    'appkey':'5fabde2445b2b751a929c2d5',
// @description    'timestamp':1607665369,
// @description    'type':'unicast',
// @description    'device_tokens':'AqfCkpFEey_nfS13tOtF6epRVhOcA0Ny9vPelVd7dGl1',
// @description    'payload':{
// @description      'display_type':'message',
// @description      'body':{
// @description        'title':'订单通知',
// @description        'text':'付款倒计时10分钟',
// @description        'custom':'{'msg_id':'123456','data':{'unread_num':100},'msg_type':20000,'send_time':1607665370,'display':true}'
// @description      }
// @description    },
// @description    'policy':{},
// @description    'production_mode':false
// @description }
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

