package pgsql

import (
	"fmt"

	"eim/internal/model"
)

func (its *Repository) InsertMessage(message *model.Message) error {
	_, err := its.db.Insert(message)
	if err != nil {
		return fmt.Errorf("insert message -> %w", err)
	}
	return nil
}

func (its *Repository) InsertMessages(messages []*model.Message) error {
	_, err := its.db.Insert(messages)
	if err != nil {
		return fmt.Errorf("insert messages -> %w", err)
	}
	return nil
}

func (its *Repository) GetMessagesByIds(msgIds []int64) ([]*model.Message, error) {
	var messages []*model.Message
	err := its.db.Where("msg_id IN ?", msgIds).Find(&messages)
	if err != nil {
		return nil, fmt.Errorf("select messages -> %w", err)
	}
	return messages, nil
}

func (its *Repository) ListHistoryMessages(filter map[string]interface{}, order []string, minSeq, maxSeq, limit, offset int64) ([]*model.Message, error) {
	var messages []*model.Message
	query := its.db.Where("")

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	if minSeq > 0 && maxSeq > 0 && maxSeq > minSeq {
		query = query.Where("seq_id >= ? AND seq_id <= ?", minSeq, maxSeq)
	}

	for _, by := range order {
		query = query.OrderBy(by)
	}

	err := query.Limit(int(limit), int(offset)).Find(&messages)
	if err != nil {
		return nil, fmt.Errorf("select messages -> %w", err)
	}
	return messages, nil
}
