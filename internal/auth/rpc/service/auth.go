package service

import (
	"context"
	"fmt"

	"github.com/gzericlee/eim/internal/auth"
	rpcmodel "github.com/gzericlee/eim/internal/auth/rpc/model"
	rpcclient "github.com/gzericlee/eim/internal/storage/rpc/client"
	oauth2lib "github.com/gzericlee/eim/pkg/oauth2"
)

type AuthService struct {
	authMode     auth.Mode
	oauth2Client oauth2lib.Client
	bizRpc       *rpcclient.BizClient
}

func NewAuthService(authMode auth.Mode, oauth2Client oauth2lib.Client, bizRpc *rpcclient.BizClient) *AuthService {
	return &AuthService{
		authMode:     authMode,
		oauth2Client: oauth2Client,
		bizRpc:       bizRpc,
	}
}

func (its *AuthService) CheckToken(ctx context.Context, args *rpcmodel.AuthArgs, reply *rpcmodel.AuthReply) error {
	authenticator := auth.NewAuther(its.authMode, its.oauth2Client, its.bizRpc)

	biz, err := authenticator.CheckToken(args.Token)
	if err != nil {
		return fmt.Errorf("check token -> %w", err)
	}

	reply.Biz = biz

	return nil
}
