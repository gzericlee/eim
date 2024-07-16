package client

import (
	"context"
	"fmt"
	"time"

	rpcxclient "github.com/smallnest/rpcx/client"

	"github.com/gzericlee/eim/internal/model"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type FileThumbClient struct {
	*rpcxclient.XClientPool
}

func (its *FileThumbClient) InsertFileThumb(thumb *model.FileThumb) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "InsertFileThumb", &rpcmodel.FileThumbArgs{FileThumb: thumb}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call InsertFileThumb -> %w", err)
	}

	return nil
}

func (its *FileThumbClient) DeleteFileThumb(thumbId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := its.Get().Call(ctx, "DeleteFileThumb", &rpcmodel.FileThumbArgs{FileThumb: &model.FileThumb{ThumbId: thumbId}}, &rpcmodel.EmptyReply{})
	if err != nil {
		return fmt.Errorf("call DeleteFileThumb -> %w", err)
	}

	return nil
}

func (its *FileThumbClient) GetFileThumb(fileId int64, spec string) (*model.FileThumb, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	reply := &rpcmodel.FileThumbReply{}
	err := its.Get().Call(ctx, "GetFileThumb", &rpcmodel.FileThumbArgs{FileThumb: &model.FileThumb{FileId: fileId, ThumbSpec: spec}}, reply)
	if err != nil {
		return nil, fmt.Errorf("call GetFileThumb -> %w", err)
	}

	return reply.FileThumb, nil
}
