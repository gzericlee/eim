package auth

import (
	"eim/internal/auth/internal/basic"
	"eim/internal/auth/internal/sso"
	"eim/model"
)

type mode string

const (
	SSOMode   mode = "sso"
	BasicMode mode = "basic"
)

type authenticator interface {
	CheckToken(token string) (*model.User, error)
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
