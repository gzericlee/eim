package fileflex

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"

	authrpc "eim/internal/auth/rpc"
	"eim/internal/fileflex/router"
	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/log"
	eimmetrics "eim/pkg/metrics"
	"eim/pkg/middleware"
)

type Config struct {
	Ip            string
	Port          int
	AuthRpc       *authrpc.Client
	StorageRpc    *storagerpc.Client
	MinioEndpoint string
}

type HttpServer struct {
	server *http.Server
}

func (its *HttpServer) Run(cfg Config) error {
	gin.SetMode("release")

	engine := gin.New()

	ginMiddleware := middleware.NewGinMiddleware(cfg.AuthRpc)

	engine.Use(gin.Recovery(), ginMiddleware.LogFormatter(), ginMiddleware.Auth)

	router.RegisterAPIRoutes(engine, cfg.StorageRpc, cfg.MinioEndpoint)

	routeInfo := engine.Routes()
	for _, ri := range routeInfo {
		log.Info(fmt.Sprintf("%-6s %-25s --> %s", ri.Method, ri.Path, ri.Handler))
	}

	its.server = &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: engine}

	eimmetrics.EnableMetrics(32003)

	return its.server.ListenAndServe()
}

func (its *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return its.server.Shutdown(ctx)
}
