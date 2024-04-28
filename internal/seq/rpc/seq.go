package rpc

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"

	"eim/pkg/log"
)

type Request struct {
	Id string
}

type Reply struct {
	Number int64
}

type Seq struct {
	etcdClient *clientv3.Client
}

func (its *Seq) Number(ctx context.Context, req *Request, reply *Reply) error {
	if req.Id == "" {
		return fmt.Errorf("id is empty")
	}

	session, err := concurrency.NewSession(its.etcdClient, concurrency.WithTTL(10))
	if err != nil {
		return fmt.Errorf("failed to create etcd session: %v", err)
	}
	defer session.Close()

	mutex := concurrency.NewMutex(session, fmt.Sprintf("/%s/%s", basePath, req.Id))
	for i := 0; i < 3; i++ {
		err = mutex.Lock(ctx)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %v", err)
	}
	defer mutex.Unlock(ctx)

	resp, err := its.etcdClient.Get(ctx, fmt.Sprintf("/%s/%s", basePath, req.Id))
	if err != nil {
		return fmt.Errorf("failed to get seq: %v", err)
	}

	var seq int64
	if resp.Count == 0 {
		seq = 1
	} else {
		seq = resp.Kvs[0].Version + 1
	}

	_, err = its.etcdClient.Put(ctx, fmt.Sprintf("/%s/%s", basePath, req.Id), "")
	if err != nil {
		return fmt.Errorf("failed to update seq: %v", err)
	}

	reply.Number = seq

	log.Info("Get seq", zap.String("id", req.Id), zap.Int64("seq", seq))

	return nil
}
