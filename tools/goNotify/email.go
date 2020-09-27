//@Description todo
//@Author 凌云（何玉福）heyf@youzu.com
//@DateTime 2020/9/3 5:53 下午

package notify

import (
	"gopkg.in/gomail.v2"
)

var (
	host     = "u2mail.uuzu.com"
	username = "platform_monitor"
	password = "62hJkC5F"
)

type Email struct {
	content
}

func (e Email) Send() error {
	em := gomail.NewMessage()
	em.SetBody("text/html", string(e.Content))
	em.SetHeader("From", "alert@youzu.com")
	em.SetHeader("To", e.To)
	em.SetHeader("Subject", "告警通知")
	d := gomail.NewDialer(host, 25, username, password)
	if err := d.DialAndSend(em); err != nil {
		return err
	}
	return nil

}
