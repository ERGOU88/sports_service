package im

func Init() ImInterface {
	Im = NewImRealize()
	return Im
}

type ImInterface interface {
	/*
		腾讯云注册用户
		@param userId 用户id
		@param name 昵称
		@param avatar 头像
	*/
	AddUser(userId, name, avatar string) (string, error) //ok

	/*
	    腾讯云创建群组
	    @param groupType 群组类型 必填  Private 同新版本中的 Work（好友工作群）Public ChatRoom，同新版本中的 Meeting（临时会议群）AVChatRoom 直播
	    @param owner 群主
	    @param name 必填 群名称
	    @param introduction 群简介 最长240字节 1个汉字3字节
	    @param notification 群公共 最长300字节 1个汉字3字节
	    @param faceUrl 群头像 URL，最长100字节
	*/
	CreateGroup(groupType, owner, name, introduction, notification, faceUrl string) (string, error)
}
