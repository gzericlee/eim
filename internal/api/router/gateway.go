package router

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"eim/internal/api/handler"
	"eim/internal/model"
	"eim/internal/redis"
)

func regGatewayAPIs(redisManager *redis.Manager) *restful.WebService {
	tags := []string{"网关"}

	ws := &restful.WebService{}
	ws.Path("/gateway").Produces(restful.MIME_JSON, restful.MIME_OCTET).Doc("网关管理")

	gatewayHandler := handler.GatewayHandler{RedisManager: redisManager}

	ws.Route(ws.POST("/list").To(gatewayHandler.List).
		Doc("获取网关列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]*model.Gateway{}))

	return ws
}
