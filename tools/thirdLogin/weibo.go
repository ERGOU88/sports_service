package thirdLogin

import (
	"github.com/parnurzeal/gorequest"
	"log"
	"net/url"
	"fmt"
)

type Weibo struct {}

// 微博授权用户信息
type WeiboInfo struct {
	Request   string `json:"request"`
	ErrorCode string `json:"error_code"`
	Error     string `json:"error"`

	Id              int64  `json:"id"`                                    // 微博授权的uid
	UserId          int64  `json:"user_id,string"`                        // 用户id
	AccessToken     string `json:"access_token"`                          // 微博token
	ExpireIn        int64  `json:"expire_in"`                             // token过期时间
	IdStr           string `json:"id_str"`                                // id字符串
	ScreenName      string `json:"screen_name"`                           // 用户昵称
	Name            string `json:"name"`                                  // 友好显示名称
	Province        string `json:"province"`                              // 用户所在省级ID
	City            string `json:"city"`                                  // 用户所在城市ID
	Location        string `json:"location"`                              // 用户所在地
	Description     string `json:"description"`                           // 用户个人描述
	Url             string `json:"url"`                                   // 用户博客地址
	ProfileImageUrl string `json:"profile_image_url"`                     // 用户头像地址（中图），50×50像素
	ProfileUrl      string `json:"profile_url"`                           // 用户的微博统一URL地址
	Domain          string `json:"domain"`                                // 用户的个性化域名
	Weihao          string `json:"weihao"`                                // 用户的微号
	Gender          string `json:"gender"`                                // 性别，m：男、f：女、n：未知
	FollowersCount  int    `json:"followers_count"`                       // 粉丝数
	FriendsCount    int    `json:"friends_count"`                         // 关注数
	StatusesCount   int    `json:"statuses_count"`                        // 微博数
	FavouritesCount int    `json:"favourites_count"`                      // 收藏数
	CreatedAt       string `json:"created_at"`                            // 用户创建（注册）时间
	AllowAllActMsg  bool   `json:"allow_all_act_msg"`                     // 是否允许所有人给我发私信，true：是，false：否
	GeoEnabled      bool   `json:"geo_enabled"`                           // 是否允许标识用户的地理位置，true：是，false：否
	Verified        bool   `json:"verified"`                              // 是否是微博认证用户，即加V用户，true：是，false：否
	VerifiedType    int    `json:"verified_type"`                         // 暂未支持
	Remark          string `json:"remark"`                                // 用户备注信息，只有在查询用户关系时才返回此字段
	Status          struct {
		Geo struct {
			Longitude    string `json:"longitude"`                        // 经度坐标
			Latitude     string `json:"latitude"`                         // 维度坐标
			City         string `json:"city"`                             // 所在城市的城市代码
			Province     string `json:"province"`                         // 所在省份的省份代码
			CityName     string `json:"city_name"`                        // 所在城市的城市名称
			ProvinceName string `json:"province_name"`                    // 所在省份的省份名称
			Address      string `json:"address"`                          // 所在的实际地址，可以为空
			Pinyin       string `json:"pinyin"`                           // 地址的汉语拼音，不是所有情况都会返回该字段
			More         string `json:"more"`                             // 更多信息，不是所有情况都会返回该字段
		} `json:"geo"`
	} `json:"status"`
	AllowAllComment  bool   `json:"allow_all_comment"`                    // 是否允许所有人对我的微博进行评论，true：是，false：否
	AvatarLarge      string `json:"avatar_large"`                         // 用户头像地址（大图），180×180像素
	AvatarHd         string `json:"avatar_hd"`                            // 用户头像地址（高清），高清头像原图
	VerifiedReason   string `json:"verified_reason"`                      // 认证原因
	FollowMe         bool   `json:"follow_me"`                            // 该用户是否关注当前登录用户，true：是，false：否
	OnlineStatus     int    `json:"online_status"`                        // 用户的在线状态，0：不在线、1：在线
	BiFollowersCount int    `json:"bi_followers_count"`                   // 用户的互粉数
	Lang             string `json:"lang"`                                 // 用户当前的语言版本，zh-cn：简体中文，zh-tw：繁体中文，en：英语
}

// 微博实栗
func NewWeibo() *Weibo {
	return &Weibo{}
}

// 获取微博用户信息 (uid 即 openid)
func (wb *Weibo) GetWeiboUserInfo(uid int64, accessToken string) *WeiboInfo {
	// 通过拉取用户信息去验证access_token和uid的有效性
	v := url.Values{}
	v.Set("access_token", accessToken)
	v.Set("uid", fmt.Sprint(uid))
	weiboInfo := &WeiboInfo{}
	resp, body, errs := gorequest.New().Get(WEIBO_USER_INFO_URL + v.Encode()).EndStruct(weiboInfo)
	if errs != nil {
		log.Printf("weibo_trace: get weibo info err %+v", errs)
		return nil
	}

	log.Println("\nweiboInfo: ", weiboInfo)
	log.Println("\nresp: ", resp)
	log.Println("\nbody: ", string(body))
	log.Println("\nerrs: ", errs)

	if resp.StatusCode != 200 || weiboInfo.ErrorCode != "" {
		log.Printf("weibo_trace: request failed, errCode:%d, error:%s", weiboInfo.ErrorCode, weiboInfo.Error)
		return nil
	}

	return weiboInfo
}
