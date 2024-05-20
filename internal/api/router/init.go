package router

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"

	"eim/internal/api/filter"
	"eim/internal/redis"
	"eim/util/log"
)

func RegisterAPIRoutes(redisManager *redis.Manager) error {
	restful.DefaultContainer.Add(regGatewayAPIs(redisManager))

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject,
	}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("swagger-ui/dist"))))

	log.Info("Get the API using /apidocs.json")
	log.Info("Open Swagger UI using /apidocs")

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{""},
		AllowedHeaders: []string{"Content-Type", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "token"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer,
	}

	restful.Filter(cors.Filter)
	restful.Filter(restful.OPTIONSFilter())
	restful.Filter(filter.NCSALogFormat())

	return nil
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Enterprise Instant Messaging",
			Description: "EIM Api服务",
			Version:     "1.0.0",
		},
	}
}
