package model

import (
	"time"

	"github.com/gzericlee/eim/internal/model"
)

type GatewayArgs struct {
	Gateway    *model.Gateway
	Expiration time.Duration
}

type GatewaysReply struct {
	Gateways []*model.Gateway
}
