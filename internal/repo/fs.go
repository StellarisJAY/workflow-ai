package repo

import (
	"context"
	"github.com/StellrisJAY/workflow-ai/internal/model"
)

type FileRepo struct {
	*Repository
}

func NewFileRepo(repo *Repository) *FileRepo {
	return &FileRepo{repo}
}

func (fr *FileRepo) Create(ctx context.Context, file *model.File) error {
	return fr.DB(ctx).Create(file).Error
}

func (fr *FileRepo) Get(ctx context.Context, id int64) (*model.File, error) {
	var file *model.File
	err := fr.DB(ctx).Table(model.File{}.TableName()).Where("id = ?", id).Scan(&file).Error
	return file, err
}

func (fr *FileRepo) GetFileByMd5(ctx context.Context, md5 string) (*model.File, error) {
	var file *model.File
	err := fr.DB(ctx).Table(model.File{}.TableName()).Where("md5 = ?", md5).Scan(&file).Error
	return file, err
}
