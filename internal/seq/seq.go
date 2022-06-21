package seq

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"eim/global"
	"eim/internal/redis"
)

type seq struct {
	id       string
	ch       chan int64
	min, max int64
	locker   sync.RWMutex
}

func newSeq(id string) *seq {
	var s = &seq{}
	s.locker.Lock()
	defer s.locker.Unlock()

	var key = "seq_" + id

	if obj, exist := global.SystemCache.Get(key); exist {
		s = obj.(*seq)
	} else {
		s.id = id
		s.ch = make(chan int64, 1)
		go s.generate()
		global.SystemCache.Save(key, s)
	}

	return s
}

func (its *seq) Get() int64 {
	select {
	case id := <-its.ch:
		return id
	}
}

func (its *seq) generate() {
	_ = its.reload()
	for {
		if its.min >= its.max {
			_ = its.reload()
		}
		its.min++
		its.ch <- its.min
	}
}

func (its *seq) reload() error {
	its.locker.Lock()
	defer its.locker.Unlock()
	for {
		seq, err := redis.GetSegmentSeq(its.id)
		if err != nil {
			global.Logger.Error("Error getting Seq", zap.String("id", its.id), zap.Error(err))
			time.Sleep(time.Second)
			continue
		}
		its.min = seq.MaxId
		its.max = seq.MaxId + int64(seq.Step)
		global.Logger.Info("Reload new seq segment", zap.String("id", its.id), zap.Int64("min", its.min), zap.Int64("max", its.max))
		return nil
	}
}
