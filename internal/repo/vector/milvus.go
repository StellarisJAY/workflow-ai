package vector

import (
	"context"
	"errors"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/milvus-io/milvus/client/v2/column"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"strconv"
)

type MilvusVectorStore struct {
	schema   *entity.Schema
	client   *milvusclient.Client
	config   config.Config
	embedder embeddings.Embedder
	kbId     int64
}

type MilvusVectorStoreFactory struct {
	config *config.Config
}

const milvusStoreDbName = "workflow_ai"
const milvusStoreCollectionName = "kb_documents"

type MilvusOption func(*MilvusVectorStore)

func WithEmbedder(embedder embeddings.Embedder) MilvusOption {
	return func(s *MilvusVectorStore) {
		s.embedder = embedder
	}
}

func newMilvusVectorStoreFactory(config *config.Config) *MilvusVectorStoreFactory {
	return &MilvusVectorStoreFactory{config: config}
}

func (m *MilvusVectorStoreFactory) MakeVectorStore(_ context.Context, kbId int64, embedder embeddings.Embedder) (Store, error) {
	store, err := newMilvusVectorStore(*m.config, kbId, WithEmbedder(embedder))
	if err != nil {
		return nil, err
	}
	return store, nil
}

func newMilvusVectorStore(config config.Config, kbId int64, options ...MilvusOption) (*MilvusVectorStore, error) {
	store := &MilvusVectorStore{config: config, kbId: kbId}
	for _, option := range options {
		option(store)
	}
	client, err := milvusclient.New(context.Background(), &milvusclient.ClientConfig{
		Address:  config.Milvus.Address,
		Username: config.Milvus.Username,
		Password: config.Milvus.Password,
	})
	if err != nil {
		return nil, err
	}
	store.client = client
	return store, nil
}

func (m *MilvusVectorStore) init(dim int64) error {

	if err := m.createDatabase(); err != nil {
		return err
	}
	if err := m.createCollection(dim); err != nil {
		return err
	}
	if err := m.createIndex(); err != nil {
		return err
	}
	return nil
}

func (m *MilvusVectorStore) createDatabase() error {
	db, _ := m.client.DescribeDatabase(context.Background(), milvusclient.NewDescribeDatabaseOption(milvusStoreDbName))
	if db == nil {
		err := m.client.CreateDatabase(context.Background(), milvusclient.NewCreateDatabaseOption(milvusStoreDbName))
		if err != nil {
			return err
		}
	}
	return m.client.UseDatabase(context.Background(), milvusclient.NewUseDatabaseOption(milvusStoreDbName))
}

func (m *MilvusVectorStore) createCollection(dim int64) error {
	c, _ := m.client.DescribeCollection(context.Background(), milvusclient.NewDescribeCollectionOption(milvusStoreCollectionName))
	if c != nil {
		return nil
	}
	// 创建collection定义
	fileIdField := entity.NewField().WithName("file_id").WithDataType(entity.FieldTypeInt64)
	chunkId := entity.NewField().WithName("chunk_id").WithDataType(entity.FieldTypeInt64).WithIsAutoID(true).WithIsPrimaryKey(true)
	kbIdFiled := entity.NewField().WithName("kb_id").WithDataType(entity.FieldTypeInt64)
	embeddingField := entity.NewField().WithName("embedding").WithDataType(entity.FieldTypeFloatVector).WithDim(dim)
	// 在chunkData字段上做全文搜索索引 倒排索引： 词->文档
	chunkDataField := entity.NewField().WithName("chunk_data").WithDataType(entity.FieldTypeVarChar).
		WithMaxLength(4096).WithEnableAnalyzer(true).WithEnableMatch(true).
		WithAnalyzerParams(map[string]any{"tokenizer": "jieba"}) // 设置中文分词器
	// 创建稀疏向量字段，用于全文搜索
	sparseField := entity.NewField().WithName("sparse_vector").WithDataType(entity.FieldTypeSparseVector)
	// 创建BM25函数，将chunk_data转换成稀疏向量
	bm25Function := entity.NewFunction().WithName("bm25").
		WithType(entity.FunctionTypeBM25).
		WithInputFields("chunk_data").
		WithOutputFields("sparse_vector")

	m.schema = entity.NewSchema().WithName(milvusStoreCollectionName).
		WithDescription("kb documents").
		WithField(fileIdField).
		WithField(kbIdFiled).
		WithField(embeddingField).
		WithField(chunkDataField).
		WithField(chunkId).
		WithField(sparseField).
		WithFunction(bm25Function)
	// 创建相似度搜索索引
	embeddingIndex := milvusclient.NewCreateIndexOption(milvusStoreCollectionName, "embedding",
		index.NewAutoIndex(entity.IP))
	// 创建BM25索引
	sparseIndex := milvusclient.NewCreateIndexOption(milvusStoreCollectionName, "sparse_vector",
		index.NewSparseInvertedIndex(entity.BM25, 0.2))
	err := m.client.CreateCollection(context.Background(),
		milvusclient.NewCreateCollectionOption(milvusStoreCollectionName, m.schema).
			WithIndexOptions(embeddingIndex, sparseIndex))
	if err != nil {
		return err
	}
	// 加载collection，否则无法搜索
	collection, err := m.client.LoadCollection(context.Background(), milvusclient.NewLoadCollectionOption(milvusStoreCollectionName))
	if err != nil {
		return err
	}
	return collection.Await(context.Background())
}

func (m *MilvusVectorStore) createIndex() error {
	return nil
}

func (m *MilvusVectorStore) SimilaritySearch(ctx context.Context, query string, n int, threshold float32) ([]*model.KbSearchReturnDocument, error) {
	queryVector, err := m.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	if err := m.client.UseDatabase(ctx, milvusclient.NewUseDatabaseOption(milvusStoreDbName)); err != nil {
		return nil, err
	}
	vectors := []entity.Vector{
		entity.FloatVector(queryVector),
	}
	idxParam := index.NewHNSWAnnParam(100)
	idxParam.WithRadius(float64(threshold))
	option := milvusclient.NewSearchOption(milvusStoreCollectionName, n, vectors).
		WithFilter(fmt.Sprintf("kb_id == %d", m.kbId)).
		WithAnnParam(idxParam).
		WithANNSField("embedding"). // 指定相似度搜索的字段
		WithOutputFields("file_id", "chunk_data", "chunk_id")
	result, err := m.client.Search(context.Background(), option)
	if err != nil {
		return nil, err
	}
	return m.convertSearchResult(result)
}

func (m *MilvusVectorStore) AddDocuments(ctx context.Context, docs []schema.Document) ([]string, error) {
	documents := make([]string, len(docs))
	for i, doc := range docs {
		documents[i] = doc.PageContent
	}
	vectors, err := m.embedder.EmbedDocuments(ctx, documents)
	if err != nil {
		return nil, err
	}
	if len(vectors) != len(docs) {
		return nil, errors.New("embedder invalid output")
	}
	if err := m.init(int64(len(vectors[0]))); err != nil {
		return nil, err
	}
	for i, doc := range docs {
		embedding := entity.FloatVector(vectors[i])
		columns := []column.Column{
			column.NewColumnInt64("file_id", []int64{doc.Metadata["fileId"].(int64)}),
			column.NewColumnInt64("kb_id", []int64{doc.Metadata["kbId"].(int64)}),
			column.NewColumnFloatVector("embedding", embedding.Dim(), [][]float32{vectors[i]}),
			column.NewColumnVarChar("chunk_data", []string{doc.PageContent}),
		}
		option := milvusclient.NewColumnBasedInsertOption(milvusStoreCollectionName, columns...)
		_, err := m.client.Insert(ctx, option)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (m *MilvusVectorStore) FulltextSearch(ctx context.Context, query string, n int) ([]*model.KbSearchReturnDocument, error) {
	if err := m.client.UseDatabase(ctx, milvusclient.NewUseDatabaseOption(milvusStoreDbName)); err != nil {
		return nil, err
	}
	indexParam := index.NewSparseAnnParam()
	indexParam.WithRadius(0.5)
	// query以Varchar类型传入，会经过数据库的bm25函数转换成稀疏向量，然后在稀疏向量索引上搜索
	option := milvusclient.NewSearchOption(milvusStoreCollectionName, n, []entity.Vector{entity.Text(query)}).
		WithFilter(fmt.Sprintf("kb_id == %d", m.kbId)).
		WithAnnParam(indexParam).
		WithANNSField("sparse_vector"). // 利用稀疏向量进行全文搜索
		WithOutputFields("file_id", "chunk_data", "chunk_id")
	result, err := m.client.Search(context.Background(), option)
	if err != nil {
		return nil, err
	}
	return m.convertSearchResult(result)
}

func (m *MilvusVectorStore) HybridSearch(ctx context.Context, query string, n int, threshold float32,
	option model.HybridSearchOption) ([]*model.KbSearchReturnDocument, error) {
	if err := m.client.UseDatabase(ctx, milvusclient.NewUseDatabaseOption(milvusStoreDbName)); err != nil {
		return nil, err
	}
	queryVector, err := m.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	filter := fmt.Sprintf("kb_id == %d", m.kbId)
	// 全文搜索配置
	sparseAnnParam := index.NewSparseAnnParam()
	sparseAnnParam.WithRadius(0.5)
	sparseRequest := milvusclient.NewAnnRequest("sparse_vector", n, entity.Text(query)).
		WithFilter(filter).WithAnnParam(sparseAnnParam).WithANNSField("sparse_vector")
	// 相似度搜索配置
	denseParams := index.NewHNSWAnnParam(100)
	denseParams.WithRadius(float64(threshold))
	denseRequest := milvusclient.NewAnnRequest("embedding", n, entity.FloatVector(queryVector)).
		WithAnnParam(denseParams).
		WithFilter(filter)
	searchOption := milvusclient.NewHybridSearchOption(milvusStoreCollectionName, n, sparseRequest, denseRequest).
		WithOutputFields("file_id", "chunk_data", "chunk_id")
	// 权重reranker，将全文搜索和相似度搜索的结果进行加权融合
	if option.WeightedRerank {
		reranker := milvusclient.NewWeightedReranker([]float64{option.SparseWeight, option.DenseWeight})
		searchOption.WithReranker(reranker)
	}
	result, err := m.client.HybridSearch(ctx, searchOption)
	if err != nil {
		return nil, err
	}
	return m.convertSearchResult(result)
}

func (m *MilvusVectorStore) ListChunks(ctx context.Context, fileId int64, paged bool, page, pageSize int) ([]*model.KbSearchReturnDocument, int, error) {
	if err := m.client.UseDatabase(ctx, milvusclient.NewUseDatabaseOption(milvusStoreDbName)); err != nil {
		return nil, 0, err
	}
	expr := fmt.Sprintf("file_id == %d", fileId)

	option := milvusclient.NewQueryOption(milvusStoreCollectionName).
		WithOutputFields("file_id", "chunk_data", "chunk_id").
		WithFilter(expr)
	total := 0
	// 分页查询，获取总数
	if paged {
		option = option.WithOffset((page - 1) * pageSize).WithLimit(pageSize)
		// 不指定outputFields，只获取总数
		totalOption := milvusclient.NewQueryOption(milvusStoreCollectionName).WithFilter(expr)
		res, err := m.client.Query(ctx, totalOption)
		if err != nil {
			return nil, 0, err
		}
		total = res.ResultCount
	}
	results, err := m.client.Query(ctx, option)
	if err != nil {
		return nil, 0, err
	}
	if !paged {
		total = results.ResultCount
	}
	chunks, err := m.convertSearchResult([]milvusclient.ResultSet{results})
	return chunks, total, err
}

func (m *MilvusVectorStore) Delete(ctx context.Context, fileId int64) error {
	if err := m.client.UseDatabase(ctx, milvusclient.NewUseDatabaseOption(milvusStoreDbName)); err != nil {
		return err
	}
	option := milvusclient.NewDeleteOption(milvusStoreCollectionName).WithExpr(fmt.Sprintf("file_id == %d", fileId))
	_, err := m.client.Delete(ctx, option)
	return err
}

func (m *MilvusVectorStore) Close() {
	_ = m.client.Close(context.Background())
}

func (m *MilvusVectorStore) convertSearchResult(results []milvusclient.ResultSet) ([]*model.KbSearchReturnDocument, error) {
	if len(results) == 0 {
		return nil, nil
	}
	rs := results[0]
	documents := make([]*model.KbSearchReturnDocument, rs.ResultCount)
	// 记录每个字段的位置
	fieldIdxMap := make(map[string]int)
	for i, field := range rs.Fields {
		fieldIdxMap[field.Name()] = i
	}
	for i := range rs.ResultCount {
		fileId, err := rs.Fields[fieldIdxMap["file_id"]].GetAsInt64(i)
		if err != nil {
			return nil, err
		}
		chunkId, err := rs.Fields[fieldIdxMap["chunk_id"]].GetAsInt64(i)
		if err != nil {
			return nil, err
		}
		chunkData, err := rs.Fields[fieldIdxMap["chunk_data"]].GetAsString(i)
		if err != nil {
			return nil, err
		}
		documents[i] = &model.KbSearchReturnDocument{
			ChunkId: strconv.FormatInt(chunkId, 10),
			FileId:  strconv.FormatInt(fileId, 10),
			Content: chunkData,
			Score:   1.0,
		}
		// 相似度搜索会返回Scores
		if len(rs.Scores) == rs.ResultCount {
			documents[i].Score = rs.Scores[i]
		}
	}
	return documents, nil
}
