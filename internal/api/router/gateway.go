package router

import (
	"github.com/gin-gonic/gin"

	"eim/internal/api/handler"
	"eim/internal/redis"
)

func regGatewayAPIs(engine *gin.Engine, redisManager *redis.Manager) {
	gatewayHandler := handler.GatewayHandler{RedisManager: redisManager}

	gateways := engine.Group("/gateways")
	{
		gateways.GET("/", gatewayHandler.List)
	}
}
