package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type SegmentClient struct {
	*rpcxclient.XClientPool
}

func (its *SegmentClient) GetSegment(bizId, tenantId string) (*model.Segment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.SegmentReply{}
	err := its.Get().Call(ctx, "GetSegment", &rpcmodel.SegmentArgs{BizId: bizId, TenantId: tenantId}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetSegment -> %w", err)
	}

	return reply.Segment, nil
}
