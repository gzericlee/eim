package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type FileClient struct {
	*rpcxclient.XClientPool
}

func (its *FileClient) InsertFile(file *model.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "InsertFile", &rpcmodel.FileArgs{File: file}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertFile -> %w", err)
	}

	return nil
}

func (its *FileClient) DeleteFile(fileId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "DeleteFile", &rpcmodel.FileArgs{File: &model.File{FileId: fileId}}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call DeleteFile -> %w", err)
	}

	return nil
}

func (its *FileClient) GetFile(fileId int64) (*model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.FileReply{}
	err := its.Get().Call(ctx, "GetFile", &rpcmodel.FileArgs{File: &model.File{FileId: fileId}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetFile -> %w", err)
	}

	return reply.File, nil
}
