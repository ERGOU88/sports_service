//@Description 通知
//@Author 凌云（何玉福）heyf@youzu.com
//@DateTime 2020/9/3 5:50 下午

package notify

type notify interface {
	Send() error
}

type content struct {
	ServiceName string
	Content     []byte
	To          string
}
