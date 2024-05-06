package rpc

import "context"

type Request struct {
	Id string
}

type Reply struct {
	Number int64
}

type Seq interface {
	Number(ctx context.Context, req *Request, reply *Reply) error
}
