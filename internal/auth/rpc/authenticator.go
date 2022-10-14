package rpc

import (
	"eim/internal/auth/pkg/basic"
	"eim/internal/auth/pkg/sso"
	"eim/internal/types"
)

type mode string

const (
	SSOMode   mode = "sso"
	BasicMode mode = "basic"
)

type authenticator interface {
	CheckToken(token string) (*types.User, error)
}

func newAuthenticator(mode mode) authenticator {
	//TODO 参数
	switch mode {
	case SSOMode:
		{
			return &sso.Authenticator{}
		}
	case BasicMode:
		{
			return &basic.Authenticator{}
		}
	default:
		{
			return nil
		}
	}
}
