package router

import (
	"github.com/gin-gonic/gin"

	"eim/internal/fileflex/handler"
	"eim/internal/redis"
)

func regUploadAPIs(engine *gin.Engine, redisManager *redis.Manager) {
	uploadHandler := handler.UploadHandler{RedisManager: redisManager}

	upload := engine.Group("/upload")
	{
		upload.POST("", uploadHandler.Upload)
	}
}
