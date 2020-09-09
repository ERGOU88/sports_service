package config

import (
	"fmt"
	"github.com/spf13/viper"
	"errors"
)

type Address struct {
	Ip   string "ip"
	Port int    "port"
}

type Addresses []Address

type MysqlType struct {
	Master    string   "master"
	Slave     []string "slave"
	MaxIdle   int      "max_idle"
	MaxActive int      "max_active"
}

type RedisOption struct {
	DBId      int    "dbid"
	Name      string "name"
	MaxIdle   int    "max_idle"
	MaxActive int    "max_active"
}

type RedisType struct {
	Master  Address       "master"
	Slave   Addresses     "slave"
	Options []RedisOption "dboption"
}

type Config struct {
	// 公网服务
	PublicAddr string "public_addr"

	// 日志
	Log struct {
		Dir       string "dir"
		Level     int    "level"
		ShowColor bool   "show_color"
	} "log"

	// mysql
	Mysql struct {
		Main MysqlType "main"
	} "mysql"

	// redis连接
	Redis struct {
		// 主服务
		Main RedisType "main"
	} "redis"

	Debug bool "debug"
}

var Global Config

func (c *Config) Load(confFile string) error {
	// 实例
	v := viper.New()
	// 设置完成配置文件
	v.SetConfigFile(confFile)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("no such config file")
		}

		return err
	}

	if err := v.Unmarshal(&Global); err != nil {
		fmt.Printf("err:%s",err)
		return err
	}

	return nil
}


