package job

import (
	"sports_service/server/global/consts"
	"sports_service/server/models"
	"sports_service/server/global/app/log"
	"time"
	"sports_service/server/dao"
)

// 检测banner列表 是否上架/是否过期（每5分钟）
func CheckBanners() {
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

// 执行检测任务 检测banner 是否上架/是否已过期
func checkBannerStatus() {
	list := getAllBanner()
	if list == nil {
		return
	}

	now := int(time.Now().Unix())
	for _, banner := range list {
		if banner.StartTime > now {
			continue
		}
		//  展示开始时间 <= 当前时间 且 展示结束时间 > 当前时间 状态设置为上架
		if banner.StartTime <= now && banner.EndTime > now {
			banner.Status = consts.HAS_LAUNCHED
		}

		// 展示结束时间 <= 当前时间 状态设置为已过期
		if banner.EndTime <= now {
			banner.Status = consts.NO_LAUNCHED
		}

		// 更新banner状态
		updateBannerStatus(banner)
	}
}

// 获取所有banner
func getAllBanner() []*models.Banner {
	var list []*models.Banner
	if err := dao.AppEngine.Find(&list); err != nil {
		return nil
	}

	return list
}

// 更新banner状态
func updateBannerStatus(banner *models.Banner) error {
	if _, err := dao.AppEngine.Where("id=?", banner.Id).Update(banner); err != nil {
		log.Log.Errorf("job_trace: update banner status err:%v", err)
		return err
	}

	return nil
}

