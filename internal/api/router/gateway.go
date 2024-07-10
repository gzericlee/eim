package router

import (
	"github.com/gin-gonic/gin"

	"eim/internal/api/handler"
	storagerpc "eim/internal/storage/rpc"
)

func regGatewayAPIs(engine *gin.Engine, storageRpc *storagerpc.Client) {
	gatewayHandler := handler.NewGatewayHandler(storageRpc)

	gateways := engine.Group("/gateways")
	{
		gateways.GET("/", gatewayHandler.List)
	}
}
