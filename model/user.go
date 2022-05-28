package model

import "eim/pkg/json"

type User struct {
	UserId   string `json:"userId"`
	LoginId  string `json:"loginId"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	Company  string `json:"company"`
	SeqId    int64  `json:"seqId"`
}

func (its *User) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *User) Serialize() ([]byte, error) {
	return json.Marshal(its)
}
