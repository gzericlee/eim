package mq

type Channel string

const (
	MessageAckChannel      Channel = "message_ack_channel"
	MessageDispatchChannel Channel = "message_dispatch_channel"
)
