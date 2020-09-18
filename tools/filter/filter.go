package filter

import (
	"bytes"
	"fmt"
	filter "github.com/antlinker/go-dirtyfilter"
	"github.com/antlinker/go-dirtyfilter/store"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// FilterManage 敏感词过滤器
var FilterManage *filter.DirtyManager

type FilterModel struct {
}

func NewFilterModel() *FilterModel {
	return new(FilterModel)
}

// todo: 应抽独立服务
func init() {
	c := NewFilterModel()
	c.FlushDirtyData([]string{"./dirty.txt"})
	return
}

// FlushDirtyData 初始化过滤器
func (c *FilterModel) FlushDirtyData(filepath []string) {
	buffer := new(bytes.Buffer)
	for _, filterPath := range filepath {
		f, err := os.Open(filterPath)
		if err != nil {
			log.Printf("open file err:%s filePath:%s", err, filterPath)
			continue
		}
		text, err := ioutil.ReadAll(f)
		if err != nil {
			log.Printf("readall err:%s, %s", err, filterPath)
			continue
		}
		// log.Printf("text:%s", text)
		buffer.WriteString(string(text))
	}
	dirtyArray := strings.Split(buffer.String(), "\n")
	deal := make([]string, 0)
	for _, dirty := range dirtyArray {
		trimTxt := strings.Trim(dirty, " ")
		if trimTxt == "" {
			continue
		}
		deal = append(deal, trimTxt)

	}
	memStore, err := store.NewMemoryStore(store.MemoryConfig{
		DataSource: deal,
	})
	if err != nil {
		panic(err)
	}
	FilterManage = filter.NewDirtyManager(memStore)
}

// ValidateText 校验敏感词
func (c *FilterModel) ValidateText(text string) (isPass bool, err error) {
	var filters []string
	filters, err = FilterManage.Filter().Filter(text, '*', '@', ' ', '&', '%', '^', '!', '#', '$', '(', ')', '+', '-')
	if err != nil {
		return
	}

	if len(filters) > 0 {
		err = fmt.Errorf("该文本中包含污秽词 :%s", filters)
		return
	}
	isPass = true
	return
}

// 敏感词替换成 **
func (c *FilterModel) ReplaceDirtyText(text string) string {
	var filters []string
	filters, _ = FilterManage.Filter().Filter(text, '*', '@', ' ', '&', '%', '^', '!', '#', '$', '(', ')', '+', '-')

	if len(filters) > 0 {
		for _, filter := range filters {
			text = strings.Replace(text, filter, c.getReplaceText(), -1)
		}
	}
	return text
}

func (c *FilterModel) getReplaceText() string {
	return "**"
}
