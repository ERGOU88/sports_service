package cinformation

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/app/errdef"
	"sports_service/server/global/app/log"
	"sports_service/server/global/consts"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/minformation"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
	"sports_service/server/util"
	"time"
)

type InformationModule struct {
	context     *gin.Context
	engine      *xorm.Session
	user        *muser.UserModel
	information *minformation.InformationModel
	attention   *mattention.AttentionModel
	like        *mlike.LikeModel
	collect     *mcollect.CollectModel
}

func New(c *gin.Context) InformationModule {
	socket := dao.AppEngine.NewSession()
	defer socket.Close()
	return InformationModule{
		context: c,
		user: muser.NewUserModel(socket),
		information: minformation.NewInformationModel(socket),
		attention: mattention.NewAttentionModel(socket),
		like: mlike.NewLikeModel(socket),
		collect: mcollect.NewCollectModel(socket),
		engine: socket,
	}
}

// 获取资讯列表
func (svc *InformationModule) GetInformationList(page, size int) (int, []*minformation.InformationResp) {
	offset := (page - 1) * size
	list, err := svc.information.GetInformationList(offset, size)
	if err != nil {
		return errdef.INFORMATION_LIST_FAIL, []*minformation.InformationResp{}
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*minformation.InformationResp{}
	}

	now := int(time.Now().Unix())
	resp := make([]*minformation.InformationResp, len(list))
	for index, information := range list {
		info := &minformation.InformationResp{
			Id: information.Id,
			Cover: information.Cover,
			Title: information.Title,
			CreateAt: information.CreateAt,
			//JumpUrl: information.JumpUrl,
			UserId: information.UserId,
		}

		if user := svc.user.FindUserByUserid(info.UserId); user != nil {
			info.Avatar = user.Avatar
			info.NickName = user.NickName
		}

		ok, err := svc.information.GetInformationStatistic(fmt.Sprint(info.Id))
		if !ok && err != nil {
			log.Log.Error("information_trace: get information statistic fail, id:%d, err:%s", info.Id, err)
		}

		if !ok && err == nil {
			svc.information.Statistic.NewsId = info.Id
			svc.information.Statistic.CreateAt = now
			svc.information.Statistic.UpdateAt = now
			// 初始化视频统计数据
			if _, err = svc.information.AddInformationStatistic(); err != nil {
				log.Log.Errorf("information_trace: add statistic id:%d, err:%s", info.Id, err)
				return errdef.INFORMATION_LIST_FAIL, []*minformation.InformationResp{}
			}
		}

		if ok {
			info.FabulousNum = svc.information.Statistic.FabulousNum
			info.CommentNum = svc.information.Statistic.CommentNum
			info.ShareNum = svc.information.Statistic.ShareNum
		}

		resp[index] = info
	}

	return errdef.SUCCESS, resp
}

// 获取资讯详情
func (svc *InformationModule) GetInformationDetail(id, userId string) (int, *minformation.InformationResp) {
	if id == "" {
		return errdef.INVALID_PARAMS, nil
	}

	ok, err := svc.information.GetInformationById(id)
	if !ok || err != nil {
		return errdef.INFORMATION_NOT_EXISTS, nil
	}

	node, err := util.GetHtmlNode(svc.information.Information.Content)
	if err != nil {
		log.Log.Errorf("information_trace: get body content fail, id:%s, err:%s", id, err)
		return errdef.INFORMATION_DETAIL_FAIL, nil
	}

	resp := &minformation.InformationResp{
		Id: svc.information.Information.Id,
		Cover: svc.information.Information.Cover,
		Title: svc.information.Information.Title,
		CreateAt: svc.information.Information.CreateAt,
		//JumpUrl: svc.information.Information.JumpUrl,
		UserId: svc.information.Information.UserId,
		Content: util.RenderNode(node),
	}

	if user := svc.user.FindUserByUserid(resp.UserId); user != nil {
		resp.Avatar = user.Avatar
		resp.NickName = user.NickName
	}

	ok, err = svc.information.GetInformationStatistic(fmt.Sprint(resp.Id))
	if !ok && err != nil {
		log.Log.Error("information_trace: get information statistic fail, id:%d, err:%s", resp.Id, err)
	}

	if ok {
		resp.FabulousNum = svc.information.Statistic.FabulousNum
		resp.CommentNum = svc.information.Statistic.CommentNum
		resp.BrowseNum = svc.information.Statistic.BrowseNum
		resp.ShareNum = svc.information.Statistic.ShareNum
	}

	now := int(time.Now().Unix())
	// 增加资讯浏览总数
	if err := svc.information.UpdateInformationBrowseNum(resp.Id, now, 1); err != nil {
		log.Log.Errorf("information_trace: update video browse num err:%s", err)
	}

	if userId == "" {
		log.Log.Error("information_trace: user no login")
		return errdef.SUCCESS, resp
	}

	// 获取用户信息
	if user := svc.user.FindUserByUserid(userId); user != nil {

		// 用户是否浏览过
		browse := svc.information.GetUserBrowseInformation(userId, consts.TYPE_INFORMATION, resp.Id)
		if browse != nil {
			svc.information.Browse.CreateAt = now
			svc.information.Browse.UpdateAt = now
			// 已有浏览记录 更新用户浏览的时间
			if err := svc.information.UpdateUserBrowseInformation(userId, consts.TYPE_INFORMATION, resp.Id); err != nil {
				log.Log.Errorf("information_trace: update user browse record err:%s", err)
			}
		} else {
			svc.information.Browse.CreateAt = now
			svc.information.Browse.UpdateAt = now
			svc.information.Browse.UserId = userId
			svc.information.Browse.ComposeId = resp.Id
			svc.information.Browse.ComposeType = consts.TYPE_INFORMATION
			// 添加用户浏览的资讯记录
			if err := svc.information.RecordUserBrowseRecord(); err != nil {
				log.Log.Errorf("information_trace: record user browse record err:%s", err)
			}
		}
	}

	// 是否关注
	if attentionInfo := svc.attention.GetAttentionInfo(userId, resp.UserId); attentionInfo != nil {
		resp.IsAttention = attentionInfo.Status
	}

	// 获取点赞的信息
	if likeInfo := svc.like.GetLikeInfo(userId, resp.Id, consts.LIKE_TYPE_INFORMATION); likeInfo != nil {
		resp.IsLike = likeInfo.Status
	}

	return errdef.SUCCESS, resp
}

//func (svc *InformationModule) GetBodyContent(content string) (string, error) {
//	type body struct {
//		Content string `xml:",innerxml"`
//	}
//
//	type html struct {
//		Body body `xml:"body"`
//	}
//
//	h := html{}
//	err := xml.NewDecoder(bytes.NewBuffer([]byte(content))).Decode(&h)
//	if err != nil {
//		fmt.Println("error", err)
//		return "", err
//	}
//
//	return h.Body.Content, nil
//}
