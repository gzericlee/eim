package router

import (
	"github.com/gin-gonic/gin"

	"eim/internal/minio"
	storagerpc "eim/internal/storage/rpc"
)

func RegisterAPIRoutes(engine *gin.Engine, storageRpc *storagerpc.Client, minioManager *minio.Manager) {
	regGatewayAPIs(engine, storageRpc)
	regTenantAPIs(engine, storageRpc, minioManager)
}
