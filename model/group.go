package model

import "eim/pkg/json"

type Group struct {
	GroupId   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

type GroupMember struct {
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
	SeqId   int64  `json:"seqId"`
}

func (its *Group) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *Group) Serialize() ([]byte, error) {
	return json.Marshal(its)
}

func (its *GroupMember) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *GroupMember) Serialize() ([]byte, error) {
	return json.Marshal(its)
}
