package mysql

import (
	"gorm.io/gorm"

	"eim/internal/model"
)

type Repository struct {
	db *gorm.DB
}

func (its *Repository) SaveMessages(messages []*model.Message) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
