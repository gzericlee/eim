package seq

import (
	"context"
	"fmt"

	"github.com/smallnest/rpcx/server"
)

type Request struct {
	UserId string
}

type Reply struct {
	Id int64
}

func id(ctx context.Context, req *Request, reply *Reply) error {
	reply.Id = newId(req.UserId).Get()
	return nil
}

func InitSeqServer(ip string, port int) error {
	s := server.NewServer()
	err := s.RegisterFunction("seq.id.service", id, "")
	if err != nil {
		return err
	}
	err = s.Serve("tcp", fmt.Sprintf("%v:%v", ip, port))
	return err
}
