package websocket

import (
	"time"

	"github.com/lesismal/nbio/nbhttp/websocket"
	"go.uber.org/zap"

	"eim/pkg/log"
)

func (its *Server) heartbeat(conn *websocket.Conn, s string) {
	err := conn.SetReadDeadline(time.Now().Add(its.keepaliveTime))
	if err != nil {
		log.Error("Error set read deadline", zap.Error(err))
		_ = conn.Close()
		return
	}
	err = conn.WriteMessage(websocket.PongMessage, []byte(time.Now().String()))
	if err != nil {
		log.Error("Error send pong message", zap.Error(err))
		_ = conn.Close()
		return
	}
	its.IncrHeartbeatTotal(1)
}
