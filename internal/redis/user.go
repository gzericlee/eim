package redis

import (
	"fmt"

	"eim/model"
)

func SaveUser(user *model.User) error {
	key := fmt.Sprintf("%v@%v:info", user.LoginId, user.Company)
	body, _ := user.Serialize()
	return rdsClient.Set(key, body, 0)
}

func GetUser(loginId, company string) (*model.User, error) {
	key := fmt.Sprintf("%v@%v:info", loginId, company)
	result, err := rdsClient.Get(key)
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	err = user.Deserialize([]byte(result))
	return user, err
}
