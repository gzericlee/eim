package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"eim/internal/redis"
	"eim/pkg/log"
)

type GatewayHandler struct {
	RedisManager *redis.Manager
}

func (its *GatewayHandler) List(c *gin.Context) {
	gateways, err := its.RedisManager.GetGateways()
	if err != nil {
		log.Error("Error get gateways", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gateways)
}
