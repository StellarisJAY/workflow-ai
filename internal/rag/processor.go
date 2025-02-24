package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/repo/fs"
	"github.com/StellrisJAY/workflow-ai/internal/repo/vector"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"log"
	"time"
)

type DocumentProcessor struct {
	kbRepo             *repo.KnowledgeBaseRepo
	fs                 fs.FileStore
	workerCancel       []context.CancelFunc
	taskChan           chan int64 // 任务队列 kbProcessTaskId
	workers            int
	llmRepo            *repo.LLMRepo
	vectorStoreFactory vector.Factory
}

func NewDocumentProcessor(workers int, kbRepo *repo.KnowledgeBaseRepo, fs fs.FileStore, llmRepo *repo.LLMRepo,
	factory vector.Factory) *DocumentProcessor {
	proc := &DocumentProcessor{
		kbRepo:             kbRepo,
		workerCancel:       make([]context.CancelFunc, workers),
		taskChan:           make(chan int64, 100),
		workers:            workers,
		fs:                 fs,
		llmRepo:            llmRepo,
		vectorStoreFactory: factory,
	}
	for i := 0; i < workers; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		proc.workerCancel[i] = cancel
		go proc.worker(ctx)
	}
	return proc
}

func (d *DocumentProcessor) SubmitTask(taskId int64) {
	d.taskChan <- taskId
}

func (d *DocumentProcessor) SimilaritySearch(ctx context.Context, kbId int64, input string, scoreThreshold float32,
	n int) ([]*model.KbSearchReturnDocument, error) {
	kb, err := d.kbRepo.Detail(ctx, kbId)
	if err != nil {
		return nil, err
	}
	llm, err := d.llmRepo.GetDetail(ctx, kb.EmbeddingModel)
	if err != nil {
		return nil, err
	}
	embeddingModel, err := makeEmbeddingModel(llm)
	if err != nil {
		return nil, err
	}
	embedder, err := embeddings.NewEmbedder(embeddingModel, embeddings.WithBatchSize(20))
	if err != nil {
		return nil, err
	}
	vectorStore, err := d.vectorStoreFactory.MakeVectorStore(ctx, kbId, embedder)
	if err != nil {
		return nil, err
	}
	defer vectorStore.Close()
	documents, err := vectorStore.SimilaritySearch(ctx, input, n, scoreThreshold)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (d *DocumentProcessor) FulltextSearch(ctx context.Context, kbId int64, input string, n int) ([]*model.KbSearchReturnDocument, error) {
	vectorStore, err := d.vectorStoreFactory.MakeVectorStore(ctx, kbId, nil)
	if err != nil {
		return nil, err
	}
	defer vectorStore.Close()
	documents, err := vectorStore.FulltextSearch(ctx, input, n)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (d *DocumentProcessor) ListChunks(ctx context.Context, kbId int64, fileId int64, page int, pageSize int) ([]*model.KbSearchReturnDocument, int, error) {
	vectorStore, err := d.vectorStoreFactory.MakeVectorStore(ctx, kbId, nil)
	if err != nil {
		return nil, 0, err
	}
	defer vectorStore.Close()
	return vectorStore.ListChunks(ctx, fileId, page, pageSize)
}

// splitDocument 拆分文档
func (d *DocumentProcessor) splitDocument(ctx context.Context, file *model.KnowledgeBaseFile, options *model.KbFileProcessOptions) ([]schema.Document, error) {
	var separators []string
	if err := json.Unmarshal([]byte(options.Separators), &separators); err != nil {
		return nil, err
	}
	data, err := d.fs.Download(ctx, file.Url)
	if err != nil {
		return nil, err
	}
	var metadata map[string]any
	if err := json.Unmarshal([]byte(file.Metadata), &metadata); err != nil {
		return nil, err
	}
	// TODO 不同的文件类型使用不同的 loader
	text := documentloaders.NewText(bytes.NewReader(data))
	splitter := textsplitter.NewMarkdownTextSplitter(textsplitter.WithSeparators(separators),
		textsplitter.WithChunkSize(options.ChunkSize))
	chunks, err := text.LoadAndSplit(ctx, splitter)
	if err != nil {
		return nil, err
	}
	for i, chunk := range chunks {
		chunk.Metadata["fileId"] = file.Id
		chunk.Metadata["order"] = i
	}
	return chunks, nil
}

// embedDocument 嵌入文档
func (d *DocumentProcessor) embedDocument(ctx context.Context, kb *model.KnowledgeBaseDetailDTO, chunks []schema.Document) ([]string, error) {
	llm, err := d.llmRepo.GetDetail(ctx, kb.EmbeddingModel)
	if err != nil {
		return nil, err
	}
	embeddingModel, err := makeEmbeddingModel(llm)
	if err != nil {
		return nil, err
	}
	embedder, err := embeddings.NewEmbedder(embeddingModel, embeddings.WithBatchSize(20))
	if err != nil {
		return nil, err
	}
	store, err := d.vectorStoreFactory.MakeVectorStore(ctx, kb.Id, embedder)
	if err != nil {
		return nil, err
	}
	defer store.Close()
	chunkIds, err := store.AddDocuments(ctx, chunks)
	return chunkIds, err
}

func (d *DocumentProcessor) worker(ctx context.Context) {
	for {
		select {
		case taskId := <-d.taskChan:
			d.handleTask(ctx, taskId)
		case <-ctx.Done():
			return
		}
	}
}

func (d *DocumentProcessor) handleTask(ctx context.Context, taskId int64) {
	task, err := d.kbRepo.GetFileProcessTask(ctx, taskId)
	if err != nil {
		log.Println("handleTask err:", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			task.Status = model.KbFileProcessStatusFailed
			task.Error = err.Error()
			task.CompleteTime = time.Now()
			if err := d.kbRepo.UpdateFileProcessTask(ctx, task); err != nil {
				log.Println("handleTask err:", err)
			}
			if err := d.kbRepo.UpdateFileStatus(ctx, task.FileId, model.KbFileFailed); err != nil {
				log.Println("handleTask err:", err)
			}
		}
	}()
	kb, err := d.kbRepo.Detail(ctx, task.KbId)
	if err != nil {
		panic(err)
	}
	file, err := d.kbRepo.GetFileDetail(ctx, task.FileId)
	if err != nil {
		panic(err)
	}
	options, err := d.kbRepo.GetFileProcessOptions(ctx, task.FileId)
	if err != nil {
		panic(err)
	}
	// 更新到splitting状态
	task.Status = model.KbFileProcessStatusSplitting
	if err := d.kbRepo.UpdateFileProcessTask(ctx, task); err != nil {
		panic(err)
	}
	if err := d.kbRepo.UpdateFileStatus(ctx, task.FileId, model.KbFileProcessing); err != nil {
		panic(err)
	}
	// step1 加载并拆分文档
	chunks, err := d.splitDocument(ctx, file, options)
	if err != nil {
		panic(err)
	}
	// 更新到embedding状态
	task.Status = model.KbFileProcessStatusEmbedding
	if err := d.kbRepo.UpdateFileProcessTask(ctx, task); err != nil {
		log.Println("handleTask err:", err)
	}
	// step2 嵌入文档
	_, err = d.embedDocument(ctx, kb, chunks)
	if err != nil {
		panic(err)
	}
	task.Status = model.KbFileProcessStatusCompleted
	task.CompleteTime = time.Now()
	if err := d.kbRepo.UpdateFileProcessTask(ctx, task); err != nil {
		log.Println("handleTask err:", err)
	}
	if err := d.kbRepo.UpdateFileStatus(ctx, file.Id, model.KbFileProcessed); err != nil {
		log.Println("handleTask err:", err)
	}
}

func makeEmbeddingModel(llm *model.LLMDetailDTO) (embeddings.EmbedderClient, error) {
	var embeddingModel embeddings.EmbedderClient
	var err error
	switch llm.ApiType {
	case string(model.ApiTypeOpenAI):
		embeddingModel, err = openai.New(openai.WithEmbeddingModel(llm.Code),
			openai.WithBaseURL(llm.BaseUrl),
			openai.WithToken(llm.ApiKey))
	case string(model.ApiTypeOllama):
		embeddingModel, err = ollama.New(ollama.WithModel(llm.Code), ollama.WithServerURL(llm.BaseUrl))
	default:
		err = fmt.Errorf("unknown api type: %s", llm.ApiType)
	}
	return embeddingModel, err
}
