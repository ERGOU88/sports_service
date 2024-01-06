package util

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"time"
)

type ExportExcel struct {
	Title   []interface{}
	Data    [][]interface{}
	Content io.ReadSeeker
}

func NewExcelExporter(title []interface{}, data [][]interface{}) *ExportExcel {
	return &ExportExcel{
		Title: title,
		Data:  data,
	}
}
func (svc *ExportExcel) Export() *ExportExcel {
	wb := excelize.NewFile()
	sheetName := wb.GetSheetName(wb.GetActiveSheetIndex())
	// 生成流写入对象
	streamSheet, err := wb.NewStreamWriter(sheetName)
	if err != nil {
		fmt.Println(err)
	}
	// 插入表头
	streamSheet.SetRow("1", svc.Title)
	for i, data := range svc.Data {
		streamSheet.SetRow(fmt.Sprintf("A%d", i+1), data)
	}
	// 执行了 flush 才算是写进去了
	streamSheet.Flush()
	var buffer bytes.Buffer
	_ = wb.Write(&buffer)
	svc.Content = bytes.NewReader(buffer.Bytes())
	return svc
}
func (svc *ExportExcel) ToResponse(c *gin.Context, fileName string) {
	fileName = fmt.Sprintf("%s%s%s.xlsx", time.Now().Format("2006-01-02 15:04:05"), `-`, fileName)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeContent(c.Writer, c.Request, fileName, time.Now(), svc.Content)
}
