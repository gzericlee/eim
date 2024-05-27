package pprof

import (
	"net"
	"net/http"
	_ "net/http/pprof"

	"go.uber.org/zap"

	"eim/util/log"
)

func EnablePProf() {
	go func() {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Error("Listen failed", zap.Error(err))
			return
		}

		log.Info("PProf server started successfully", zap.String("addr", listener.Addr().String()))

		if err := http.Serve(listener, nil); err != nil {
			log.Error("Could not start PProf server", zap.Error(err))
		}
	}()
}
