package job

import (
	"sports_service/server/global/app/log"
	"time"
)

// 检测banner列表 是否上架/是否过期（每5分钟）
func CheckOrder() {
	ticker := time.NewTicker(time.Minute * 3)
	defer ticker.Stop()

	for {
		select {
		case <- ticker.C:
			log.Log.Debugf("开始检测banner 是否上架/是否已过期")
			// 检测banner是否上架/是否已过期
			checkBannerStatus()
			log.Log.Debugf("检测完毕")
		}
	}

}
