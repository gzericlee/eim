package basic

import (
	"encoding/base64"
	"fmt"
	"strings"

	"eim/internal/redis"
	"eim/internal/types"
)

type Authenticator struct {
}

func (its *Authenticator) CheckToken(token string) (*types.User, error) {
	c, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return nil, fmt.Errorf("invalid token")
	}
	passwd := cs[s+1:]

	var loginId, company string
	uc := strings.IndexByte(cs[:s], '@')
	if uc < 0 {
		loginId = cs[:s]
	} else {
		loginId = cs[:s][:uc]
		company = cs[:s][uc+1:]
	}

	user, err := redis.GetUser(loginId, company)
	if err != nil {
		return nil, err
	}

	if user.Password != passwd {
		return nil, fmt.Errorf("password is incorrect")
	}

	return user, nil
}
