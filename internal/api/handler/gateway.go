package handler

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"go.uber.org/zap"

	"eim/internal/redis"
	"eim/pkg/log"
)

type GatewayHandler struct {
	RedisManager *redis.Manager
}

func (its *GatewayHandler) List(request *restful.Request, response *restful.Response) {
	gateways, err := its.RedisManager.GetGateways()
	if err != nil {
		log.Error("Error get gateways", zap.Error(err))
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}
	_ = response.WriteAsJson(gateways)
}
