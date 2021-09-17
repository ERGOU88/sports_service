package live

func Init() ILive {
	Live = NewLiveRealize()
	return Live
}

type ILive interface {
	/*
		云直播 生成推流地址
		@param roomId 房间id
	    @param expireTm 过期时长 [秒]
	*/
	GenPushStream(roomId string, expireTm int64) (string, string) //ok

	/*
	    云直播  生成拉流地址
		@param roomId 房间id
		@param expireTm 过期时长 [秒]
	*/
	GenPullStream(roomId string, expireTm int64) *PullStreamInfo
}

