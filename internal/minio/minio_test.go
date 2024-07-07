package minio

import "testing"

var manager *Manager

func init() {
	var err error
	manager, err = NewManager(&Config{
		Endpoint:        "127.0.0.1:9000",
		AccessKeyId:     "minioadmin",
		SecretAccessKey: "minioadmin",
		UseSSL:          false,
	})
	if err != nil {
		panic(err)
	}
}

func TestManager_CreateUser(t *testing.T) {
	err := manager.CreateUser("lirui", "pass@word1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestManager_CreateBucket(t *testing.T) {
	err := manager.CreateBucket("bingo")
	if err != nil {
		t.Fatal(err)
	}
}

func TestManager_DetachBucketPolicy(t *testing.T) {
	err := manager.DetachBucketPolicy("bingo", "lirui")
	if err != nil {
		t.Fatal(err)
	}
}

func TestManager_AttachBucketPolicy(t *testing.T) {
	err := manager.AttachBucketPolicy("bingo", "lirui")
	if err != nil {
		t.Fatal(err)
	}
}
