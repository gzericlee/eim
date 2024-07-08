package api

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"

	"eim/internal/fileflex/router"
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
	gin.SetMode("release")

	engine := gin.Default()

	engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.DateTime),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.RegisterAPIRoutes(engine, cfg.RedisManager)

	its.server = &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: engine}

	eimmetrics.EnableMetrics(32003)

	return its.server.ListenAndServe()
}

func (its *HttpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return its.server.Shutdown(ctx)
}
