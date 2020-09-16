package mvideo

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/global/app/log"
	"sports_service/server/models"
	"fmt"
)

type VideoModel struct {
	Videos     *models.Videos
	Engine     *xorm.Session

}

// 实栗
func NewVideoModel(engine *xorm.Session) *VideoModel {
	return &VideoModel{
		Videos: new(models.Videos),
		Engine: engine,
	}
}

// 分页获取 用户发布的视频列表
func (m *VideoModel) GetUserPublishVideos(offset, page int) {
	return
}

// 通过id查询视频
func (m *VideoModel) FindVideoById(videoId string) *models.Videos {
	ok, err := m.Engine.Where("video_id=?").Get(m.Videos)
	if !ok || err != nil {
		return nil
	}

	return m.Videos
}

// 通过视频id查询视频列表
func (m *VideoModel) FindVideoListByIds(videoIds string, offset, size int) []*models.Videos {
	var list []*models.Videos
	sql := fmt.Sprintf("SELECT * FROM videos WHERE video_id in(%s) AND status=0 ORDER BY is_top DESC, is_recommend DESC, sortorder DESC, id LIMIT ?, ?", videoIds)
	if err := m.Engine.SQL(sql, offset, size).Find(&list); err != nil {
		log.Log.Errorf("video_trace: get video list err:%s", err)
		return nil
	}

	return list
}


