package im

import (
	"errors"
	"fmt"
	"github.com/tencentyun/tls-sig-api-v2-golang/tencentyun"
	"sports_service/server/util"
)

type imRealize struct {
	ImAppId     int
	ImAppKey    string
	Identifier  string
}

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
		GroupId   string `json:"GroupId"`            // 群id
	}

	Response struct {
		BaseResponse
	}
)

const (
	EXPIRE_TM_DAY = 86400
)

func NewImRealize(appId int, appKey, identifier string) *imRealize {
	return &imRealize{
		ImAppId: appId,
		ImAppKey: appKey,
		Identifier: identifier,
	}
}

func (im *imRealize) AddUser(userId, name, avatar string) (string, error) {
	sig, err := im.GenSig(im.Identifier, EXPIRE_TM_DAY)
	if err != nil {
		return "", err
	}

	userSig, err := im.GenSig(userId, EXPIRE_TM_DAY * 90)
	if err != nil {
		return "", err
	}

	url := GenRequestUrl(im.ImAppId, im.Identifier, sig, addUserUrl)
	param := map[string]interface{}{
		"Identifier": userId,
		"Nick": name,
		"FaceUrl": avatar,
	}

	resp, err := HttpPostBody(url, param)
	if err != nil {
		return "", fmt.Errorf("resp.Body:%s, Failed to get response data! error: %s", string(resp), err.Error())
	}

	var response Response
	if err := util.JsonFast.Unmarshal(resp, &response); err != nil {
		return "", fmt.Errorf("Failed to unmarshal! error: %s", err.Error())
	}

	if response.ActionStatus != "OK" || response.ErrorCode != 0 {
		return "", errors.New(response.ErrorInfo)
	}

	return userSig, nil
}

// 创建群组
// Private	支持，同新版本中的 Work（好友工作群）
// Public	支持
// ChatRoom	支持，同新版本中的 Meeting（临时会议群）
// AVChatRoom	支持
func (im *imRealize) CreateGroup(groupType, owner, name, introduction, notification, faceUrl string) (string, error) {
	if groupType == "" || name == "" {
		return "", errors.New("invalid param")
	}

	sig, err := im.GenSig(im.Identifier, EXPIRE_TM_DAY)
	if err != nil {
		return "", err
	}

	url := GenRequestUrl(im.ImAppId, im.Identifier, sig, createGroupUrl)
	param := map[string]interface{}{
		"Owner_Account": owner,
		"Type": groupType,
		"FaceUrl": faceUrl,
		"Name": name,
		"Introduction": introduction,
		"Notification": notification,
	}

	resp, err := HttpPostBody(url, param)
	if err != nil {
		return "", fmt.Errorf("resp.Body:%s, Failed to get response data! error: %s", string(resp), err.Error())
	}

	var response Response
	if err := util.JsonFast.Unmarshal(resp, &response); err != nil {
		return "", fmt.Errorf("Failed to unmarshal! error: %s", err.Error())
	}

	if response.ActionStatus != "OK" || response.ErrorCode != 0 {
		return "", errors.New(response.ErrorInfo)
	}

	return response.GroupId, nil
}

// 生成签名
func (im *imRealize) GenSig(userId string, expireTm int) (string, error) {
	sig, err := tencentyun.GenUserSig(im.ImAppId, im.ImAppKey, userId, expireTm)
	if err != nil {
		return "", err
	}

	return sig, nil
}
