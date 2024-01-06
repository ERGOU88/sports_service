package util

import (
	"bufio"
	"io"
	"mime/multipart"
)

// ReadFileByLine 读取文件内容并按行分割
func ReadFileByLine(file *multipart.FileHeader) ([]string, error) {
	// 打开文件
	fileContent, err := file.Open()
	if err != nil {
		return nil, err
	}
	var lines []string
	reader := bufio.NewReader(fileContent)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}
