package handler

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"

	"eim/internal/redis"
)

type GatewayHandler struct {
}

func (its *GatewayHandler) List(request *restful.Request, response *restful.Response) {
	gateways, err := redis.GetGateways()
	if err != nil {
		_ = response.WriteError(http.StatusInternalServerError, err)
		return
	}
	_ = response.WriteAsJson(gateways)
}
