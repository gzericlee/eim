package router

import (
	"github.com/gin-gonic/gin"

	"eim/internal/redis"
)

func RegisterAPIRoutes(engine *gin.Engine, redisManager *redis.Manager) {
	regUploadAPIs(engine, redisManager)
}
