package redis

import (
	"fmt"

	"eim/internal/types"
	"eim/pkg/json"
)

func SaveGroupMember(member *types.GroupMember) error {
	var users []string
	key := fmt.Sprintf("%v:members", member.GroupId)

	result, _ := rdsClient.Get(key)
	if result != "" {
		err := json.Unmarshal([]byte(result), &users)
		if err != nil {
			return err
		}
	}

	users = append(users, member.UserId)

	body, err := json.Marshal(users)
	if err != nil {
		return err
	}

	return rdsClient.Set(key, body, 0)
}

func GetGroupMembers(groupId string) ([]string, error) {
	key := fmt.Sprintf("%v:members", groupId)
	result, err := rdsClient.Get(key)
	if err != nil {
		return nil, err
	}
	var users []string
	err = json.Unmarshal([]byte(result), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
