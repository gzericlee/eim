package rpc

import (
	"context"

	"eim/internal/config"
	"eim/internal/model"
	"eim/internal/redis"
)

type Request struct {
	Token string
}

type Reply struct {
	User *model.User
}

type Authentication struct {
	RedisManager *redis.Manager
}

func (its *Authentication) CheckToken(ctx context.Context, req *Request, reply *Reply) error {
	authenticator := NewAuthenticator(Mode(config.SystemConfig.AuthSvr.Mode), its.RedisManager)
	user, err := authenticator.CheckToken(req.Token)
	reply.User = user
	return err
}
