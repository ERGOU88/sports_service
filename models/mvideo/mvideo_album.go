package mvideo

import (
	"sports_service/server/models"
)

// 创建专辑请求参数
type CreateAlbumParam struct {
	AlbumName     string    `binding:"required" json:"album_name"`    // 专辑名称
}

// 视频添加到专辑中
type AddVideoToAlbumParam struct {
	VideoId       string    `binding:"required" json:"video_id"`     // 视频id
	AlbumId       string    `binding:"required" json:"album_id"`

}

// 创建视频专辑
func (m *VideoModel) CreateVideoAlbum() (int64, error) {
	return m.Engine.InsertOne(m.Album)
}

// 通过专辑id获取专辑信息
func (m *VideoModel) GetVideoAlbumById(id string) (*models.VideoAlbum, error) {
	m.Album = new(models.VideoAlbum)
	ok, err := m.Engine.Where("id=?", id).Get(m.Album)
	if !ok || err != nil {
		return nil, err
	}

	return m.Album, nil
}

// 获取用户添加的专辑列表
func (m *VideoModel) GetVideoAlbumByUser(userId string) ([]*models.VideoAlbum, error) {
	var list []*models.VideoAlbum
	if err := m.Engine.Where("user_id=? AND status=0", userId).Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}

// 将视频添加到某个专辑
func (m *VideoModel) AddVideoToAlbum() (int64, error) {
	return m.Engine.InsertOne(m.AlbumDetail)
}

// 通过专辑id 获取专辑下的视频列表 [按添加先后排序]
func (m *VideoModel) GetVideoListByAlbumId(albumId string) ([]*models.VideoAlbumDetail, error) {
	var list []*models.VideoAlbumDetail
	if err := m.Engine.Where("album_id=? AND status=0", albumId).Asc("id").Find(&list); err != nil {
		return nil, err
	}

	return list, nil
}
