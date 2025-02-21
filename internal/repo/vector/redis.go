package vector

import (
	"context"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/redisvector"
)

type RedisVectorStoreFactory struct {
	config *config.Config
}

func newRedisVectorStoreFactory(config *config.Config) *RedisVectorStoreFactory {
	return &RedisVectorStoreFactory{config}
}

func (r *RedisVectorStoreFactory) MakeVectorStore(ctx context.Context, kbId int64, embedder embeddings.Embedder) (vectorstores.VectorStore, error) {
	store, err := redisvector.New(ctx,
		redisvector.WithEmbedder(embedder),
		redisvector.WithIndexName(indexName(kbId), true),
		redisvector.WithConnectionURL(fmt.Sprintf("redis://%s:%s", r.config.Redis.Host, r.config.Redis.Port)))
	if err != nil {
		return nil, err
	}
	return store, nil
}
