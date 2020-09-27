//@Description todo
//@Author 凌云（何玉福）heyf@youzu.com
//@DateTime 2020/9/3 7:38 下午

package notify

import "testing"

func TestSms_Send(t *testing.T) {
	s := &Sms{}
	s.To = "17721473132"
	s.ServiceName = "test"
	s.Content = []byte("test")
	t.Log(s.Send())
}
