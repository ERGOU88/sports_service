package job

import (
	"sports_service/dao"
	"sports_service/global/app/log"
	"sports_service/global/rdskey"
	third "sports_service/tools/thirdLogin"
	"time"
)

// 刷新小程序全局唯一后台接口调用凭据
func FlushAppletAccessToken() {
	if err := flushAccessToken(); err != nil {
		return
	}

	ticker := time.NewTicker(time.Second * 7000)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := flushAccessToken(); err != nil {
				log.Log.Errorf("job_trace: flush access token fail, err:%s", err)
			}
		}
	}
}

func flushAccessToken() error {
	rds := dao.NewRedisDao()
	b, err := rds.EXISTS(rdskey.WECHAT_ACCESS_TOKEN)
	if err != nil {
		log.Log.Errorf("job_trace: rds exists err:%s", err)
		return err
	}

	if !b {
		wx := third.NewWechat()
		info := wx.GetAppletAccessToken()
		log.Log.Infof("applet info:%+v", info)
		if info.AccessToken != "" && info.ExpiresIn > 0 {
			if err := rds.SETEX(rdskey.WECHAT_ACCESS_TOKEN, info.ExpiresIn, info.AccessToken); err != nil {
				return err
			}
		}
	}

	return nil
}
