package model

import (
	"time"

	"eim/pkg/json"
)

type Device struct {
	DeviceId      string     `json:"deviceId" bson:"deviceId"`
	UserId        string     `json:"userId" bson:"userId"`
	DeviceType    string     `json:"deviceType" bson:"deviceType"`
	DeviceVersion string     `json:"deviceVersion" bson:"deviceVersion"`
	GatewayIp     string     `json:"gateway_ip" bson:"gatewayIp"`
	OnlineAt      *time.Time `json:"onlineAt" bson:"onlineAt"`
	OfflineAt     *time.Time `json:"offlineAt" bson:"offlineAt"`
	State         int        `json:"state" bson:"state"`
}

func (its *Device) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *Device) Serialize() ([]byte, error) {
	return json.Marshal(its)
}
