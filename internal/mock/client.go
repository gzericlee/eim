package mock

import (
	"encoding/json"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"github.com/gzericlee/eim/internal/gateway/protocol"
	"github.com/gzericlee/eim/internal/model"
	"github.com/gzericlee/eim/pkg/log"
)

type Client struct {
	token         string
	userId        string
	deviceId      string
	deviceName    string
	deviceType    string
	deviceVersion string
	conn          *websocket.Conn
	connected     chan bool
}

func (its *Server) upgrade(cli *Client) *websocket.Upgrader {
	u := websocket.NewUpgrader()
	u.SetPongHandler(func(conn *websocket.Conn, s string) {
		its.hbCount.Add(1)
	})

	var ackHandler = func(conn *websocket.Conn, msg *model.Message) {
		err := conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(protocol.Ack, []byte(strconv.FormatInt(msg.MsgId, 10))))
		if err != nil {
			log.Error("Error sending ack", zap.Error(err))
			return
		}
		its.ackCount.Add(1)
	}

	u.OnMessage(func(conn *websocket.Conn, messageType websocket.MessageType, data []byte) {
		cmd, data := protocol.WebsocketCodec.Decode(data)
		switch cmd {
		case protocol.Ack:
			{
				its.ackCount.Add(1)
			}
		case protocol.Connected:
			{
				cli.connected <- true
			}
		case protocol.Message:
			{
				msg := &model.Message{}
				err := proto.Unmarshal(data, msg)
				if err != nil {
					log.Error("Error unmarshal message", zap.Error(err))
					its.invalidCount.Add(1)
					return
				}
				ackHandler(conn, msg)
				its.msgCount.Add(1)
			}
		case protocol.Messages:
			{
				var msgs []*model.Message
				err := json.Unmarshal(data, &msgs)
				if err != nil {
					log.Error("Error unmarshal messages", zap.Error(err))
					its.invalidCount.Add(1)
					return
				}
				for _, msg := range msgs {
					ackHandler(conn, msg)
				}
				its.msgCount.Add(int64(len(msgs)))
			}
		}
	})

	u.OnClose(func(conn *websocket.Conn, err error) {
		its.connectedCount.Add(-1)
		log.Error("Connection closed", zap.Error(err))
	})

	return u
}
