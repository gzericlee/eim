package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) (*Repository, error) {
	repository := &Repository{db: db}
	err := repository.initIndexes()
	return repository, err
}
