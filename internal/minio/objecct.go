package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

func (its *Manager) UploadObject(bucketName, objectName string, reader io.ReadCloser) error {
	_, err := its.minioClient.PutObject(context.Background(), bucketName, objectName, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("add user -> %w", err)
	}
	return err
}

func (its *Manager) ShareObject(bucketName, objectName string, days int) (string, error) {
	expiryDuration := time.Duration(days) * 24 * time.Hour

	presignedUrl, err := its.minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expiryDuration, nil)
	if err != nil {
		return "", fmt.Errorf("generating presigned URL: %w", err)
	}

	return presignedUrl.String(), nil
}
