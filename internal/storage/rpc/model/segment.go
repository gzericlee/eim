package model

import "github.com/gzericlee/eim/internal/model"

type SegmentArgs struct {
	BizId    string
	TenantId string
}

type SegmentReply struct {
	Segment *model.Segment
}
