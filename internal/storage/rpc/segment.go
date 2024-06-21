package rpc

import (
	"context"
	"fmt"
	"time"

	"eim/internal/database"
	"eim/internal/model"
	"eim/util/log"
)

type SegmentArgs struct {
	BizId    string
	TenantId string
}

type SegmentReply struct {
	Segment *model.Segment
}

type Segment struct {
	database database.IDatabase
}

func (its *Segment) GetSegment(ctx context.Context, args *SegmentArgs, reply *SegmentReply) error {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	segment, err := its.database.GetSegment(args.BizId, args.TenantId)
	if err != nil {
		return fmt.Errorf("get segment -> %w", err)
	}

	reply.Segment = segment

	return nil
}
