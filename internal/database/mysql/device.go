package mysql

import (
	"eim/internal/model"
)

func (its *Repository) SaveDevice(device *model.Device) error {
	return its.db.Save(device).Error
}
