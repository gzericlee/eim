package api

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/emicklei/go-restful/v3"
	"go.uber.org/zap"

	"eim/internal/api/router"
	"eim/pkg/log"
)

type Config struct {
	Port int
}

type HttpServer struct {
	server *http.Server
}

func (its *HttpServer) Run(cfg Config) error {
	err := router.RegisterAPIRoute()
	if err != nil {
		return err
	}

	its.server = &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: restful.DefaultContainer}

	log.Info("ApiServer starting", zap.String("port", strconv.Itoa(cfg.Port)))

	return its.server.ListenAndServe()
}

func (its *HttpServer) Stop() error {
	return its.server.Shutdown(nil)
}
