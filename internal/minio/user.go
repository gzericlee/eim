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

func (its *Manager) RemoveUser(userName string) error {
	err := its.adminClient.RemoveUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("remove user -> %w", err)
	}

	return nil
}
