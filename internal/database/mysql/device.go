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
