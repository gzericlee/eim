package router

import (
	"github.com/gin-gonic/gin"

	storagerpc "eim/internal/storage/rpc"
)

func RegisterAPIRoutes(engine *gin.Engine, storageRpc *storagerpc.Client, minioEndpoint string) {
	regUploadAPIs(engine, storageRpc, minioEndpoint)
}
