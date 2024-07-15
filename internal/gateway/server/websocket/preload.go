package websocket

import (
	"go.uber.org/zap"

	"github.com/gzericlee/eim/internal/gateway/session"
	"github.com/gzericlee/eim/pkg/log"
)

func (its *Server) preload(sess *session.Session) {
	if sess == nil {
		return
	}

	user := sess.GetUser()

	//缓存预热
	_, err := its.deviceRpc.GetDevices(user.BizId, user.TenantId)
	if err != nil {
		log.Error("Error get devices", zap.Error(err))
	}

	_, _ = its.seqRpc.IncrId(user.BizId, user.TenantId)
}
