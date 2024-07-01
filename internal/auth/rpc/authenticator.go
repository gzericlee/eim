package rpc

import (
	"eim/internal/auth/pkg/basic"
	"eim/internal/auth/pkg/oauth2"
	"eim/internal/model"
	storagerpc "eim/internal/storage/rpc"
)

type Mode string

const (
	OAuth2Mode Mode = "auth2"
	BasicMode  Mode = "basic"
)

type IAuthenticator interface {
	CheckToken(token string) (*model.Biz, error)
}

func NewAuthenticator(mode Mode, storageRpc *storagerpc.Client) IAuthenticator {
	//TODO 参数
	switch mode {
	case OAuth2Mode:
		{
			return &oauth2.Authenticator{}
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
