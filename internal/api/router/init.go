package router

import (
	"fmt"

	"github.com/emicklei/go-restful/v3"

	"eim/internal/api/filter"
	"eim/internal/redis"
	"eim/pkg/log"
)

func RegisterAPIRoutes(redisManager *redis.Manager) error {
	restful.DefaultContainer.Add(regGatewayAPIs(redisManager))

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{""},
		AllowedHeaders: []string{"Content-Type", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "token"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer,
	}

	wss := restful.DefaultContainer.RegisteredWebServices()
	for _, ws := range wss {
		log.Info(fmt.Sprintf("------------------------------> %v <------------------------------", ws.Documentation()))
		for _, route := range ws.Routes() {
			log.Info(route.String())
		}
	}

	restful.Filter(cors.Filter)
	restful.Filter(restful.OPTIONSFilter())
	restful.Filter(filter.LogFormat())

	return nil
}
