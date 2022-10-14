package maindb

import "eim/internal/types"

func (its *tidbRepository) SaveMessage(message *types.Message) error {
	return tidb.Save(message).Error
}
