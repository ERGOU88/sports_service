//@Description todo
//@Author 凌云（何玉福）heyf@youzu.com
//@DateTime 2020/9/3 7:15 下午

package notify

import "testing"

func TestEmail_Send(t *testing.T) {
	em := Email{
		content: content{
			To:          "heyf@youzu.com",
			Content:     []byte("test"),
			ServiceName: "test",
		},
	}
	t.Log(em.Send())
}
