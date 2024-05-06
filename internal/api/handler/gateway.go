package handler

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"

	"eim/internal/redis"
)

type GatewayHandler struct {
	RedisManager *redis.Manager
}

func (its *GatewayHandler) List(request *restful.Request, response *restful.Response) {
	gateways, err := its.RedisManager.GetGateways()
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}
	_ = response.WriteAsJson(gateways)
}
