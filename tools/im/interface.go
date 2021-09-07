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
}
