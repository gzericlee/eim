package router

import (
	"github.com/gin-gonic/gin"

	"eim/internal/fileflex/handler"
	storagerpc "eim/internal/storage/rpc"
)

func regUploadAPIs(engine *gin.Engine, storageRpc *storagerpc.Client, minioEndpoint string) {
	uploadHandler := handler.NewUploadHandler(storageRpc, minioEndpoint)

	upload := engine.Group("/upload")
	{
		upload.POST("", uploadHandler.Upload)
	}
}
