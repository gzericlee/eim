package model

import "github.com/gzericlee/eim/internal/model"

type DeviceArgs struct {
	Device *model.Device
}

type UserArgs struct {
	UserId   string
	TenantId string
	DeviceId string
}

type DevicesReply struct {
	Devices []*model.Device
}

type DeviceReply struct {
	Device *model.Device
}
