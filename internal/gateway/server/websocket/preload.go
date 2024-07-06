package websocket

import "eim/internal/gateway/session"

func (its *Server) preload(sess *session.Session) {
	if sess == nil {
		return
	}

	user := sess.GetUser()

	//缓存预热
	_, _ = its.storageRpc.GetDevices(user.BizId, user.TenantId)
	_, _ = its.seqRpc.IncrId(user.BizId, user.TenantId)
}
