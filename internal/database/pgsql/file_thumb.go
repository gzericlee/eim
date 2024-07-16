package pgsql

import (
	"fmt"

	"github.com/gzericlee/eim/internal/model"
)

func (its *Repository) InsertFileThumb(thumb *model.FileThumb) error {
	_, err := its.db.Insert(thumb)
	if err != nil {
		return fmt.Errorf("insert file thumb -> %w", err)
	}
	return nil
}

func (its *Repository) DeleteFileThumb(thumbId int64) error {
	_, err := its.db.Where("thumb_id = ?", thumbId).Delete(&model.FileThumb{})
	if err != nil {
		return fmt.Errorf("delete file thumb -> %w", err)
	}
	return nil
}

func (its *Repository) GetFileThumb(fileId int64, thumbSpec string) (*model.FileThumb, error) {
	thumb := &model.FileThumb{}
	_, err := its.db.Where("file_id = ? AND thumb_spec = ?", fileId, thumbSpec).Get(thumb)
	if err != nil {
		return nil, fmt.Errorf("select file thumb -> %w", err)
	}
	return thumb, nil
}
