package api

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/emicklei/go-restful/v3"

	"eim/internal/api/router"
	"eim/internal/redis"
	eimmetrics "eim/pkg/metrics"
)

type Config struct {
	Ip           string
	Port         int
	RedisManager *redis.Manager
}

type HttpServer struct {
	server *http.Server
}

func (its *HttpServer) Run(cfg Config) error {
	err := router.RegisterAPIRoutes(cfg.RedisManager)
	if err != nil {
		return fmt.Errorf("register api routes -> %w", err)
	}

	its.server = &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: restful.DefaultContainer}

	eimmetrics.EnableMetrics(32003)

	return its.server.ListenAndServe()
}

func (its *HttpServer) Stop() error {
	return its.server.Shutdown(nil)
}
