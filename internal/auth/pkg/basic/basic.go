package basic

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"eim/internal/model"
	storagerpc "eim/internal/storage/rpc"
	"eim/util/log"
)

type Authenticator struct {
	StorageRpc *storagerpc.Client
}

func (its *Authenticator) CheckToken(token string) (*model.User, error) {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

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

	var loginId, tenantId string
	uc := strings.IndexByte(cs[:s], '@')
	if uc < 0 {
		loginId = cs[:s]
	} else {
		loginId = cs[:s][:uc]
		tenantId = cs[:s][uc+1:]
	}

	user, err := its.StorageRpc.GetUser(loginId, tenantId)
	if err != nil {
		return nil, fmt.Errorf("get user -> %w", err)
	}

	if user.Password != passwd {
		return nil, fmt.Errorf("password is incorrect")
	}

	return user, nil
}
