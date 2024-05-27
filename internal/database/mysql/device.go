package mysql

import (
	"fmt"

	"eim/internal/model"
)

func (its *Repository) SaveDevice(device *model.Device) error {
	err := its.db.Save(device).Error
	if err != nil {
		return fmt.Errorf("save device -> %w", err)
	}
	return nil
}

func (its *Repository) GetDevice(userId, deviceId string) (*model.Device, error) {
	var device *model.Device
	err := its.db.Model(&model.Device{}).Where("user_id = ? AND device_id = ?", userId, deviceId).First(&device).Error
	if err != nil {
		return nil, fmt.Errorf("find one device -> %w", err)
	}
	return device, nil
}

func (its *Repository) GetDevices(userId string) ([]*model.Device, error) {
	var devices []*model.Device
	err := its.db.Where("user_id = ?", userId).Find(&devices).Error
	if err != nil {
		return nil, fmt.Errorf("find user devices -> %w", err)
	}
	return devices, nil
}
