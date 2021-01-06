package cuser

import (
  "sports_service/server/global/app/errdef"
  "sports_service/server/global/consts"
  "strings"
  "fmt"
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
func (svc *UserModule) VersionUp(versionCode string) (int, bool, interface{}) {
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

  var update interface{}
  if plt == int32(consts.ANDROID_PLATFORM) {
    update = &AndroidUpgrade{
      AndroidVersion: latestPkg.Version,
      AndroidPromptHeader: "发现新版本",
      AndroidPromptBody: fmt.Sprintf("赛赛V%s版本（%s）", latestPkg.Version, latestPkg.Size),
      AndroidPromptButton: "立即更新",
      AndroidUpgradeURL: latestPkg.UpgradeUrl,
    }
  }

  if plt == int32(consts.IOS_PLATFORM) {
    update = &IOSUpgrade{
      IOSVersion: latestPkg.Version,
      IOSPromptHeader: "发现新版本",
      IOSPromptBody: fmt.Sprintf("赛赛V%s版本（%s）", latestPkg.Version, latestPkg.Size),
      IOSPromptButton: "立即更新",
      IOSUpgradeURL: latestPkg.UpgradeUrl,
    }
  }

  return errdef.SUCCESS, isForce, update
}
