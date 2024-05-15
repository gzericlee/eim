package mongodb

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) (*Repository, error) {
	repository := &Repository{db: db}
	err := repository.initIndexes()
	if err != nil {
		return nil, fmt.Errorf("init indexes -> %w", err)
	}
	return repository, nil
}
