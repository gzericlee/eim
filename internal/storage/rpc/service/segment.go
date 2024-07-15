package service

import (
	"context"
	"fmt"

	"github.com/gzericlee/eim/internal/database"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type SegmentService struct {
	database database.IDatabase
}

func NewSegmentService(database database.IDatabase) *SegmentService {
	return &SegmentService{
		database: database,
	}
}

func (its *SegmentService) GetSegment(ctx context.Context, args *rpcmodel.SegmentArgs, reply *rpcmodel.SegmentReply) error {
	segment, err := its.database.GetSegment(args.BizId, args.TenantId)
	if err != nil {
		return fmt.Errorf("get segment -> %w", err)
	}

	reply.Segment = segment

	return nil
}
