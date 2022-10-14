package maindb

import "eim/internal/types"

func (its *tidbRepository) SaveDevice(device *types.Device) error {
	return tidb.Save(device).Error
}
