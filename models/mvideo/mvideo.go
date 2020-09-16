package mvideo

import (
	"github.com/go-xorm/xorm"
	"sports_service/server/models"
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

