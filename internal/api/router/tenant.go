package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gzericlee/eim/internal/api/handler"
	"github.com/gzericlee/eim/internal/minio"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

func regTenantAPIs(engine *gin.Engine, tenantRpc *storagerpc.TenantClient, minioManager *minio.Manager) {
	tenantHandler := handler.NewTenantHandler(tenantRpc, minioManager)

	tenants := engine.Group("/tenants")
	{
		tenants.POST("", tenantHandler.Register)
		tenants.PUT("/:tenantId", tenantHandler.Update)
		tenants.PATCH("/:tenantId/enable_fileflex", tenantHandler.EnableFileFlex)
	}
}
