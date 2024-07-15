package pprof

import (
	"net"
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"

	"github.com/gzericlee/eim/pkg/log"
)

func EnablePProf() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Error("Listen failed", zap.Error(err))
		return
	}

	go func() {
		if err := http.Serve(listener, nil); err != nil {
			log.Error("Could not start PProf server", zap.Error(err))
		}
	}()

	log.Info("PProf server started successfully", zap.String("addr", listener.Addr().String()))
}
