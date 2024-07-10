package middleware

import authrpc "eim/internal/auth/rpc"

type ginMiddleware struct {
	authRpc *authrpc.Client
}

func NewGinMiddleware(authRpc *authrpc.Client) *ginMiddleware {
	return &ginMiddleware{authRpc: authRpc}
}
