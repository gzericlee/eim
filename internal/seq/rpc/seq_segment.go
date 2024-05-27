package rpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	storagerpc "eim/internal/storage/rpc"
	"eim/pkg/snowflake"
	"eim/util/log"
)

type segmentSeq struct {
	generator  *snowflake.Generator
	storageRpc *storagerpc.Client
	cache      sync.Map
}

func (its *segmentSeq) IncrementId(ctx context.Context, req *Request, reply *Reply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	var incr *incrementer
	if obj, exist := its.cache.Load(req.BizId); exist {
		incr = obj.(*incrementer)
	} else {
		incr = newIncrementer(req.BizId, its.storageRpc)
		its.cache.Store(req.BizId, incr)
	}
	reply.Number = incr.Get()
	return nil
}

func (its *segmentSeq) SnowflakeId(ctx context.Context, req *Request, reply *Reply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	reply.Number = its.generator.NextId()
	return nil
}

type incrementer struct {
	bizId      string
	ch         chan int64
	min, max   int64
	locker     sync.RWMutex
	storageRpc *storagerpc.Client
}

func newIncrementer(bizId string, storageRpc *storagerpc.Client) *incrementer {
	var gen = &incrementer{}

	gen.bizId = bizId
	gen.storageRpc = storageRpc
	gen.ch = make(chan int64, 1)
	go gen.generate()

	return gen
}

func (its *incrementer) Get() int64 {
	select {
	case id := <-its.ch:
		return id
	}
}

func (its *incrementer) generate() {
	_ = its.reload()
	for {
		if its.min >= its.max {
			_ = its.reload()
		}
		its.min++
		its.ch <- its.min
	}
}

func (its *incrementer) reload() error {
	its.locker.Lock()
	defer its.locker.Unlock()
	for {
		seg, err := its.storageRpc.GetSegment(its.bizId)
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
