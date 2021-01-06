package cuser

import (
  "sports_service/server/global/app/errdef"
  "sports_service/server/global/consts"
  "sports_service/server/models/mconfigure"
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
  IOSVersion          string `json:"ios_version"`
  IOSPromptHeader     string `json:"ios_prompt_header"`
  IOSPromptBody       string `json:"ios_prompt_body"`
  IOSPromptButton     string `json:"ios_prompt_button"`
  IOSUpgradeURL       string `json:"ios_upgrade_url" example:"新包下载地址"`
}

// 通过版本code获取包信息
func (svc *UserModule) VersionUp(versionCode string) (int, bool, *mconfigure.UpgradeInfo) {
  var plt int32
  appId := svc.context.GetHeader("AppId")
  if strings.Compare(appId, string(consts.IOS_APP_ID)) == 0 {
    plt = int32(consts.IOS_PLATFORM)
  }

  if strings.Compare(appId, string(consts.AND_APP_ID)) == 0 {
    plt = int32(consts.ANDROID_PLATFORM)
  }

  // 通过版本code及平台 获取当前包信息
  info := svc.configure.GetPackageDetailByVersion(plt, versionCode)
  if info == nil {
    return errdef.USER_PACKAGE_NOT_EXISTS, false, nil
  }

  // 获取平台最新包信息
  latestPkg := svc.configure.GetLatestPackageInfo(plt)
  if latestPkg == nil {
    return errdef.USER_LATEST_PACKAGE_FAIL, false, nil
  }

  var isForce bool
  // isForce == 1 需要强更
  if info.IsForce == 1 {
    isForce = true
  }

  upgrade := &mconfigure.UpgradeInfo{
    Version:  latestPkg.Version,
    VersionCode: latestPkg.VersionCode,
    Size: latestPkg.Size,
    Describe: latestPkg.Describe,
    UpgradeURL: latestPkg.UpgradeUrl,
  }

  return errdef.SUCCESS, isForce, upgrade
}
