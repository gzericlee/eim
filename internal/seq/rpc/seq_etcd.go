package rpc

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"

	"eim/pkg/idgenerator"
	"eim/util/log"
)

type etcdSeq struct {
	etcdClient *clientv3.Client
}

func (its *etcdSeq) IncrementId(ctx context.Context, req *Request, reply *Reply) error {
	bizId := req.BizId
	if bizId == "" {
		return fmt.Errorf("bizId can't be empty")
	}

	session, err := concurrency.NewSession(its.etcdClient, concurrency.WithTTL(10))
	if err != nil {
		return fmt.Errorf("new etcd session -> %w", err)
	}
	defer session.Close()

	mutex := concurrency.NewMutex(session, fmt.Sprintf("%s/%s", basePath, bizId))
	for i := 0; i < 3; i++ {
		err = mutex.Lock(ctx)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return fmt.Errorf("new etcd locker -> %w", err)
	}
	defer mutex.Unlock(ctx)

	resp, err := its.etcdClient.Get(ctx, fmt.Sprintf("%s/%s", basePath, bizId))
	if err != nil {
		return fmt.Errorf("etcd get -> %w", err)
	}

	var seq int64
	if resp.Count == 0 {
		seq = 1
	} else {
		seq = resp.Kvs[0].Version + 1
	}

	_, err = its.etcdClient.Put(ctx, fmt.Sprintf("%s/%s", basePath, bizId), "")
	if err != nil {
		return fmt.Errorf("etcd put -> %w", err)
	}

	reply.Number = seq

	log.Debug("Get seq", zap.String("bizId", bizId), zap.Int64("seq", seq))

	return nil
}

func (its *etcdSeq) SnowflakeId(ctx context.Context, req *Request, reply *Reply) error {
	reply.Number = idgenerator.NextId()
	return nil
}
