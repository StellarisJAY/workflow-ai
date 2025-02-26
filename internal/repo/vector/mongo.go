package vector

import (
	"context"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/mongovector"
)

type MongoVectorStoreFactory struct {
	conf *config.Config
}

func NewMongoVectorStoreFactory(conf *config.Config) *MongoVectorStoreFactory {
	return &MongoVectorStoreFactory{conf}
}

func (m *MongoVectorStoreFactory) MakeVectorStore(ctx context.Context, kbId int64, embedder embeddings.Embedder) (Store, error) {
	//TODO implement me
	panic("implement me")
}

type MongoVector struct {
	kbId  int64
	store *mongovector.Store
}

const mongoDbName = "workflowai"

func NewMongoVector(kbId int64, embedder embeddings.Embedder) *MongoVector {
	return nil
}

func (m *MongoVector) SimilaritySearch(ctx context.Context, query string, n int, threshold float32) ([]*model.KbSearchReturnDocument, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoVector) AddDocuments(ctx context.Context, docs []schema.Document) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoVector) FulltextSearch(ctx context.Context, query string, n int) ([]*model.KbSearchReturnDocument, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoVector) ListChunks(ctx context.Context, fileId int64, page, pageSize int) ([]*model.KbSearchReturnDocument, int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoVector) Close() {
	//TODO implement me
	panic("implement me")
}
