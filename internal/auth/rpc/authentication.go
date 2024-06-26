package rpc

import (
	"context"
	"fmt"

	"eim/internal/config"
	"eim/internal/model"
	storagerpc "eim/internal/storage/rpc"
)

type Request struct {
	Token string
}

type Reply struct {
	Biz *model.Biz
}

type Authentication struct {
	StorageRpc *storagerpc.Client
}

func (its *Authentication) CheckToken(ctx context.Context, req *Request, reply *Reply) error {
	authenticator := NewAuthenticator(Mode(config.SystemConfig.AuthSvr.Mode), its.StorageRpc)
	biz, err := authenticator.CheckToken(req.Token)
	if err != nil {
		return fmt.Errorf("check token -> %w", err)
	}
	reply.Biz = biz
	return nil
}
