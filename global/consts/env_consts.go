package consts

type Mode string

const (
	ModeLocal Mode     = "local"   // 本地
	ModeDev   Mode     = "dev"     // 开发模式
	ModeTest  Mode     = "test"    // 测试模式
	ModeProd  Mode     = "prod"    // 生产模式
)

type AppId string

const (
	WEB_APP_ID       AppId =  "aunGaE4h"
	IOS_APP_ID       AppId =  "5EewXD1i"
	AND_APP_ID       AppId =  "mj4mQaop"
)

const (
  FORMAT_INFO     = "2006-01-02 15:04"
  FORMAT_DATE     = "2006-01-02"
  FORMAT_TM       = "2006-01-02 15:04:05"
  FORMAT_WX_TM    = "20060102150405"
  FORMAT_DATE_STR = "01月02日 15:04"
  FORMAT_MONTH    = "2006-01"
)

