package authenticator

import "eim/model"

type Authenticator interface {
	VerificationToken(token string) (*model.User, error)
}
