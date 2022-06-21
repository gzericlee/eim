package model

import (
	"time"

	"eim/pkg/json"
)

type Seq struct {
	Id       string    `json:"id"`
	MaxId    int64     `json:"maxId"`
	Step     int       `json:"step"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}

func (its *Seq) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *Seq) Serialize() ([]byte, error) {
	return json.Marshal(its)
}
