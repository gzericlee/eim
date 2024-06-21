package model

const (
	AndroidDevice = "android"
	IPhoneDevice  = "iphone"
	WindowsDevice = "windows"
	LinuxDevice   = "linux"
	WebDevice     = "web"
)

const (
	OfflineState = 0
	OnlineState  = 1
	LogoutState  = -1
)

const (
	ToUser    = 1
	ToGroup   = 2
	ToService = 3

	FromUser    = 1
	FromService = 2
)

const (
	TextMessage     = 1
	ImageMessage    = 2
	FileMessage     = 3
	AudioMessage    = 4
	VideoMessage    = 5
	LocationMessage = 6
	ControlMessage  = 7
)

const (
	BizGroup   = "group"
	BizUser    = "user"
	BizService = "service"
)
