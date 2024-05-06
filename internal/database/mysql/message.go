package mysql

import "eim/internal/model"

func (its *Repository) SaveMessage(message *model.Message) error {
	return its.db.Save(message).Error
}
