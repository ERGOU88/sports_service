package cinformation

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"sports_service/server/dao"
	"sports_service/server/global/backend/errdef"
	"sports_service/server/global/backend/log"
	"sports_service/server/models/mattention"
	"sports_service/server/models/mcollect"
	"sports_service/server/models/minformation"
	"sports_service/server/models/mlike"
	"sports_service/server/models/muser"
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
	condition := ""

	list, err := svc.information.GetInformationList(condition, offset, size)
	if err != nil {
		return errdef.INFORMATION_LIST_FAIL, []*minformation.InformationResp{}
	}

	if len(list) == 0 {
		return errdef.SUCCESS, []*minformation.InformationResp{}
	}

	resp := make([]*minformation.InformationResp, len(list))
	for index, information := range list {
		info := &minformation.InformationResp{
			Id: information.Id,
			Cover: information.Cover,
			Title: information.Title,
			CreateAt: information.CreateAt,
			//JumpUrl: information.JumpUrl,
			UserId: information.UserId,
			Content: information.Content,
			Describe: information.Describe,
			PubType: information.PubType,
			Status: information.Status,
		}

		if user := svc.user.FindUserByUserid(info.UserId); user != nil {
			info.Avatar = user.Avatar
			info.NickName = user.NickName
		}

		ok, err := svc.information.GetInformationStatistic(fmt.Sprint(info.Id))
		if !ok && err != nil {
			log.Log.Error("information_trace: get information statistic fail, id:%d, err:%s", info.Id, err)
		}

		if ok {
			info.FabulousNum = svc.information.Statistic.FabulousNum
			info.CommentNum = svc.information.Statistic.CommentNum
			info.ShareNum = svc.information.Statistic.ShareNum
			info.BrowseNum = svc.information.Statistic.BrowseNum
		}

		resp[index] = info
	}

	return errdef.SUCCESS, resp
}

func (svc *InformationModule) DeleteInformation(id string) int {
	cols := "status"
	condition := fmt.Sprintf("id=%s", id)
	svc.information.Information.Status = 3
	if _, err := svc.information.UpdateInfo(condition, cols); err != nil {
		return errdef.INFORMATION_DELETE_FAIL
	}

	return errdef.SUCCESS
}