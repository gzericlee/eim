package service

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	rpcmodel "github.com/gzericlee/eim/internal/seq/rpc/model"

	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
	"github.com/gzericlee/eim/pkg/log"
	"github.com/gzericlee/eim/pkg/snowflake"
)

type SeqService struct {
	generator  *snowflake.Generator
	segmentRpc *storagerpc.SegmentClient
	cache      sync.Map
}

func NewSeqService(segmentRpc *storagerpc.SegmentClient, generator *snowflake.Generator) *SeqService {
	return &SeqService{
		generator:  generator,
		segmentRpc: segmentRpc,
	}
}

func (its *SeqService) IncrId(ctx context.Context, req *rpcmodel.Request, reply *rpcmodel.Reply) error {
	var incr *incrementer
	if obj, exist := its.cache.Load(req.BizId); exist {
		incr = obj.(*incrementer)
	} else {
		incr = newIncrementer(req.BizId, req.TenantId, its.segmentRpc)
		its.cache.Store(req.BizId, incr)
	}
	reply.Number = incr.Get()
	return nil
}

func (its *SeqService) SnowflakeId(ctx context.Context, req *rpcmodel.Request, reply *rpcmodel.Reply) error {
	reply.Number = its.generator.NextId()
	return nil
}

type incrementer struct {
	bizId      string
	tenantId   string
	ch         chan int64
	min, max   int64
	locker     sync.RWMutex
	segmentRpc *storagerpc.SegmentClient
}

func newIncrementer(bizId, tenantId string, storageRpc *storagerpc.SegmentClient) *incrementer {
	var gen = &incrementer{}

	gen.bizId = bizId
	gen.tenantId = tenantId
	gen.segmentRpc = storageRpc
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
		seg, err := its.segmentRpc.GetSegment(its.bizId, its.tenantId)
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
