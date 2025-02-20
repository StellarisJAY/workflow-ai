package service

import (
	"context"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/repo/fs"
	"github.com/bwmarrin/snowflake"
	"io"
	"log"
	"time"
)

type KnowledgeBaseService struct {
	kbRepo    *repo.KnowledgeBaseRepo
	snowflake *snowflake.Node
	tm        *repo.TransactionManager
	fs        fs.FileStore
}

func NewKnowledgeBaseService(kbRepo *repo.KnowledgeBaseRepo, snowflake *snowflake.Node,
	tm *repo.TransactionManager, fs fs.FileStore) *KnowledgeBaseService {
	return &KnowledgeBaseService{kbRepo: kbRepo, snowflake: snowflake, tm: tm, fs: fs}
}

func (k *KnowledgeBaseService) Create(ctx context.Context, kb *model.KnowledgeBase) error {
	kb.Id = k.snowflake.Generate().Int64()
	kb.AddTime = time.Now()
	kb.AddUser = 1
	kb.Enable = true
	return k.kbRepo.Insert(ctx, kb)
}

func (k *KnowledgeBaseService) Update(ctx context.Context, kb *model.KnowledgeBase) error {
	return k.kbRepo.Update(ctx, kb)
}

func (k *KnowledgeBaseService) List(ctx context.Context, query *model.KnowledgeBaseQuery) ([]*model.KnowledgeBaseListDTO, int, error) {
	result, total, err := k.kbRepo.List(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	ids := make([]int64, len(result))
	for i, item := range result {
		ids[i] = item.Id
	}
	fileStats, err := k.kbRepo.ListFileCount(ctx, ids)
	if err != nil {
		return nil, 0, err
	}

	m := make(map[int64]*model.KbFileCountSize)
	for _, stat := range fileStats {
		m[stat.KbId] = stat
	}

	for _, kb := range result {
		stat := m[kb.Id]
		if stat != nil {
			kb.Size = stat.Size
			kb.DocumentCount = int(stat.Count)
		}
	}
	return result, total, nil
}

func (k *KnowledgeBaseService) Detail(ctx context.Context, id int64) (*model.KnowledgeBaseDetailDTO, error) {
	detail, err := k.kbRepo.Detail(ctx, id)
	if err != nil {
		return nil, err
	}
	documentCount, size, err := k.kbRepo.CountFile(ctx, id)
	if err != nil {
		return nil, err
	}
	detail.DocumentCount = int(documentCount)
	detail.Size = size
	return detail, nil
}

func (k *KnowledgeBaseService) UploadFile(ctx context.Context, file *model.KnowledgeBaseFile, reader io.Reader) error {
	return k.tm.Tx(ctx, func(ctx context.Context) error {
		file.Id = k.snowflake.Generate().Int64()
		file.AddTime = time.Now()
		file.AddUser = 1
		file.Status = model.KbFileStatusUnavailable
		if err := k.kbRepo.InsertFile(ctx, file); err != nil {
			return err
		}
		path := fmt.Sprintf("%d/%s", file.KbId, file.Name)
		file.Url = path
		if err := k.fs.Upload(ctx, path, reader); err != nil {
			return err
		}
		file.Status = model.KbFileUploaded
		if err := k.kbRepo.UpdateFile(ctx, file); err != nil {
			log.Println("update file status after upload failed", file.Id, err)
		}
		return nil
	})
}

func (k *KnowledgeBaseService) ListFile(ctx context.Context, kbId int64, query *model.KbFileQuery) ([]*model.KbFileListDTO, int, error) {
	files, total, err := k.kbRepo.ListFile(ctx, kbId, query)
	if err != nil {
		return nil, 0, err
	}
	for _, file := range files {
		file.StatusName = file.Status.String()
	}
	return files, total, nil
}
