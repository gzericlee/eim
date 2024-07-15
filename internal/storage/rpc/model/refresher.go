package model

import "github.com/gzericlee/eim/internal/model"

const (
	ActionInsert = "insert"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

type RefreshDevicesArgs struct {
	Key    string
	Device *model.Device
	Action string
}

type RefreshBizArgs struct {
	Key    string
	Biz    *model.Biz
	Action string
}

type RefreshTenantArgs struct {
	Key    string
	Tenant *model.Tenant
	Action string
}

type RefreshBizMembersArgs struct {
	Key       string
	BizMember *model.BizMember
	Action    string
}
