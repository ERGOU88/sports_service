package cuser

import (
	"sports_service/global/app/errdef"
	"sports_service/global/consts"
	"sports_service/models/mconfigure"
	"strconv"
	"strings"
)

type AndroidUpgrade struct {
	AndroidVersion      string `json:"android_version" example:"版本"`
	AndroidPromptHeader string `json:"android_prompt_header" example:"发现新版本"`
	AndroidPromptBody   string `json:"android_prompt_body" example:"您的APP当前不是最新版本"`
	AndroidPromptButton string `json:"android_prompt_button" example:"请前往更新"`
	AndroidUpgradeURL   string `json:"android_upgrade_url" example:"新包下载地址"`
}

type IOSUpgrade struct {
	IOSVersion      string `json:"ios_version"`
	IOSPromptHeader string `json:"ios_prompt_header"`
	IOSPromptBody   string `json:"ios_prompt_body"`
	IOSPromptButton string `json:"ios_prompt_button"`
	IOSUpgradeURL   string `json:"ios_upgrade_url" example:"新包下载地址"`
}

// 通过版本code获取包信息
// iOS version 例如：v1.0.1,v1.0.2
// Android version 例如：1001 1002
func (svc *UserModule) VersionUp(version string) (int, *mconfigure.UpgradeInfo) {
	var (
		plt  int32
		code int
	)

	appId := svc.context.GetHeader("AppId")
	if appId == "" {
		appId = svc.context.Query("Appid")
	}

	if strings.Compare(appId, string(consts.IOS_APP_ID)) == 0 {
		plt = int32(consts.IOS_PLATFORM)
		pkgInfo := svc.configure.GetPackageByVersion(version, plt)
		if pkgInfo != nil {
			code = pkgInfo.VersionCode
		}
	}

	if strings.Compare(appId, string(consts.AND_APP_ID)) == 0 {
		plt = int32(consts.ANDROID_PLATFORM)
		code, _ = strconv.Atoi(version)
	}

	// 通过版本code及平台 获取当前包信息
	//info := svc.configure.GetPackageDetailByVersion(plt, versionCode)
	//if info == nil {
	//  return errdef.USER_PACKAGE_NOT_EXISTS, false, nil
	//}

	// 获取平台最新包信息
	latestPkg := svc.configure.GetLatestPackageInfo(plt)
	if latestPkg == nil {
		return errdef.USER_LATEST_PACKAGE_FAIL, nil
	}

	upgrade := &mconfigure.UpgradeInfo{
		Version:     latestPkg.Version,
		VersionCode: latestPkg.VersionCode,
		Size:        latestPkg.Size,
		Describe:    latestPkg.Describe,
		UpgradeURL:  latestPkg.UpgradeUrl,
	}

	if latestPkg.VersionCode > code {
		upgrade.HasNewPkg = true
	}

	// isForce == 1 需要强更
	if latestPkg.IsForce == 1 {
		upgrade.IsForce = true
	}

	return errdef.SUCCESS, upgrade
}
