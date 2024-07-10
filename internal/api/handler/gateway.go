package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/log"
)

type GatewayHandler struct {
	storageRpc *storagerpc.Client
}

func NewGatewayHandler(storageRpc *storagerpc.Client) *GatewayHandler {
	return &GatewayHandler{storageRpc: storageRpc}
}

func (its *GatewayHandler) List(c *gin.Context) {
	gateways, err := its.storageRpc.GetGateways()
	if err != nil {
		log.Error("Error get gateways", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gateways)
}
