package handler

import (
	"github.com/gin-gonic/gin"

	"eim/internal/redis"
)

type UploadHandler struct {
	RedisManager *redis.Manager
}

func (its *UploadHandler) Upload(c *gin.Context) {

}
