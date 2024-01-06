package tencentCloud

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sports_service/util"
	"time"
)

type BucketURI string

func (uri BucketURI) String() string {
	if string(uri) == "" {
		return ""
	}
	u, err := url.Parse(string(uri))
	if err != nil {
		return ""
	}

	now := time.Now().Unix()
	sign := util.Md5String(fmt.Sprintf("%s%s%d", CDN_SECRET, u.Path, now))
	return fmt.Sprintf("%s%s?sign=%s&t=%d", CDN_HOST, u.Path, sign, now)
}

func (uri BucketURI) MarshalJSON() ([]byte, error) {
	return json.Marshal(uri.String())
}
