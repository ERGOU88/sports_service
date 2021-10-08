package minformation

// 更新资讯信息
func (m *InformationModel) UpdateInfo(condition, cols string) (int64, error) {
	return m.Engine.Where(condition).Cols(cols).Update(m.Information)
}
