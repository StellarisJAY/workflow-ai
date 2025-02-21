package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/rag"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/repo/fs"
	"github.com/bwmarrin/snowflake"
	"io"
	"log"
	"strconv"
	"time"
)

type KnowledgeBaseService struct {
	kbRepo    *repo.KnowledgeBaseRepo
	snowflake *snowflake.Node
	tm        *repo.TransactionManager
	fs        fs.FileStore
	processor *rag.DocumentProcessor
}

func NewKnowledgeBaseService(kbRepo *repo.KnowledgeBaseRepo, snowflake *snowflake.Node,
	tm *repo.TransactionManager, fs fs.FileStore, processor *rag.DocumentProcessor) *KnowledgeBaseService {
	return &KnowledgeBaseService{kbRepo: kbRepo, snowflake: snowflake, tm: tm, fs: fs, processor: processor}
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
		path := fmt.Sprintf("%d/%s", file.KbId, file.Name)
		file.Url = path
		file.Id = k.snowflake.Generate().Int64()
		file.AddTime = time.Now()
		file.AddUser = 1
		file.Status = model.KbFileStatusUnavailable
		file.Metadata = "{}"
		// 添加文件
		if err := k.kbRepo.InsertFile(ctx, file); err != nil {
			return err
		}
		// 创建默认的解析选项
		options := model.DefaultKbFileProcessOptions()
		options.FileId = file.Id
		if err := k.kbRepo.InsertFileProcessOptions(ctx, &options); err != nil {
			return err
		}

		// 上传文件数据
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

func (k *KnowledgeBaseService) Delete(ctx context.Context, fileId int64) error {
	return k.tm.Tx(ctx, func(ctx context.Context) error {
		file, err := k.kbRepo.GetFileDetail(ctx, fileId)
		if err != nil {
			return err
		}
		// 删除文件
		if err := k.kbRepo.DeleteFile(ctx, fileId); err != nil {
			return err
		}
		// 删除文件解析选项
		if err := k.kbRepo.DeleteFileProcessOptions(ctx, fileId); err != nil {
			return err
		}
		// 删除文件数据
		if err := k.fs.Delete(ctx, file.Url); err != nil {
			return err
		}
		return nil
	})
}

func (k *KnowledgeBaseService) DownloadFile(ctx context.Context, fileId int64) ([]byte, string, error) {
	detail, err := k.kbRepo.GetFileDetail(ctx, fileId)
	if err != nil {
		return nil, "", err
	}
	data, err := k.fs.Download(ctx, detail.Url)
	if err != nil {
		return nil, "", err
	}
	return data, detail.Name, nil
}

func (k *KnowledgeBaseService) GetFileProcessOptions(ctx context.Context, fileId int64) (*model.KbFileProcessOptions, error) {
	return k.kbRepo.GetFileProcessOptions(ctx, fileId)
}

func (k *KnowledgeBaseService) UpdateFileProcessOptions(ctx context.Context, dto *model.KbFileProcessOptionsUpdateDTO) error {
	separators, err := json.Marshal(dto.Separators)
	if err != nil {
		return err
	}
	options := model.KbFileProcessOptions{
		FileId:     dto.FileId,
		ChunkSize:  dto.ChunkSize,
		Separators: string(separators),
	}
	return k.kbRepo.UpdateFileProcessOptions(ctx, &options)
}

func (k *KnowledgeBaseService) ProcessFile(ctx context.Context, fileId int64) error {
	file, err := k.kbRepo.GetFileDetail(ctx, fileId)
	if err != nil {
		return err
	}
	task := model.KbFileProcessTask{
		Id:           k.snowflake.Generate().Int64(),
		KbId:         file.KbId,
		FileId:       fileId,
		Status:       model.KbFileProcessStatusQueued,
		AddTime:      time.Now(),
		CompleteTime: time.Now(),
	}
	if err := k.kbRepo.InsertFileProcessTask(ctx, &task); err != nil {
		return err
	}
	k.processor.SubmitTask(task.Id)
	return nil
}

func (k *KnowledgeBaseService) SimilaritySearch(ctx context.Context, request *model.KbSearchRequest) (*model.KbSearchResult, error) {
	documents, err := k.processor.SimilaritySearch(ctx, request.KbId, request.Input, request.ScoreThreshold, request.Count)
	if err != nil {
		return nil, err
	}
	ids := make(map[string]struct{})
	for _, document := range documents {
		ids[document.FileId] = struct{}{}
	}
	idList := make([]int64, 0, len(ids))
	for id := range ids {
		i, _ := strconv.ParseInt(id, 10, 64)
		idList = append(idList, i)
	}
	files, err := k.kbRepo.GetFilesInIdList(ctx, idList)
	if err != nil {
		return nil, err
	}
	result := &model.KbSearchResult{
		Documents: documents,
		Files:     files,
	}
	return result, nil
}
