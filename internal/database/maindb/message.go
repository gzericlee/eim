package maindb

import "eim/model"

func (its *tidbRepository) SaveMessage(message *model.Message) error {
	return tidb.Save(message).Error
}
