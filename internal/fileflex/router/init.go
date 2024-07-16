package router

import (
	"github.com/gin-gonic/gin"

	seqrpc "github.com/gzericlee/eim/internal/seq/rpc/client"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

func RegisterAPIRoutes(engine *gin.Engine, tenantRpc *storagerpc.TenantClient, seqRpc *seqrpc.SeqClient, fileRpc *storagerpc.FileClient, minioEndpoint, externalServiceEndpoint string) {
	regUploadAPIs(engine, tenantRpc, seqRpc, fileRpc, minioEndpoint, externalServiceEndpoint)
}
