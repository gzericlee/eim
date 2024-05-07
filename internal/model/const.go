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
	ToUser           = 1
	ToGroup          = 2
	ToServiceAccount = 3

	FromUser           = 1
	FromServiceAccount = 1
)

const (
	TextMessage = 1
)

const (
	BizGroup   = "group"
	BizUser    = "user"
	BizService = "service"
)
