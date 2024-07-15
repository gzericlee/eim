package oauth2

import (
	"fmt"

	"github.com/gzericlee/eim/internal/model"
	storagerpc "github.com/gzericlee/eim/internal/storage/rpc/client"
	"github.com/gzericlee/eim/pkg/maputil"
	"github.com/gzericlee/eim/pkg/oauth2"
)

type Auther struct {
	bizRpc       *storagerpc.BizClient
	oauth2Client oauth2.Client
}

func NewAuther(client oauth2.Client, bizRpc *storagerpc.BizClient) *Auther {
	return &Auther{oauth2Client: client, bizRpc: bizRpc}
}

func (its *Auther) CheckToken(token string) (*model.Biz, error) {
	userInfo, err := its.oauth2Client.GetUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("get oauth2 user info -> %w", err)
	}

	userId := maputil.GetString(userInfo, "userId", "")
	tenantId := maputil.GetString(userInfo, "tenant", "")

	biz, err := its.bizRpc.GetBiz(userId, tenantId)
	if err != nil {
		return nil, fmt.Errorf("get biz -> %w", err)
	}

	return biz, nil
}
