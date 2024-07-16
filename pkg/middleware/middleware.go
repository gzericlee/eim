package middleware

import (
	authrpc "github.com/gzericlee/eim/internal/auth/rpc/client"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

type ginMiddleware struct {
	tenantRpc *storagerpc.TenantClient
	authRpc   *authrpc.AuthClient
}

func NewGinMiddleware(authRpc *authrpc.AuthClient, tenantRpc *storagerpc.TenantClient) *ginMiddleware {
	return &ginMiddleware{authRpc: authRpc, tenantRpc: tenantRpc}
}
