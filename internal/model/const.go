package model

const (
	AndroidDevice = "android"
	IPhoneDevice  = "iphone"
	WindowsDevice = "windows"
	LinuxDevice   = "linux"
	WebDevice     = "web"
)

const (
	BizUser    = 1
	BizGroup   = 2
	BizService = 3
)

const (
	Offline = 0
	Online  = 1
	Logout  = -1
)

const (
	Enabled  = 1
	Disabled = 0
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
