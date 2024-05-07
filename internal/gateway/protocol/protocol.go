package protocol

const (
	Message = 10
	Ack     = 200
)

var WebsocketCodec *websocketCodec

func init() {
	WebsocketCodec = &websocketCodec{}
}
