package util

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func HttpDo(uri, method, serviceName string, query []byte, timeout int64, header http.Header) ([]byte, int, error) {
	if method == http.MethodGet {
		if strings.Index(uri, "?") > 0 {
			uri = fmt.Sprintf("%s&%s", uri, string(query))
		} else {
			uri = fmt.Sprintf("%s?%s", uri, string(query))
		}
	}

	client := http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	req, err := http.NewRequest(method, uri, bytes.NewReader(query))
	if err != nil {
		return nil, 0, err
	}

	req.Header = header
	req.Header.Set("User-Agent", "youzu-go-notify")
	req.Header.Set("Service-Name", serviceName)
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, errors.New("请求返回的状态码非200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, errors.New("读取body内容时返回error,err:" + err.Error())
	}
	defer resp.Body.Close()

	return body, 200, nil

}
