package vector

import (
	"context"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores"
)

type Factory interface {
	MakeVectorStore(ctx context.Context, kbId int64, embedder embeddings.Embedder) (vectorstores.VectorStore, error)
}

func indexName(kbId int64) string {
	return fmt.Sprintf("workflowai/kb/%d", kbId)
}

func MakeFactory(cfg config.Config) Factory {
	switch cfg.VectorStoreType {
	case "redis":
		return newRedisVectorStoreFactory(&cfg)
	default:
		panic("unsupported vector store type: " + cfg.VectorStoreType)
	}
}
