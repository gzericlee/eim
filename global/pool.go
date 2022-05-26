package global

import (
	"runtime"

	"github.com/lesismal/nbio/taskpool"
)

var SystemPool *systemPool

type systemPool struct {
	*taskpool.FixedNoOrderPool
}

func init() {
	SystemPool = &systemPool{}
	SystemPool.FixedNoOrderPool = taskpool.NewFixedNoOrderPool(runtime.NumCPU(), 1024)
}
