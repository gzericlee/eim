package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rcrowley/go-metrics"
	"go.uber.org/zap"

	"github.com/gzericlee/eim/pkg/log"
	"github.com/gzericlee/eim/pkg/netutil"
)

func EnableMetrics(port int) {
	var err error
	port, err = netutil.IncrAvailablePort(port)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("[::]:%d", port)

	metrics.RegisterRuntimeMemStats(metrics.DefaultRegistry)
	go metrics.CaptureRuntimeMemStats(metrics.DefaultRegistry, time.Second)

	go func() {
		http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
			metrics.WriteJSONOnce(metrics.DefaultRegistry, writer)
		})

		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Error("Could not start Metrics server", zap.Error(err))
		}
	}()

	log.Info("Metrics server started successfully", zap.String("addr", addr))
}
