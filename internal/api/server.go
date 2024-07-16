package api

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gzericlee/eim/internal/api/router"
	authrpc "github.com/gzericlee/eim/internal/auth/rpc/client"
	"github.com/gzericlee/eim/internal/minio"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
	"github.com/gzericlee/eim/pkg/log"
	eimmetrics "github.com/gzericlee/eim/pkg/metrics"
	"github.com/gzericlee/eim/pkg/middleware"
)

type Config struct {
	Ip           string
	Port         int
	AuthRpc      *authrpc.AuthClient
	MinioManager *minio.Manager
	GatewayRpc   *storagerpc.GatewayClient
	TenantRpc    *storagerpc.TenantClient
}

type HttpServer struct {
	server *http.Server
}

func (its *HttpServer) Run(cfg *Config) error {
	gin.SetMode("release")

	engine := gin.New()

	ginMiddleware := middleware.NewGinMiddleware(cfg.AuthRpc, cfg.TenantRpc)

	engine.Use(gin.Recovery(), ginMiddleware.LogFormatter(), ginMiddleware.Auth)

	router.RegisterAPIRoutes(engine, cfg.GatewayRpc, cfg.TenantRpc, cfg.MinioManager)

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
