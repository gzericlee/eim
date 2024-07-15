package model

import "github.com/gzericlee/eim/internal/model"

type ListFilesArgs struct {
	Filter map[string]interface{}
	Order  []string
	Limit  int64
	Offset int64
}

type FileArgs struct {
	File *model.File
}

type FileReply struct {
	File *model.File
}

type FilesReply struct {
	Files []*model.File
	Total int64
}
