package dispatch

import (
	"sync/atomic"

	"github.com/rcrowley/go-metrics"
)

var (
	msgTotal        = &atomic.Int64{}
	offlineMsgTotal = &atomic.Int64{}
	onlineMsgTotal  = &atomic.Int64{}
	savedMsgTotal   = &atomic.Int64{}
)

func init() {
	metrics.Register("dispatch_msg_total", metrics.NewFunctionalGauge(func() int64 {
		return msgTotal.Load()
	}))
	metrics.Register("dispatch_offline_msg_total", metrics.NewFunctionalGauge(func() int64 {
		return offlineMsgTotal.Load()
	}))
	metrics.Register("dispatch_online_msg_total", metrics.NewFunctionalGauge(func() int64 {
		return onlineMsgTotal.Load()
	}))
	metrics.Register("dispatch_saved_msg_total", metrics.NewFunctionalGauge(func() int64 {
		return savedMsgTotal.Load()
	}))
}
