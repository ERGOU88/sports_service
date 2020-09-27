//@Description todo
//@Author 凌云（何玉福）heyf@youzu.com
//@DateTime 2020/9/3 7:36 下午

package notify

import "testing"

func TestFeishu_Send(t *testing.T) {
	s := Feishu{
		content: content{
			To:          "heyf",
			Content:     []byte("test"),
			ServiceName: "test",
		},
	}
	t.Log(s.Send())
}
