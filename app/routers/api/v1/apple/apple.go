package apple

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sports_service/server/app/config"
	"sports_service/server/global/app/log"
)

func AppleLink(c *gin.Context) {
	//c.Header("Content-Type", "application/octet-stream")
	//// 强制浏览器下载
	//c.Header("Content-Disposition", "attachment; filename=apple-app-site-association")
	//// 浏览器下载或预览
	//c.Header("Content-Disposition", "inline;filename=apple-app-site-association")
	//c.Header("Content-Transfer-Encoding", "binary")
	//c.Header("Cache-Control", "no-cache")
	//
	//c.File(config.ApplicationConfig.AppleLinkPath)
	file, err := os.Open(config.Global.AppleLinkPath)
	if err != nil {
		log.Log.Errorf("apple_trace: open file fail, err:%s", err)
		c.String(http.StatusInternalServerError, "")
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Log.Errorf("apple_trace: file stat fail, err:%s", err)
		c.String(http.StatusInternalServerError, "")
		return
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	if _, err := file.Read(buffer); err != nil {
		log.Log.Errorf("apple_trace: read buffer fail, err:%s", err)
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(200, string(buffer))
}
