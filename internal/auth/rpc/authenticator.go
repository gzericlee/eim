package rpc

import (
	"eim/internal/auth/pkg/basic"
	"eim/internal/auth/pkg/sso"
	"eim/internal/model"
	storagerpc "eim/internal/storage/rpc"
)

type Mode string

const (
	SSOMode   Mode = "sso"
	BasicMode Mode = "basic"
)

type Authenticator interface {
	CheckToken(token string) (*model.User, error)
}

func NewAuthenticator(mode Mode, storageRpc *storagerpc.Client) Authenticator {
	//TODO 参数
	switch mode {
	case SSOMode:
		{
			return &sso.Authenticator{}
		}
	case BasicMode:
		{
			return &basic.Authenticator{StorageRpc: storageRpc}
		}
	default:
		{
			return nil
		}
	}
}
