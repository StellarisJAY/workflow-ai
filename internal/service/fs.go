package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/repo/fs"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"io"
	"time"
)

type FileService struct {
	fileRepo  *repo.FileRepo
	store     fs.FileStore
	tm        *repo.TransactionManager
	snowflake *snowflake.Node
}

func NewFileService(fileRepo *repo.FileRepo, store fs.FileStore, tm *repo.TransactionManager, snowflake *snowflake.Node) *FileService {
	return &FileService{fileRepo: fileRepo, store: store, tm: tm, snowflake: snowflake}
}

func (f *FileService) Upload(ctx context.Context, file *model.File, data io.Reader) error {
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(dataBytes)
	hash := md5.Sum(dataBytes)
	file.Md5 = hex.EncodeToString(hash[:])
	return f.tm.Tx(ctx, func(ctx context.Context) error {
		existingFile, _ := f.fileRepo.GetFileByMd5(ctx, file.Md5)
		if f != nil {
			file.Id = existingFile.Id
			return nil
		}
		file.Id = f.snowflake.Generate().Int64()
		file.AddTime = time.Now()
		file.AddUser = 1
		file.Url = "fs/" + uuid.NewString()
		err := f.fileRepo.Create(ctx, file)
		if err != nil {
			return err
		}
		err = f.store.Upload(ctx, file.Url, buffer)
		return err
	})
}

func (f *FileService) BatchUpload(ctx context.Context, files []*model.File, data []io.Reader) error {
	return f.tm.Tx(ctx, func(ctx context.Context) error {
		needUpload := make(map[string]io.Reader)
		for i, file := range files {
			dataBytes, err := io.ReadAll(data[i])
			if err != nil {
				return err
			}
			buffer := bytes.NewBuffer(dataBytes)
			hash := md5.Sum(dataBytes)
			file.Md5 = hex.EncodeToString(hash[:])
			existingFile, _ := f.fileRepo.GetFileByMd5(ctx, file.Md5)
			if existingFile != nil {
				file.Id = existingFile.Id
				continue
			}
			file.Id = f.snowflake.Generate().Int64()
			file.AddTime = time.Now()
			file.AddUser = 1
			file.Url = "fs/" + uuid.NewString()
			if err := f.fileRepo.Create(ctx, file); err != nil {
				return err
			}
			needUpload[file.Url] = buffer
		}
		for url, reader := range needUpload {
			if err := f.store.Upload(ctx, url, reader); err != nil {
				return err
			}
		}
		return nil
	})
}

func (f *FileService) Download(ctx context.Context, id int64) ([]byte, string, error) {
	file, err := f.fileRepo.Get(ctx, id)
	if err != nil {
		return nil, "", err
	}
	data, err := f.store.Download(ctx, file.Url)
	if err != nil {
		return nil, "", err
	}
	return data, file.Name, nil
}
