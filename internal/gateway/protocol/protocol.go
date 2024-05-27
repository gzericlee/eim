package protocol

const (
	Message        = 10
	OfflineMessage = 20
	Ack            = 200
)

var WebsocketCodec *websocketCodec

func init() {
	WebsocketCodec = &websocketCodec{}
}
