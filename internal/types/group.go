package types

import "eim/pkg/json"

type GroupMember struct {
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
}

func (its *GroupMember) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *GroupMember) Serialize() ([]byte, error) {
	return json.Marshal(its)
}
