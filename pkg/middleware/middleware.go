package middleware

import (
	authrpc "github.com/gzericlee/eim/internal/auth/rpc/client"
)

type ginMiddleware struct {
	authRpc *authrpc.AuthClient
}

func NewGinMiddleware(authRpc *authrpc.AuthClient) *ginMiddleware {
	return &ginMiddleware{authRpc: authRpc}
}
