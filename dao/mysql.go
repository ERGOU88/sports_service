package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	//"fmt"
	//"github.com/arthurkiller/rollingwriter"
)

var (
	Engine *xorm.EngineGroup
)

func ConnectDb(mysqlDsn string) *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", mysqlDsn)
	if err != nil {
		panic(err)
	}

	engine.ShowSQL(true)
	return engine
}

// TODO，为了避免主从同步速度问题， 现强制使用主库
// 从库使用需要手动指定
func MasterOnly() xorm.GroupPolicyHandler {
	return func(g *xorm.EngineGroup) *xorm.Engine {
		return g.Master()
	}
}

func InitXorm(masterDb string, slaveDb []string) *xorm.EngineGroup {
	var err error
	engine, err := ConnectDbs(masterDb, slaveDb)
	if err != nil {
		panic(err)
	}

	return engine
}

// ConnectDb 连接数据库，主从
func ConnectDbs(masterDsn string, slaveDsn []string) (*xorm.EngineGroup, error) {
	conns := make([]string, len(slaveDsn)+1)
	conns[0] = masterDsn
	for i, v := range slaveDsn {
		conns[i+1] = v
	}

	engineGroup, err := xorm.NewEngineGroup("mysql", conns)
	if err != nil {
		return nil, err
	}

	//config := rollingwriter.Config{
	//	LogPath:       "./logs",                    // 日志路径
	//	TimeTagFormat: "060102150405",              // 时间格式串
	//	FileName:      "mysql_exec",                // 日志文件名
	//	MaxRemain:     0,                           // 配置日志最大存留数
	//	RollingPolicy: rollingwriter.VolumeRolling, // 配置滚动策略
	//	RollingTimePattern: "* * * * * *", // 配置时间滚动策略
	//	RollingVolumeSize:  "1M",          // 配置截断文件下限大小
	//	WriterMode: "none",
	//	BufferWriterThershould: 256,
	//	Compress: true, // Compress will compress log file with gzip
	//}
	//
	//writer, err := rollingwriter.NewWriterFromConfig(&config)
	//if err != nil {
	//	panic(err)
	//}
	//engineGroup.SetLogger(xorm.NewSimpleLogger(writer))

	engineGroup.ShowSQL(true)
	engineGroup.SetMaxIdleConns(150)
	engineGroup.SetMaxOpenConns(1000)
	return engineGroup, nil
}
