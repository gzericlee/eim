package redis

import (
	"fmt"

	"eim/internal/model"
)

func (its *Manager) SaveUser(user *model.User) error {
	key := fmt.Sprintf("%v@%v:info", user.LoginId, user.TenantId)
	body, _ := user.Serialize()
	return its.rdsClient.Set(key, body, 0)
}

func (its *Manager) GetUser(loginId, tenantId string) (*model.User, error) {
	key := fmt.Sprintf("%v@%v:info", loginId, tenantId)
	result, err := its.rdsClient.Get(key)
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = user.Deserialize([]byte(result))
	return user, err
}
