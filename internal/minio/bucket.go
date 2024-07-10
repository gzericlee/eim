package minio

import (
	"context"
	"fmt"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
)

func (its *Manager) CreateBucket(name string) error {
	err := its.minioClient.MakeBucket(context.Background(), name, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("make bucket -> %w", err)
	}

	return nil
}

func (its *Manager) RemoveBucket(name string) error {
	err := its.minioClient.RemoveBucket(context.Background(), name)
	if err != nil {
		return fmt.Errorf("remove bucket -> %w", err)
	}

	return nil
}

func (its *Manager) AttachBucketPolicy(bucketName, userName string) error {
	ctx := context.Background()
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Action": ["s3:*"],
			"Resource": ["arn:aws:s3:::%s/*"]
		}]
	}`, bucketName)

	policyName := fmt.Sprintf("%s:%s:rw", bucketName, userName)

	err := its.adminClient.AddCannedPolicy(ctx, policyName, []byte(policy))
	if err != nil {
		return fmt.Errorf("add canned policy -> %w", err)
	}

	_, err = its.adminClient.AttachPolicy(ctx, madmin.PolicyAssociationReq{
		Policies: []string{policyName},
		User:     userName,
	})
	if err != nil {
		return fmt.Errorf("attach policy -> %w", err)
	}

	return nil
}

func (its *Manager) DetachBucketPolicy(bucketName, userName string) error {
	ctx := context.Background()

	policyName := fmt.Sprintf("%s:%s:rw", bucketName, userName)

	_, err := its.adminClient.DetachPolicy(ctx, madmin.PolicyAssociationReq{
		Policies: []string{policyName},
		User:     userName,
	})
	if err != nil {
		return fmt.Errorf("detach policy -> %w", err)
	}

	return nil
}
