package pgsql

import (
	"fmt"

	"eim/internal/model"
)

func (its *Repository) InsertDevice(device *model.Device) error {
	_, err := its.db.Insert(device)
	if err != nil {
		return fmt.Errorf("insert device -> %w", err)
	}
	return nil
}

func (its *Repository) UpdateDevice(device *model.Device) error {
	_, err := its.db.Where("user_id = ? AND tenant_id = ? AND device_id = ?", device.UserId, device.TenantId, device.DeviceId).Update(device)
	if err != nil {
		return fmt.Errorf("update device -> %w", err)
	}
	return nil
}

func (its *Repository) GetDevice(userId, tenantId, deviceId string) (*model.Device, error) {
	var device = &model.Device{}
	_, err := its.db.Where("user_id = ? AND tenant_id = ? AND device_id = ?", userId, tenantId, deviceId).Get(device)
	if err != nil {
		return nil, fmt.Errorf("select device -> %w", err)
	}
	return device, nil
}

func (its *Repository) GetDevicesByUser(userId, tenantId string) ([]*model.Device, error) {
	var devices []*model.Device
	err := its.db.Where("user_id = ? AND tenant_id = ?", userId, tenantId).Find(&devices)
	if err != nil {
		return nil, fmt.Errorf("select devices -> %w", err)
	}
	return devices, nil
}

func (its *Repository) DeleteDevice(userId, tenantId, deviceId string) error {
	_, err := its.db.Where("user_id = ? AND tenant_id = ? AND device_id = ?", userId, tenantId, deviceId).Delete()
	if err != nil {
		return fmt.Errorf("delete device -> %w", err)
	}
	return nil
}

func (its *Repository) ListDevices(filter map[string]interface{}, order []string, limit, offset int64) ([]*model.Device, int64, error) {
	var devices []*model.Device
	query := its.db.Where("")

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	for _, by := range order {
		query = query.OrderBy(by)
	}

	total, err := query.Limit(int(limit), int(offset)).FindAndCount(&devices)
	if err != nil {
		return nil, 0, fmt.Errorf("select devices -> %w", err)
	}

	return devices, total, nil
}
