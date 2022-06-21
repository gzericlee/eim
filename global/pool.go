package global

import (
	"runtime"

	"github.com/lesismal/nbio/taskpool"
)

var SystemPool *systemPool

type systemPool struct {
	*taskpool.FixedPool
}

func init() {
	SystemPool = &systemPool{}
	SystemPool.FixedPool = taskpool.NewFixedPool(runtime.NumCPU()*100, 1024)
}
