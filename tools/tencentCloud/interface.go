package tencentCloud

import (
	"sports_service/util"
)

type TencentCloudParam interface {
	// 用于提供访问的 method
	APIName() string

	// 返回参数列表
	Params() map[string]string
}

func marshal(obj interface{}) string {
	var bytes, err = util.JsonFast.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(bytes)
}
