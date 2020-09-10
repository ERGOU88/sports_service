package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	//"fmt"
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

	engineGroup.ShowSQL(true)
	engineGroup.SetMaxIdleConns(150)
	engineGroup.SetMaxOpenConns(1000)
	return engineGroup, nil
}
