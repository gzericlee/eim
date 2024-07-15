package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	rpcmodel "github.com/gzericlee/eim/internal/auth/rpc/model"
	"github.com/gzericlee/eim/internal/model"
)

type AuthClient struct {
	*rpcxclient.XClientPool
}

func (its *AuthClient) CheckToken(token string) (*model.Biz, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.AuthReply{}
	err := its.Get().Call(ctx, "CheckToken", &rpcmodel.AuthArgs{Token: token}, reply)
	if err != nil {
		return nil, fmt.Errorf("call CheckToken -> %w", err)
	}

	return reply.Biz, nil
}
