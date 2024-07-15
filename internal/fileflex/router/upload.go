package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gzericlee/eim/internal/fileflex/handler"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

func regUploadAPIs(engine *gin.Engine, tenantRpc *storagerpc.TenantClient, minioEndpoint string) {
	uploadHandler := handler.NewUploadHandler(tenantRpc, minioEndpoint)

	upload := engine.Group("/upload")
	{
		upload.POST("", uploadHandler.Upload)
	}
}
