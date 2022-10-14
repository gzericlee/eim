package sso

import (
	"eim/internal/types"
)

type Authenticator struct {
	endpoint     string
	clientId     string
	clientSecret string
	resourceId   string
}

func (its *Authenticator) CheckToken(token string) (*types.User, error) {
	//TODO 具体的SSO认证逻辑
	return nil, nil
}
