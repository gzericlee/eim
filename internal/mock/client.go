package mock

import (
	"encoding/json"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/internal/gateway/protocol"
	"eim/internal/model"
	"eim/util/log"
)

type client struct {
	token         string
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
		cmd, data := protocol.WebsocketCodec.Decode(data)
		switch cmd {
		case protocol.Ack:
			{
				ackCount.Add(1)
			}
		case protocol.OfflineMessage:
			{
				var msgs []*model.Message
				err := json.Unmarshal(data, &msgs)
				if err != nil {
					log.Error("Error unmarshal offline message", zap.Error(err))
					invalidCount.Add(1)
					return
				}
				for _, msg := range msgs {
					err = conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(protocol.Ack, []byte(strconv.FormatInt(msg.MsgId, 10))))
					if err != nil {
						log.Error("Error sending ack", zap.Error(err))
						return
					}
					ackCount.Add(1)
					msgCount.Add(1)
				}
			}
		case protocol.Message:
			{
				msg := &model.Message{}
				err := proto.Unmarshal(data, msg)
				if err != nil {
					log.Error("Error unmarshal message", zap.Error(err))
					invalidCount.Add(1)
					return
				}
				err = conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(protocol.Ack, []byte(strconv.FormatInt(msg.MsgId, 10))))
				if err != nil {
					log.Error("Error sending ack", zap.Error(err))
					return
				}
				msgCount.Add(1)
			}
		}
	})

	u.OnClose(func(conn *websocket.Conn, err error) {
		connectedCount.Add(-1)
	})

	return u
}
