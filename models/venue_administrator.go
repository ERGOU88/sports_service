package models

type VenueAdministrator struct {
	Id        int    `json:"id" xorm:"not null pk autoincr INT(10)"`
	UserId    string `json:"user_id" xorm:"default '' VARCHAR(255)"`
	JobNumber string `json:"job_number" xorm:"default '' comment('工号') VARCHAR(255)"`
	Mobile    int    `json:"mobile" xorm:"default 0 comment('手机号') INT(11)"`
	Name      string `json:"name" xorm:"default '' comment('用户名称') VARCHAR(255)"`
	Username  string `json:"username" xorm:"not null default '' comment('账号') VARCHAR(255)"`
	Password  string `json:"password" xorm:"not null default '' comment('密码') VARCHAR(255)"`
	Status    int    `json:"status" xorm:"default 0 comment('状态') TINYINT(1)"`
	CreateAt  int    `json:"create_at" xorm:"default 0 comment('创建时间') INT(11)"`
	UpdateAt  int    `json:"update_at" xorm:"default 0 comment('更新时间') INT(11)"`
}
