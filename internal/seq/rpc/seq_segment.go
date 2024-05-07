package rpc

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"eim/internal/database"
	"eim/pkg/idgenerator"
	"eim/pkg/log"
)

type segmentSeq struct {
	db    database.IDatabase
	cache sync.Map
}

func (its *segmentSeq) IncrementId(ctx context.Context, req *Request, reply *Reply) error {
	var gen *generator
	if obj, exist := its.cache.Load(req.BizId); exist {
		gen = obj.(*generator)
	} else {
		gen = newGenerator(req.BizId, its.db)
		its.cache.Store(req.BizId, gen)
	}
	reply.Number = gen.Get()
	return nil
}

func (its *segmentSeq) SnowflakeId(ctx context.Context, req *Request, reply *Reply) error {
	reply.Number = idgenerator.NextId()
	return nil
}

type generator struct {
	bizId    string
	ch       chan int64
	min, max int64
	locker   sync.RWMutex
	db       database.IDatabase
}

func newGenerator(bizId string, db database.IDatabase) *generator {
	var gen = &generator{}

	gen.bizId = bizId
	gen.db = db
	gen.ch = make(chan int64, 1)
	go gen.generate()

	return gen
}

func (its *generator) Get() int64 {
	select {
	case id := <-its.ch:
		return id
	}
}

func (its *generator) generate() {
	_ = its.reload()
	for {
		if its.min >= its.max {
			_ = its.reload()
		}
		its.min++
		its.ch <- its.min
	}
}

func (its *generator) reload() error {
	its.locker.Lock()
	defer its.locker.Unlock()
	for {
		seg, err := its.db.GetSegment(its.bizId)
		if err != nil {
			log.Error("Get segment failed", zap.String("bizId", its.bizId), zap.Error(err))
			time.Sleep(time.Second)
			continue
		}
		its.min = seg.MaxId
		its.max = seg.MaxId + int64(seg.Step)
		log.Debug("Reload new seq segment", zap.String("bizId", its.bizId), zap.Int64("min", its.min), zap.Int64("max", its.max))
		return nil
	}
}
