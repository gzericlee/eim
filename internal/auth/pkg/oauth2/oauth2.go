package oauth2

import (
	"eim/internal/model"
)

type Authenticator struct {
	endpoint     string
	clientId     string
	clientSecret string
	resourceId   string
}

func (its *Authenticator) CheckToken(token string) (*model.Biz, error) {
	//TODO 具体的SSO认证逻辑
	return nil, nil
}
