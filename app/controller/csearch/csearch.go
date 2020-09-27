package csearch

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
	"sports_service/server/models/mvideo"
	"strings"
	"time"
)

type SearchModule struct {
	context     *gin.Context
	engine      *xorm.Session
	collect     *mcollect.CollectModel
	user        *muser.UserModel
	video       *mvideo.VideoModel
	attention   *mattention.AttentionModel
	like        *mlike.LikeModel
}

func New(c *gin.Context) SearchModule {
	socket := dao.Engine.Context(c)
	defer socket.Close()
	return SearchModule{
		context: c,
		collect: mcollect.NewCollectModel(socket),
		user: muser.NewUserModel(socket),
		video: mvideo.NewVideoModel(socket),
		attention: mattention.NewAttentionModel(socket),
		like: mlike.NewLikeModel(socket),
		engine: socket,
	}
}

// 综合搜索（视频+用户 默认各展示3条 视频默认播放量排序）
func (svc *SearchModule) ColligateSearch(userId, name string) ([]*mvideo.VideoDetailInfo, []*muser.UserSearchResults) {
	if name == "" {
		log.Log.Errorf("search_trace: search name can't empty, name:%s", name)
		return []*mvideo.VideoDetailInfo{}, []*muser.UserSearchResults{}
	}

	// 搜索到的视频
	videos := svc.VideoSearch(userId, name, consts.VIDEO_CONDITION_PLAY, string(consts.UNLIMITED_DURATION),
		string(consts.UNLIMITED_TIME), consts.DEFAULT_SEARCH_PAGE, consts.DEFAULT_SEARCH_SIZE)
	// 搜索到的用户
	users := svc.UserSearch(userId, name, consts.DEFAULT_SEARCH_PAGE, consts.DEFAULT_SEARCH_SIZE)

	return videos, users
}

// 视频搜索
func (svc *SearchModule) VideoSearch(userId, name, sort, duration, publishTime string, page, size int) []*mvideo.VideoDetailInfo {
	if name == "" {
		log.Log.Errorf("search_trace: search name can't empty, name:%s", name)
		return []*mvideo.VideoDetailInfo{}
	}

	sortField := svc.GetSortField(sort)
	min, max := svc.GetDurationCondition(duration)
	pubTime := svc.GetPublishTimeCondition(publishTime)
	offset := (page - 1) * size

	list := svc.video.SearchVideos(name, sortField, min, max, pubTime, offset, size)
	if len(list) == 0 {
		return []*mvideo.VideoDetailInfo{}
	}

	for _, video := range list {
		// 查询用户信息
		userInfo := svc.user.FindUserByUserid(video.UserId)
		if userInfo == nil {
			log.Log.Errorf("video_trace: user not found, uid:%s", video.UserId)
			continue
		}

		video.Avatar = userInfo.Avatar
		video.Nickname = userInfo.NickName
		// 用户未登录
		if userId == "" {
			log.Log.Error("search_trace: no login")
			continue
		}

		// 是否关注
		if attentionInfo := svc.attention.GetAttentionInfo(userId, video.UserId); attentionInfo != nil {
			video.IsAttention = attentionInfo.Status
		}

		// 获取点赞的信息
		if likeInfo := svc.like.GetLikeInfo(userId, video.VideoId, consts.TYPE_VIDEO); likeInfo != nil {
			video.IsLike = likeInfo.Status
		}

		// 获取收藏的信息
		if collectInfo := svc.collect.GetCollectInfo(userId, video.VideoId, consts.TYPE_VIDEO); collectInfo != nil {
			video.IsCollect = collectInfo.Status
		}

	}

	return list
}

// 搜索用户
func (svc *SearchModule) UserSearch(userId, name string, page, size int) []*muser.UserSearchResults {
	if name == "" {
		log.Log.Errorf("search_trace: search user name can't empty, name:%s", name)
		return []*muser.UserSearchResults{}
	}

	offset := (page - 1) * size
	list := svc.user.SearchUser(name, offset, size)
	if len(list) == 0 {
		return []*muser.UserSearchResults{}
	}

	for _, user := range list {
		// 获取搜索出的用户 粉丝总数
		user.FansNum = svc.attention.GetTotalFans(user.UserId)
		// 获取搜索出的用户 作品总数
		user.WorksNum = svc.video.GetTotalPublish(user.UserId)

		if userId == "" {
			log.Log.Error("search_trace: user no login")
			continue
		}

		// 是否关注
		if attentionInfo := svc.attention.GetAttentionInfo(userId, user.UserId); attentionInfo != nil {
			user.IsAttention = int32(attentionInfo.Status)
		}

	}

	return list
}

// 标签搜索
func (svc *SearchModule) LabelSearch(userId string, labelId string, page, size int) []*mvideo.VideoDetailInfo {
	offset := (page - 1) * size
	videoIds := svc.video.GetVideoIdsByLabelId(labelId, offset, size)
	if len(videoIds) == 0 {
		log.Log.Errorf("search_trace: not found videos by label id, labelId:%s", labelId)
		return []*mvideo.VideoDetailInfo{}
	}

	vids := strings.Join(videoIds, ",")
	videos := svc.video.FindVideoListByIds(vids)
	if len(videos) == 0 {
		log.Log.Errorf("search_trace: not found videos, vids:%s", vids)
		return []*mvideo.VideoDetailInfo{}
	}

	// 重新组装数据
	list := make([]*mvideo.VideoDetailInfo, len(videos))
	for index, video := range videos {
		resp := new(mvideo.VideoDetailInfo)
		resp.VideoId = video.VideoId
		resp.Title = video.Title
		resp.Describe = video.Describe
		resp.Cover = video.Cover
		resp.VideoAddr = video.VideoAddr
		resp.IsRecommend = video.IsRecommend
		resp.IsTop = video.IsTop
		resp.VideoDuration = video.VideoDuration
		resp.VideoWidth = video.VideoWidth
		resp.VideoHeight = video.VideoHeight
		resp.CreateAt = video.CreateAt
		resp.UserId = video.UserId
		// 获取用户信息
		if user := svc.user.FindUserByUserid(video.UserId); user != nil {
			resp.Avatar = user.Avatar
			resp.Nickname = user.NickName
		}

		// 用户未登录
		if userId == "" {
			log.Log.Error("search_trace: no login")
			continue
		}

		// 是否关注
		if attentionInfo := svc.attention.GetAttentionInfo(userId, video.UserId); attentionInfo != nil {
			resp.IsAttention = attentionInfo.Status
		}

		// 获取点赞的信息
		if likeInfo := svc.like.GetLikeInfo(userId, video.VideoId, consts.TYPE_VIDEO); likeInfo != nil {
			resp.IsLike = likeInfo.Status
		}

		// 获取收藏的信息
		if collectInfo := svc.collect.GetCollectInfo(userId, video.VideoId, consts.TYPE_VIDEO); collectInfo != nil {
			resp.IsCollect = collectInfo.Status
		}

		list[index] = resp
	}

	return list
}

// 获取热门搜索（后台配置的结果集）
func (svc *SearchModule) GetHotSearch() []string {
	hot := svc.video.GetHotSearch()
	return strings.Split(hot.HotSearchContent, ",")
}

// 搜索关注的用户
func (svc *SearchModule) SearchAttentionUser(userId, name string, page, size int) []*mattention.SearchContactRes {
	if name == "" || userId == "" {
		log.Log.Errorf("search_trace: search attention name can't empty, name:%s", name)
		return []*mattention.SearchContactRes{}
	}

	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("search_trace: user not found, uid:%s", userId)
		return []*mattention.SearchContactRes{}
	}

	offset := (page - 1) * size
	return svc.attention.SearchAttentionUser(userId, name, offset, size)
}

// 搜索粉丝列表
func (svc *SearchModule) SearchFans(userId, name string, page, size int) []*mattention.SearchContactRes {
	if name == "" || userId == "" {
		log.Log.Errorf("search_trace: search fans name can't empty, name:%s", name)
		return []*mattention.SearchContactRes{}
	}

	if user := svc.user.FindUserByUserid(userId); user == nil {
		log.Log.Errorf("search_trace: user not found, uid:%s", userId)
		return []*mattention.SearchContactRes{}
	}

	offset := (page - 1) * size
	list := svc.attention.SearchFans(userId, name, offset, size)
	for _, info := range list {
		// 是否回关
		if attentionInfo := svc.attention.GetAttentionInfo(userId, info.UserId); attentionInfo != nil {
			info.IsAttention = int32(attentionInfo.Status)
		}
	}

	return list
}

// 获取时长条件
func (svc *SearchModule) GetDurationCondition(durationType string) (minDuration, maxDuration int64) {
	switch durationType {
	case string(consts.UNLIMITED_DURATION):
		return
	case string(consts.ONE_TO_FIVE_MINUTES):
		minDuration = 1 * 60
		maxDuration = 5 * 60
		return
	case string(consts.FIVE_TO_TEN_MINUTES):
		minDuration = 5 * 60
		maxDuration = 10 * 60
		return
	case string(consts.TEN_TO_HALF_HOUR):
		minDuration = 10 * 60
		maxDuration = 30 * 60
		return
	case string(consts.MORE_THAN_HALF_HOUR):
		minDuration = 30 * 60
		maxDuration = 9e10
	default:
		log.Log.Errorf("search_trace: unsupported duration condition, durationType: %s", durationType)
	}

	return
}

// 获取发布时间条件
func (svc *SearchModule) GetPublishTimeCondition(publishTime string) int64 {
	switch publishTime {
	case string(consts.UNLIMITED_TIME):
		return 0
	case string(consts.A_DAY):
		return time.Now().Unix() - 60 * 60 *24
	case string(consts.A_WEEK):
		return time.Now().Unix() - 60 * 60 * 24 * 7
	case string(consts.HALF_A_YEAR):
		return time.Now().Unix() - 365 / 2 * 24 * 60 * 60
	}

	return 0
}

// 获取排序字段
// 0 播放数 1 弹幕数 2 评论数 3 点赞数
func (svc *SearchModule) GetSortField(condition string) string {
	switch condition {
	// 播放数
	case consts.VIDEO_CONDITION_PLAY:
		return consts.CONDITION_FIELD_PLAY
	// 弹幕数
	case consts.VIDEO_CONDITION_BARRAGE:
		return consts.CONDITION_FIELD_BARRAGE
	// 点赞数
	case consts.VIDEO_CONDITION_LIKE:
		return consts.CONDITION_FIELD_LIKE
	default:
		log.Log.Errorf("search_trace: unsupported condition, condition: %s", condition)
	}

	return consts.CONDITION_FIELD_PLAY
}