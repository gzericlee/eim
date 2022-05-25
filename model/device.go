package model

import (
	"time"

	"eim/pkg/json"
)

type Device struct {
	DeviceId      string     `gorm:"primary_key;column:device_id" json:"deviceId"`
	UserId        string     `gorm:"index;column:user_id" json:"userId"`
	DeviceType    string     `gorm:"column:device_type" json:"deviceType"`
	DeviceVersion string     `gorm:"column:device_version" json:"deviceVersion"`
	GatewayIp     string     `gorm:"column:gateway_ip" json:"gateway_ip"`
	OnlineAt      *time.Time `gorm:"column:online_at" json:"onlineAt"`
	OfflineAt     *time.Time `gorm:"column:offline_at" json:"offlineAt"`
	State         int        `gorm:"column:state" json:"state"`
}

func (its *Device) Deserialize(data []byte) error {
	return json.Unmarshal(data, &its)
}

func (its *Device) Serialize() ([]byte, error) {
	return json.Marshal(its)
}

func (its *Device) TableName() string {
	return "eim_device"
}
