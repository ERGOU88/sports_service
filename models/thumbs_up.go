package models

type ThumbsUp struct {
	Id       int64  `json:"id" xorm:"pk autoincr comment('主键id') BIGINT(20)"`
	TypeId   int64  `json:"type_id" xorm:"not null comment('作品id （视频id/帖子id/评论id）') index(thumbs_up_id) index BIGINT(20)"`
	UserId   string `json:"user_id" xorm:"not null comment('用户id') index VARCHAR(60)"`
	ToUserId string `json:"to_user_id" xorm:"not null comment('被点赞的用户id') index VARCHAR(60)"`
	ZanType  int    `json:"zan_type" xorm:"not null comment('1 视频点赞 2 帖子点赞 3 评论点赞') index(thumbs_up_id) TINYINT(2)"`
	Status   int    `json:"status" xorm:"not null comment('1赞 2取消点赞') TINYINT(1)"`
	CreateAt int    `json:"create_at" xorm:"not null comment('创建时间') INT(11)"`
}
