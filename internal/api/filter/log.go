package filter

import (
	"fmt"
	"strings"
	"time"

	"github.com/emicklei/go-restful/v3"

	"eim/pkg/log"
)

func LogFormat() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		start := time.Now()
		chain.ProcessFilter(req, resp)
		tc := time.Since(start)
		log.Info(fmt.Sprintf("%s - \"%s %s %s\" status: %d size: %d time: %v",
			strings.Split(req.Request.RemoteAddr, ":")[0],
			req.Request.Method,
			req.Request.URL.RequestURI(),
			req.Request.Proto,
			resp.StatusCode(),
			resp.ContentLength(),
			tc,
		))
	}
}
