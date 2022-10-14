package types

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
	DeviceStoreTopic          = "device_store"
	MessageGroupDispatchTopic = "message_group_dispatch"
	MessageUserDispatchTopic  = "message_user_dispatch"
)

const (
	DeviceStoreChannel  = "device_store_channel"
	MessageStoreChannel = "message_store_channel"

	AckDispatchChannel     = "ack_dispatch_channel"
	MessageDispatchChannel = "message_dispatch_channel"
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
