module sports_service/server

replace github.com/go-xorm/core v0.6.3 => xorm.io/core v0.6.3

go 1.14

require (
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-xorm/core v0.6.3
	github.com/json-iterator/go v1.1.9
	github.com/rs/xid v1.2.1
	github.com/spf13/viper v1.7.1
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.5.1
	github.com/zheng-ji/goSnowFlake v0.0.0-20180906112711-fc763800eec9
	go.uber.org/zap v1.16.0
)
