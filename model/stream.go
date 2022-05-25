package model

import "eim/pkg/json"

//InternalStream 内部流转的数据流
type InternalStream struct {
	UserId   string `json:"userId"`
	DeviceId string `json:"deviceId"`
	Topic    string `json:"topic"`
	Body     []byte `json:"body"`
}

func (its *InternalStream) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *InternalStream) Serialize() ([]byte, error) {
	return json.Marshal(its)
}
