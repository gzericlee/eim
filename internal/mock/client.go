package mock

import (
	"github.com/lesismal/nbio/nbhttp/websocket"

	"eim/global"
	"eim/internal/protocol"
)

type client struct {
	userId        string
	deviceId      string
	deviceName    string
	deviceType    string
	deviceVersion string
	conn          *websocket.Conn
}

func newUpgrader() *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.SetPongHandler(func(conn *websocket.Conn, s string) {
		hbCount.Add(1)
	})
	u.OnMessage(func(conn *websocket.Conn, messageType websocket.MessageType, data []byte) {
		cmd, _ := protocol.WebsocketCodec.Decode(data)
		switch cmd {
		case protocol.Ack:
			{
				ackCount.Add(1)
				//global.Logger.Debugf("Ack: %v", string(data))
			}
		case protocol.Message:
			{
				msgCount.Add(1)
				//global.Logger.Debugf("Message: %v", string(data))
			}
		}
	})

	u.OnClose(func(conn *websocket.Conn, err error) {
		connectedCount.Add(-1)
		client := conn.Session().(*client)
		global.Logger.Debugf("Device closed: %v - %v，%v", client.userId, client.deviceId, err)
	})

	return u
}
