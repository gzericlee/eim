package router

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"eim/internal/api/handler"
	"eim/internal/types"
)

func regExampleAPIs() *restful.WebService {
	tags := []string{"网关"}

	ws := &restful.WebService{}
	ws.Path("/gateway").Produces(restful.MIME_JSON, restful.MIME_OCTET).Doc("范例")

	gatewayHandler := handler.GatewayHandler{}

	ws.Route(ws.POST("/list").To(gatewayHandler.List).
		Doc("获取网关列表").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]*types.Gateway{}))

	return ws
}
