package mongodb

import "github.com/gzericlee/eim/internal/model"

func (its *Repository) InsertFile(file *model.File) error {
	//TODO implement me
	panic("implement me")
}

func (its *Repository) DeleteFile(fileId string) error {
	//TODO implement me
	panic("implement me")
}

func (its *Repository) GetFile(fileId string) (*model.File, error) {
	//TODO implement me
	panic("implement me")
}

func (its *Repository) ListFiles(filter map[string]interface{}, order []string, limit, offset int64) ([]*model.File, int64, error) {
	//TODO implement me
	panic("implement me")
}
