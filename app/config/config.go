package config

import (
	"fmt"
	"github.com/spf13/viper"
	"errors"
)

type Address struct {
	Ip   string
	Port int
}

type Addresses []Address

type MysqlType struct {
	Master    string
	Slave     []string
	MaxIdle   int
	MaxActive int
}

type RedisOption struct {
	DBId      int
	Name      string
	MaxIdle   int
	MaxActive int
}

type RedisType struct {
	Master  Address
	Slave   Addresses
	Options []RedisOption
}

type Config struct {
	// 公网服务
	PublicAddr string
	// 性能监控地址
	PprofAddr  string
	// 开发（dev）、测试(test)、生产(prod)
	Mode       string
	// 日志
	Log struct {
		Path      string
		Level     int
		ShowColor bool
	}

	// mysql
	Mysql struct {
		Main MysqlType
	}

	// redis连接
	Redis struct {
		// 主服务
		Main RedisType
	}

	Debug bool
}

var Global Config

func (c *Config) Load(confFile string) error {
	// 实例
	v := viper.New()
	// 设置配置文件
	v.SetConfigFile(confFile)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("no such config file")
		}

		return err
	}

	if err := v.Unmarshal(&Global); err != nil {
		fmt.Printf("unmarshal err:%s",err)
		return err
	}

	return nil
}


