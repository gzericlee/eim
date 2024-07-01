package protocol

const (
	Message  = 1
	Messages = 2

	Ack       = 100
	Connected = 200
)

var WebsocketCodec *websocketCodec

func init() {
	WebsocketCodec = &websocketCodec{}
}
