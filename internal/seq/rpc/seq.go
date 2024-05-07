package rpc

import "context"

type Request struct {
	BizId string
}

type Reply struct {
	Number int64
}

type Seq interface {
	IncrementId(ctx context.Context, req *Request, reply *Reply) error
	SnowflakeId(ctx context.Context, req *Request, reply *Reply) error
}
