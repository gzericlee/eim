package minio

import (
	"fmt"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
	UseSSL          bool
	ExpiryDays      int
}

type Manager struct {
	minioClient *minio.Client
	adminClient *madmin.AdminClient

	endpoint    string
	accessKeyId string
	secretKey   string
	expiryDays  int
}

func NewManager(cfg *Config) (*Manager, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyId, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("new minio client -> %w", err)
	}

	adminClient, err := madmin.NewWithOptions(cfg.Endpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyId, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("new madmin client -> %w", err)

	}

	return &Manager{
		minioClient: minioClient,
		adminClient: adminClient,
		endpoint:    cfg.Endpoint,
		accessKeyId: cfg.AccessKeyId,
		secretKey:   cfg.SecretAccessKey,
		expiryDays:  cfg.ExpiryDays,
	}, nil
}
