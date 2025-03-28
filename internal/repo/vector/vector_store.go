package vector

import (
	"context"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
)

type Factory interface {
	MakeVectorStore(ctx context.Context, kbId int64, embedder embeddings.Embedder) (Store, error)
}

func indexName(kbId int64) string {
	return fmt.Sprintf("workflowai/kb/%d", kbId)
}

func MakeFactory(cfg config.Config) Factory {
	switch cfg.VectorStoreType {
	case "redis":
		return newRedisVectorStoreFactory(&cfg)
	case "milvus":
		return newMilvusVectorStoreFactory(&cfg)

	default:
		panic("unsupported vector store type: " + cfg.VectorStoreType)
	}
}

type Store interface {
	// SimilaritySearch 语义搜索,搜索相似度大于threshold的前n个结果
	SimilaritySearch(ctx context.Context, query string, n int, threshold float32) ([]*model.KbSearchReturnDocument, error)
	AddDocuments(ctx context.Context, docs []schema.Document) ([]string, error)
	// FulltextSearch 全文搜索
	FulltextSearch(ctx context.Context, query string, n int) ([]*model.KbSearchReturnDocument, error)
	// HybridSearch 基于权重reranker的混合搜索
	HybridSearch(ctx context.Context, query string, n int, threshold float32, option model.HybridSearchOption) ([]*model.KbSearchReturnDocument, error)
	ListChunks(ctx context.Context, fileId int64, paged bool, page, pageSize int) ([]*model.KbSearchReturnDocument, int, error)
	Delete(ctx context.Context, fileId int64) error
	Close()
}
