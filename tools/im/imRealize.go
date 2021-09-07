package im

import (
	"errors"
	"fmt"
	"sports_service/server/util"
)

type imRealize struct {}

var Im ImInterface
var (
	// 导入用户 [单个]
	addUserUrl      = "/v4/im_open_login_svc/account_import"
	// 设置资料
	portraitSetUrl  = "/v4/profile/portrait_set"
    // 创建群组
	createGroupUrl  = "/v4/group_open_http_svc/create_group"
	// 获取直播群在线人数
	getOnlineNumUrl = "/v4/group_open_http_svc/get_online_member_num"
)

type (
	BaseResponse struct {
		ActionStatus string `json:"ActionStatus"`    // OK 表示处理成功，FAIL 表示失败
		ErrorCode int `json:"ErrorCode"`
		ErrorInfo string `json:"ErrorInfo"`          // 0表示成功，非0表示失败
	}

	ResponseAddUser struct {
		BaseResponse
	}
)

const (
	EXPIRE_TM = 86400
)

func NewImRealize() *imRealize {
	return &imRealize{}
}

func (im *imRealize) AddUser(userId, name, avatar string) (string, error) {
	sig, err := GenSig(EXPIRE_TM)
	if err != nil {
		return "", err
	}

	url := GenRequestUrl(sig, addUserUrl)
	param := map[string]interface{}{
		"Identifier": userId,
		"Nick": name,
		"FaceUrl": avatar,
	}

	resp, err := HttpPostBody(url, param)
	if err != nil {
		return "", fmt.Errorf("resp.Body:%s, Failed to get response data! error: %s", string(resp), err.Error())
	}

	var response ResponseAddUser
	if err := util.JsonFast.Unmarshal(resp, &response); err != nil {
		return "", fmt.Errorf("Failed to unmarshal! error: %s", err.Error())
	}

	if response.ActionStatus != "OK" || response.ErrorCode != 0 {
		return "", errors.New(response.ErrorInfo)
	}

	return sig, nil
}
