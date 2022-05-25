package maindb

import (
	"eim/model"
)

func (its *tidbRepository) SaveDevice(device *model.Device) error {
	return tidb.Save(device).Error
}
