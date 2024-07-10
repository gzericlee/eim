package router

import (
	"github.com/gin-gonic/gin"

	"eim/internal/api/handler"
	"eim/internal/minio"
	storagerpc "eim/internal/storage/rpc"
)

func regTenantAPIs(engine *gin.Engine, storageRpc *storagerpc.Client, minioManager *minio.Manager) {
	tenantHandler := handler.NewTenantHandler(storageRpc, minioManager)

	tenants := engine.Group("/tenants")
	{
		tenants.POST("/", tenantHandler.Register)
		tenants.PUT("/:tenantId", tenantHandler.Update)
		tenants.PATCH("/:tenantId/enable_fileflex", tenantHandler.EnableFileFlex)
	}
}
