package minio

import (
	"context"
	"fmt"
)

func (its *Manager) CreateUser(userName, password string) error {
	err := its.adminClient.AddUser(context.Background(), userName, password)
	if err != nil {
		return fmt.Errorf("add user -> %w", err)
	}

	return nil
}
