package consts

type PLATFORM int32

const (
	IOS_PLATFORM        PLATFORM = iota            // iOS端
	ANDROID_PLATFORM                               // android端
	WEB_PLATFORM                                   // web端
)
