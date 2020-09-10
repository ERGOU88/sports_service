package rdskey

import (
	"fmt"
)

// make redis key
func MakeKey(key_fmt string, keys ...interface{}) string {
	return fmt.Sprintf(key_fmt, keys...)
}
