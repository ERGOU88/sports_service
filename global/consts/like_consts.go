package consts

// 1 视频点赞/视频评论、回复中@ 2 帖子点赞/帖子评论、回复中@ 3 视频评论点赞/视频评论、回复 4 帖子评论点赞/帖子评论、回复 5 发布帖子时内容带有 @用户
const (
	TYPE_VIDEOS        = 1
	TYPE_POSTS         = 2
	TYPE_VIDEO_COMMENT = 3
	TYPE_POST_COMMENT  = 4
	TYPE_PUBLISH_POST  = 5
)

// 1 已点赞 0 未点赞
const (
	ALREADY_GIVE_LIKE  = 1
	NOT_GIVE_LIKE      = 0
)


const (
	// 资讯点赞
	LIKE_TYPE_INFORMATION         = 5
	// 资讯评论点赞
	LIKE_TYPE_INFORMATION_COMMENT = 6
)
