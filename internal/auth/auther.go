package auth

import (
	"github.com/gzericlee/eim/internal/auth/basic"
	"github.com/gzericlee/eim/internal/auth/oauth2"

	"github.com/gzericlee/eim/internal/model"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
	oauth2lib "github.com/gzericlee/eim/pkg/oauth2"
)

type Mode string

const (
	OAuth2Mode Mode = "auth2"
	BasicMode  Mode = "basic"
)

type IAuther interface {
	CheckToken(token string) (*model.Biz, error)
}

func NewAuther(mode Mode, oauth2Client oauth2lib.Client, bizRpc *storagerpc.BizClient) IAuther {
	switch mode {
	case OAuth2Mode:
		{
			return oauth2.NewAuther(oauth2Client, bizRpc)
		}
	case BasicMode:
		{
			return basic.NewAuther(bizRpc)
		}
	default:
		{
			return nil
		}
	}
}
