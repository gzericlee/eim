package mysql

import (
	"fmt"

	"eim/internal/model"
)

func (its *Repository) SaveMessage(message *model.Message) error {
	err := its.db.Save(message).Error
	if err != nil {
		return fmt.Errorf("save message -> %w", err)
	}
	return nil
}

func (its *Repository) GetMessagesByIds(msgIds []int64) ([]*model.Message, error) {
	//TODO implement me
	panic("implement me")
}
