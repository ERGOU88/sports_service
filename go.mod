module sports_service/server

go 1.14

replace github.com/go-xorm/core v0.6.3 => xorm.io/core v0.6.3

require (
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/core v0.6.3
	github.com/go-xorm/xorm v0.7.9
	github.com/json-iterator/go v1.1.10
	github.com/parnurzeal/gorequest v0.2.16
	github.com/rs/xid v1.2.1
	github.com/spf13/viper v1.7.1
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.5.1
	github.com/zheng-ji/goSnowFlake v0.0.0-20180906112711-fc763800eec9
	go.uber.org/zap v1.16.0
	moul.io/http2curl v1.0.0 // indirect
)
