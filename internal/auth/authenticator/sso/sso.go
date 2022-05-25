package sso

import "eim/model"

type Authenticator struct {
	endpoint     string
	clientId     string
	clientSecret string
	resourceId   string
}

func (its *Authenticator) VerificationToken(token string) (*model.User, error) {
	//TODO 具体的SSO认证逻辑
	return nil, nil
}
