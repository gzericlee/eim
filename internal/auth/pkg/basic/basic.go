package basic

import (
	"encoding/base64"
	"fmt"
	"strings"

	"eim/internal/model"
	storagerpc "eim/internal/storage/rpc"
)

type Authenticator struct {
	StorageRpc *storagerpc.Client
}

func (its *Authenticator) CheckToken(token string) (*model.Biz, error) {
	c, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("decode token -> %w", err)
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return nil, fmt.Errorf("invalid token")
	}
	passwd := cs[s+1:]

	var bizId, tenantId string
	uc := strings.IndexByte(cs[:s], '@')
	if uc < 0 {
		bizId = cs[:s]
	} else {
		bizId = cs[:s][:uc]
		tenantId = cs[:s][uc+1:]
	}

	biz, err := its.StorageRpc.GetBiz(bizId, tenantId)
	if err != nil {
		return nil, fmt.Errorf("get user -> %w", err)
	}

	if biz.Attributes != nil {
		if password := biz.Attributes["password"]; password != "" {
			if password == passwd {
				return biz, nil
			}
			return nil, fmt.Errorf("password is incorrect")
		}
	}

	return nil, fmt.Errorf("password is not set")
}
