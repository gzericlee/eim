package service

import (
	"context"
	"fmt"

	"github.com/gzericlee/eim/internal/database"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type FileThumbService struct {
	database database.IDatabase
}

func NewFileThumbService(database database.IDatabase) *FileThumbService {
	return &FileThumbService{
		database: database,
	}
}

func (its *FileThumbService) InsertFileThumb(ctx context.Context, args *rpcmodel.FileThumbArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.InsertFileThumb(args.FileThumb)
	if err != nil {
		return fmt.Errorf("insert file thumb -> %w", err)
	}
	return nil
}

func (its *FileThumbService) DeleteFileThumb(ctx context.Context, args *rpcmodel.FileThumbArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.DeleteFileThumb(args.FileThumb.ThumbId)
	if err != nil {
		return fmt.Errorf("delete file thumb -> %w", err)
	}
	return nil
}

func (its *FileThumbService) GetFileThumb(ctx context.Context, args *rpcmodel.FileThumbArgs, reply *rpcmodel.FileThumbReply) error {
	fileThumb, err := its.database.GetFileThumb(args.FileThumb.FileId, args.FileThumb.ThumbSpec)
	if err != nil {
		return fmt.Errorf("get file thumb -> %w", err)
	}
	reply.FileThumb = fileThumb

	return nil
}
