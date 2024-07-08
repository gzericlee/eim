package minio

import (
	"os"
	"testing"
)

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

func TestManager_UploadObject(t *testing.T) {
	file, err := os.Open("test.txt")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = manager.UploadObject("bingo", "test.txt", file)
	if err != nil {
		t.Fatal(err)
		return
	}
	url, err := manager.ShareObject("bingo", "test.txt", 1)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(url)
}
