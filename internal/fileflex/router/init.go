package router

import (
	"github.com/gin-gonic/gin"

	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

func RegisterAPIRoutes(engine *gin.Engine, tenantRpc *storagerpc.TenantClient, minioEndpoint string) {
	regUploadAPIs(engine, tenantRpc, minioEndpoint)
}
