package oauth2

import (
	"fmt"
	"time"

	"eim/internal/model"
	"eim/util/log"
)

type Authenticator struct {
	endpoint     string
	clientId     string
	clientSecret string
	resourceId   string
}

func (its *Authenticator) CheckToken(token string) (*model.Biz, error) {
	now := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("Function time duration %v", time.Since(now)))
	}()

	//TODO 具体的SSO认证逻辑
	return nil, nil
}
