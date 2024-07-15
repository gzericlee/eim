package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gzericlee/eim/internal/minio"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

func RegisterAPIRoutes(engine *gin.Engine, gatewayRpc *storagerpc.GatewayClient, tenantRpc *storagerpc.TenantClient, minioManager *minio.Manager) {
	regGatewayAPIs(engine, gatewayRpc)
	regTenantAPIs(engine, tenantRpc, minioManager)
}
