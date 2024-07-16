package pgsql

import (
	"fmt"

	"github.com/gzericlee/eim/internal/model"
)

func (its *Repository) InsertFile(file *model.File) error {
	_, err := its.db.Insert(file)
	if err != nil {
		return fmt.Errorf("insert file -> %w", err)
	}
	return nil
}

func (its *Repository) DeleteFile(fileId int64) error {
	_, err := its.db.Where("file_id = ?", fileId).Delete(&model.File{})
	if err != nil {
		return fmt.Errorf("delete file -> %w", err)
	}
	return nil
}

func (its *Repository) GetFile(fileId int64) (*model.File, error) {
	file := &model.File{}
	_, err := its.db.Where("file_id = ?", fileId).Get(file)
	if err != nil {
		return nil, fmt.Errorf("select file -> %w", err)
	}
	return file, nil
}

func (its *Repository) ListFiles(filter map[string]interface{}, order []string, limit, offset int64) ([]*model.File, int64, error) {
	var files []*model.File
	query := its.db.Where("")

	for key, value := range filter {
		query = query.Where(key+" = ?", value)
	}

	for _, by := range order {
		query = query.OrderBy(by)
	}

	total, err := query.Limit(int(limit), int(offset)).FindAndCount(&files)
	if err != nil {
		return nil, 0, fmt.Errorf("select files -> %w", err)
	}

	return files, total, nil
}
