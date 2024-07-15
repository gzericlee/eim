package service

import (
	"context"
	"fmt"

	"github.com/gzericlee/eim/internal/database"
	rpcmodel "github.com/gzericlee/eim/internal/storage/rpc/model"
)

type FileService struct {
	database database.IDatabase
}

func NewFileService(database database.IDatabase) *FileService {
	return &FileService{
		database: database,
	}
}

func (its *FileService) InsertFile(ctx context.Context, args *rpcmodel.FileArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.InsertFile(args.File)
	if err != nil {
		return fmt.Errorf("insert file -> %w", err)
	}
	return nil
}

func (its *FileService) GetFile(ctx context.Context, args *rpcmodel.FileArgs, reply *rpcmodel.FileReply) error {
	file, err := its.database.GetFile(args.File.FileId)
	if err != nil {
		return fmt.Errorf("get file -> %w", err)
	}
	reply.File = file
	return nil
}

func (its *FileService) DeleteFile(ctx context.Context, args *rpcmodel.FileArgs, reply *rpcmodel.EmptyReply) error {
	err := its.database.DeleteFile(args.File.FileId)
	if err != nil {
		return fmt.Errorf("delete file -> %w", err)
	}
	return nil
}

func (its *FileService) ListFiles(ctx context.Context, args *rpcmodel.ListFilesArgs, reply *rpcmodel.FilesReply) error {
	files, total, err := its.database.ListFiles(args.Filter, args.Order, args.Limit, args.Offset)
	if err != nil {
		return fmt.Errorf("list files -> %w", err)
	}
	reply.Files = files
	reply.Total = total
	return nil
}
