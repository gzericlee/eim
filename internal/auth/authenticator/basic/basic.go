package basic

import "eim/model"

type Authenticator struct {
	AdminUser   string
	AdminPasswd string
}

func (its *Authenticator) VerificationToken(token string) (*model.User, error) {
	//TODO 具体的Basic认证逻辑
	return nil, nil
}
