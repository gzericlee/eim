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
	DeviceStoreTopic = "device_store"

	MessageDispatchTopic = "message_dispatch"
	MessageSendTopic     = "message_send"
	MessageStoreTopic    = "message_store"
)

const (
	DeviceStoreChannel  = "device_store_channel"
	MessageStoreChannel = "message_store_channel"

	AckDispatchChannel     = "ack_dispatch_channel"
	MessageDispatchChannel = "message_dispatch_channel"
)

const (
	ToUser    = 1
	ToGroup   = 2
	ToService = 3

	FromUser    = 1
	FromService = 1
)

const (
	TextMessage = 1
)
