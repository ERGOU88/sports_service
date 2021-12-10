package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"io"
	"sports_service/server/global/consts"
	"time"
	"fmt"
	//"github.com/arthurkiller/rollingwriter"
	//"sports_service/server/app/config"
)

var (
	AppEngine   *xorm.EngineGroup
	VenueEngine *xorm.EngineGroup
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

func InitXorm(masterDb, logPath, mode string, slaveDb []string, maxIdle, maxActive int) *xorm.EngineGroup {
	var err error
	engine, err := ConnectDbs(masterDb, logPath, mode, slaveDb, maxIdle, maxActive)
	if err != nil {
		panic(err)
	}

	return engine
}

// ConnectDb 连接数据库，主从
func ConnectDbs(masterDsn, logPath, mode string, slaveDsn []string, maxIdle, maxActive int) (*xorm.EngineGroup, error) {
	conns := make([]string, len(slaveDsn)+1)
	conns[0] = masterDsn
	for i, v := range slaveDsn {
		conns[i+1] = v
	}

	engineGroup, err := xorm.NewEngineGroup("mysql", conns)
	if err != nil {
		return nil, err
	}


	engineGroup.SetLogger(xorm.NewSimpleLogger(getWriter(logPath)))
	if mode != string(consts.ModeProd) {
		engineGroup.ShowSQL(true)
	}

	engineGroup.SetMaxIdleConns(maxIdle)
	engineGroup.SetMaxOpenConns(maxActive)
	return engineGroup, nil
}

// 日志文件切割
func getWriter(filename string) io.Writer {
	// 保存30天内的日志，每24小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		fmt.Sprintf(filename, "%Y%m%d"),
		//rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour * 24 * 90),
		rotatelogs.WithRotationTime(time.Hour * 24),
	)

	if err != nil {
		panic(err)
	}

	return hook
}
