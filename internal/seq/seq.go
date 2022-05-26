package seq

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"eim/global"
	"eim/internal/redis"
)

type Id struct {
	userId   string
	ch       chan int64
	min, max int64
	locker   sync.RWMutex
}

func newId(userId string) *Id {
	var id = &Id{}
	id.locker.Lock()
	defer id.locker.Unlock()

	var key = "seq_" + userId

	if obj, exist := global.SystemCache.Get(key); exist {
		id = obj.(*Id)
	} else {
		id.userId = userId
		id.ch = make(chan int64, 1)
		go id.generate()
		global.SystemCache.Save(key, id)
	}

	return id
}

func (its *Id) Get() int64 {
	select {
	case id := <-its.ch:
		return id
	}
}

func (its *Id) generate() {
	_ = its.reload()
	for {
		if its.min >= its.max {
			_ = its.reload()
		}
		its.min++
		its.ch <- its.min
	}
}

func (its *Id) reload() error {
	its.locker.Lock()
	defer its.locker.Unlock()
	for {
		seq, err := redis.GetSegmentSeq(its.userId)
		if err != nil {
			global.Logger.Error("Error getting Seq", zap.String("userId", its.userId), zap.Error(err))
			time.Sleep(time.Second)
			continue
		}
		its.min = seq.MaxId
		its.max = seq.MaxId + int64(seq.Step)
		global.Logger.Info("Reload new seq segment", zap.String("userId", its.userId), zap.Int64("min", its.min), zap.Int64("max", its.max))
		return nil
	}
}
