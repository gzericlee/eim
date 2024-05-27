package stats

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"eim/pkg/cache"
)

var m sync.Map

type Node struct {
	Evicted, Updated, Added, GetMiss, GetHit, DelMiss, DelHit uint64
}

func (s *Node) HitRate() float64 {
	if s.GetHit == 0 && s.GetMiss == 0 {
		return 0.0
	}
	return float64(s.GetHit) / float64(s.GetHit+s.GetMiss)
}

func Bind(pool string, caches ...*cache.Cache) error {
	v, _ := m.LoadOrStore(pool, &Node{})
	for _, c := range caches {
		c.Inspect(func(action int, _ string, _ *interface{}, _ []byte, status int) {
			atomic.AddUint64((*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(v.(*Node)))+uintptr(status+action*2-1)*unsafe.Sizeof(&status))), 1)
		})
	}
	return nil
}

func All() *sync.Map {
	return &m
}
