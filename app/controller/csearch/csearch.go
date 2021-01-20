package csearch

import (
  "github.com/gin-gonic/gin"
  "github.com/go-xorm/xorm"
  "sports_service/server/dao"
  "sports_service/server/global/app/errdef"
  "sports_service/server/global/app/log"
  "sports_service/server/global/consts"
  "sports_service/server/models"
  "sports_service/server/models/mattention"
  "sports_service/server/models/mcollect"
  "sports_service/server/models/mlike"
  "sports_service/server/models/muser"
  "sports_service/server/models/mvideo"
  "sports_service/server/util"
  "time"
  "fmt"
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
  socket := dao.Engine.NewSession()
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

// 综合搜索（视频+用户 默认视频取10条 用户取20条 视频默认播放量排序） 如果视频和用户都未搜索到 则推荐两个视频
func (svc *SearchModule) ColligateSearch(userId, name string) ([]*mvideo.VideoDetailInfo, []*muser.UserSearchResults, []*mvideo.VideoDetailInfo) {
  if name == "" {
		log.Log.Errorf("search_trace: search name can't empty, name:%s", name)
		return []*mvideo.VideoDetailInfo{}, []*muser.UserSearchResults{}, []*mvideo.VideoDetailInfo{}
	}

  length := util.GetStrLen([]rune(name))
  if length > 20 {
    log.Log.Errorf("search_trace: invalid search name len, len:%s", length)
    return []*mvideo.VideoDetailInfo{}, []*muser.UserSearchResults{}, []*mvideo.VideoDetailInfo{}
  }

	if userId != "" {
	  if err := svc.video.RecordHistorySearch(userId, name); err != nil {
	    log.Log.Errorf("search_trace: record history search err:%s", err)
    }
	}

	// 搜索到的视频
	videos := svc.VideoSearch(userId, name, consts.VIDEO_CONDITION_PLAY, string(consts.UNLIMITED_DURATION),
		string(consts.UNLIMITED_TIME), consts.DEFAULT_SEARCH_VIDEO_PAGE, consts.DEFAULT_SEARCH_VIDEO_SIZE)
	// 搜索到的用户
	users := svc.UserSearch(userId, name, consts.DEFAULT_SEARCH_USER_PAGE, consts.DEFAULT_SEARCH_USER_SIZE)

	var recommend []*mvideo.VideoDetailInfo
	if len(videos) == 0 && len(users) == 0 {
    recommend = svc.RecommendVideo()
  }

  if len(recommend) == 0 {
    recommend = []*mvideo.VideoDetailInfo{}
  }

	return videos, users, recommend
}

// 推荐视频 默认取两条
func (svc *SearchModule) RecommendVideo() []*mvideo.VideoDetailInfo {
  offset := util.GenerateRandnum(0, int(svc.video.GetVideoTotalCount())-10)
  videos := svc.video.GetRecommendVideos(int32(offset), 2)
  if videos == nil {
    log.Log.Error("search_trace: get recommend video fail")
    return []*mvideo.VideoDetailInfo{}
  }

  for _, val := range videos {
    val.Title = util.TrimHtml(val.Title)
    val.Describe = util.TrimHtml(val.Describe)
     user := svc.user.FindUserByUserid(val.UserId)
     if user != nil {
       val.Nickname = user.NickName
       val.Avatar = user.Avatar
     }
  }

  return videos
}


// 视频搜索 todo: 限制搜索的字符数
func (svc *SearchModule) VideoSearch(userId, name, sort, duration, publishTime string, page, size int) []*mvideo.VideoDetailInfo {
	if name == "" {
		log.Log.Errorf("search_trace: search name can't empty, name:%s", name)
		return []*mvideo.VideoDetailInfo{}
	}

	length := util.GetStrLen([]rune(name))
	if length > 20 {
	  log.Log.Errorf("search_trace: invalid search name len, len:%s", length)
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

    video.Title = util.TrimHtml(video.Title)
    video.Describe = util.TrimHtml(video.Describe)

		video.Avatar = userInfo.Avatar
		video.Nickname = userInfo.NickName
		video.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)

		if len(video.PlayInfo) == 0 {
		  video.PlayInfo = []*mvideo.PlayInfo{}
    }

    if len(video.Labels) == 0 {
      video.Labels = []*models.VideoLabels{}
    }

    // 获取视频统计数据
    info := svc.video.GetVideoStatistic(fmt.Sprint(video.VideoId))
    if info != nil {
      video.BarrageNum = info.BarrageNum
    }

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
		if likeInfo := svc.like.GetLikeInfo(userId, video.VideoId, consts.TYPE_VIDEOS); likeInfo != nil {
			video.IsLike = likeInfo.Status
		}

		// 获取收藏的信息
		if collectInfo := svc.collect.GetCollectInfo(userId, video.VideoId, consts.TYPE_VIDEO); collectInfo != nil {
			video.IsCollect = collectInfo.Status
		}

	}

	return list
}

// 搜索用户 todo：限制搜索的字符数
func (svc *SearchModule) UserSearch(userId, name string, page, size int) []*muser.UserSearchResults {
	if name == "" {
		log.Log.Errorf("search_trace: search user name can't empty, name:%s", name)
		return []*muser.UserSearchResults{}
	}

  length := util.GetStrLen([]rune(name))
  if length > 20 {
    log.Log.Errorf("search_trace: invalid search name len, len:%s", length)
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
	//videoIds := svc.video.GetVideoIdsByLabelId(labelId, offset, size)
	//if len(videoIds) == 0 {
	//	log.Log.Errorf("search_trace: not found videos by label id, labelId:%s", labelId)
	//	return []*mvideo.VideoDetailInfo{}
	//}
  //
	//vids := strings.Join(videoIds, ",")
  //log.Log.Debugf("videoIds:%v, vids:%s", videoIds, vids)
	//videos := svc.video.FindVideoListByIds(vids)
	//if len(videos) == 0 {
	//	log.Log.Errorf("search_trace: not found videos, vids:%s", vids)
	//	return []*mvideo.VideoDetailInfo{}
	//}

	// 通过标签id 获取同标签视频列表
	videos := svc.video.SearchVideosByLabelId(labelId, offset, size)
  if len(videos) == 0 {
  	log.Log.Errorf("search_trace: not found videos, labelId:%s", labelId)
  	return []*mvideo.VideoDetailInfo{}
  }

	// 重新组装数据
	list := make([]*mvideo.VideoDetailInfo, len(videos))
	for index, video := range videos {
		resp := new(mvideo.VideoDetailInfo)
		resp.VideoId = video.VideoId
		resp.Title = util.TrimHtml(video.Title)
		resp.Describe = util.TrimHtml(video.Describe)
		resp.Cover = video.Cover
		resp.VideoAddr = svc.video.AntiStealingLink(video.VideoAddr)
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

    // 获取视频统计数据
    info := svc.video.GetVideoStatistic(fmt.Sprint(video.VideoId))
    if info != nil {
       resp.BarrageNum = info.BarrageNum
       resp.BrowseNum = info.BrowseNum
    }

		// 用户未登录
		if userId != "" {
      // 是否关注
      if attentionInfo := svc.attention.GetAttentionInfo(userId, video.UserId); attentionInfo != nil {
        resp.IsAttention = attentionInfo.Status
      }

      // 获取点赞的信息
      if likeInfo := svc.like.GetLikeInfo(userId, video.VideoId, consts.TYPE_VIDEOS); likeInfo != nil {
        resp.IsLike = likeInfo.Status
      }

      // 获取收藏的信息
      if collectInfo := svc.collect.GetCollectInfo(userId, video.VideoId, consts.TYPE_VIDEO); collectInfo != nil {
        resp.IsCollect = collectInfo.Status
      }
		}

		resp.Labels = []*models.VideoLabels{}
		resp.PlayInfo = []*mvideo.PlayInfo{}

		list[index] = resp
	}

	return list
}

// 获取热门搜索（后台配置的结果集）
func (svc *SearchModule) GetHotSearch() []*models.HotSearch {
	return  svc.video.GetHotSearch()
}

// 历史搜索记录
func (svc *SearchModule) GetHistorySearch(userId string) []string {
  history := svc.video.GetHistorySearch(userId)
  if history == nil {
    history = []string{}
  }

  return history
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
	list := svc.attention.SearchAttentionUser(userId, name, offset, size)
	if len(list) == 0 {
	  return []*mattention.SearchContactRes{}
  }

  return list
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
  if len(list) == 0 {
    return []*mattention.SearchContactRes{}
  }

	for _, info := range list {
		// 是否回关
		if attentionInfo := svc.attention.GetAttentionInfo(userId, info.UserId); attentionInfo != nil {
			info.IsAttention = int32(attentionInfo.Status)
		}
	}

	return list
}

// 清空搜索历史
func (svc *SearchModule) CleanSearchHistory(userId string) int {
  if err := svc.video.CleanHistorySearch(userId); err != nil {
    log.Log.Errorf("search_trace: clean history search err:%s", err)
    return errdef.SEARCH_CLEAN_HISTORY_FAIL
  }

  return errdef.SUCCESS
}

// 获取时长条件（数据库存储的视频时长是毫秒）
func (svc *SearchModule) GetDurationCondition(durationType string) (minDuration, maxDuration int64) {
	switch durationType {
	case string(consts.UNLIMITED_DURATION):
		return
	case string(consts.ONE_TO_FIVE_MINUTES):
		minDuration = 1 * 60 * 1000
		maxDuration = 5 * 60 * 1000
		return
	case string(consts.FIVE_TO_TEN_MINUTES):
		minDuration = 5 * 60 * 1000
		maxDuration = 10 * 60 * 1000
		return
	case string(consts.TEN_TO_HALF_HOUR):
		minDuration = 10 * 60 * 1000
		maxDuration = 30 * 60 * 1000
		return
	case string(consts.MORE_THAN_HALF_HOUR):
		minDuration = 30 * 60 * 1000
		maxDuration = 9e10 * 1000
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

