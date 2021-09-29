package mcommunity

// 添加板块
type AddSection struct {
	SectionName string `json:"section_name"`
	Sortorder   int    `json:"sortorder"`
}

// 添加社区板块
func (m *CommunityModel) AddCommunitySection() {

}
