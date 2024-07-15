package router

import (
	"github.com/gin-gonic/gin"

	"github.com/gzericlee/eim/internal/api/handler"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
)

func regGatewayAPIs(engine *gin.Engine, gatewayRpc *storagerpc.GatewayClient) {
	gatewayHandler := handler.NewGatewayHandler(gatewayRpc)

	gateways := engine.Group("/gateways")
	{
		gateways.GET("", gatewayHandler.List)
	}
}
