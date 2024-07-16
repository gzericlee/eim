package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gzericlee/eim/internal/fileflex/handler"
	seqrpc "github.com/gzericlee/eim/internal/seq/rpc/client"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

func regUploadAPIs(engine *gin.Engine, tenantRpc *storagerpc.TenantClient, seqRpc *seqrpc.SeqClient, fileRpc *storagerpc.FileClient, minioEndpoint, externalServiceEndpoint string) {
	uploadHandler := handler.NewUploadHandler(tenantRpc, seqRpc, fileRpc, minioEndpoint, externalServiceEndpoint)

	upload := engine.Group("/upload")
	{
		upload.POST("", uploadHandler.Upload)
	}
}
