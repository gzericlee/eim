package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
	"github.com/gzericlee/eim/pkg/log"
)

type GatewayHandler struct {
	gatewayRpc *storagerpc.GatewayClient
}

func NewGatewayHandler(gatewayRpc *storagerpc.GatewayClient) *GatewayHandler {
	return &GatewayHandler{gatewayRpc: gatewayRpc}
}

func (its *GatewayHandler) List(c *gin.Context) {
	gateways, err := its.gatewayRpc.GetGateways()
	if err != nil {
		log.Error("Error get gateways", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gateways)
}
