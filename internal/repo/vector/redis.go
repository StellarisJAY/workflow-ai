package vector

import (
	"context"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/redis/rueidis"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/redisvector"
)

type RedisVectorStoreFactory struct {
	config *config.Config
}

func newRedisVectorStoreFactory(config *config.Config) *RedisVectorStoreFactory {
	return &RedisVectorStoreFactory{config}
}

func (r *RedisVectorStoreFactory) MakeVectorStore(ctx context.Context, kbId int64, embedder embeddings.Embedder) (Store, error) {
	connectUrl := fmt.Sprintf("redis://%s:%s", r.config.Redis.Host, r.config.Redis.Port)
	opts, err := rueidis.ParseURL(connectUrl)
	if err != nil {
		return nil, err
	}
	cli, err := rueidis.NewClient(opts)
	if err != nil {
		return nil, err
	}
	s := &RedisVectorStore{cli: cli, kbId: kbId}
	if embedder != nil {
		options := []redisvector.Option{
			redisvector.WithIndexName(indexName(kbId), true),
			redisvector.WithConnectionURL(connectUrl),
			redisvector.WithEmbedder(embedder),
		}
		store, err := redisvector.New(ctx, options...)
		if err != nil {
			return nil, err
		}
		s.store = store
	}
	return s, nil
}

type RedisVectorStore struct {
	store *redisvector.Store
	cli   rueidis.Client
	kbId  int64
}

func (r *RedisVectorStore) SimilaritySearch(ctx context.Context, query string, n int, threshold float32) ([]*model.KbSearchReturnDocument, error) {
	// redis使用的向量距离是 1 - COSINE，然后用距离作为相似度得分升序排序。使得分最小时相似度最大。所以这里需要转换一下，使得分最大时，相似度最大.
	// 查询语句如下：
	// FT.SEARCH index
	// @content_vector:[VECTOR_RANGE $distance_threshold $vector]=>{$yield_distance_as: distance} 查询向量距离小于threshold的结果
	// SORTBY distance ASC
	// DIALECT 2
	// LIMIT 0 10
	// PARAMS 4 vector ... distance_threshold 0.5
	//
	// 向量范围查询：@field:[VECTOR_RANGE radius $vector]=>{$YIELD_DISTANCE_AS: dist_field}

	// langchaingo v0.1.13 修改了redisvector的距离算法，不需要再手动取 1-threshold
	// 参考PR：https://github.com/tmc/langchaingo/pull/1003
	docs, err := r.store.SimilaritySearch(ctx, query, n, vectorstores.WithScoreThreshold(threshold))
	if err != nil {
		return nil, err
	}
	res := make([]*model.KbSearchReturnDocument, len(docs))
	for i, doc := range docs {
		res[i] = &model.KbSearchReturnDocument{
			Content: doc.PageContent,
			Score:   1.0 - doc.Score, // 转换为相似度得分
			ChunkId: doc.Metadata["id"].(string),
			FileId:  doc.Metadata["fileId"].(string),
		}
	}
	return res, nil
}

func (r *RedisVectorStore) AddDocuments(ctx context.Context, docs []schema.Document) ([]string, error) {
	return r.store.AddDocuments(ctx, docs)
}

func (r *RedisVectorStore) FulltextSearch(ctx context.Context, query string, n int) ([]*model.KbSearchReturnDocument, error) {
	cmd := r.cli.B().FtSearch().
		Index(indexName(r.kbId)).
		Query(fmt.Sprintf(`@content:%s`, query)).
		Limit().OffsetNum(0, int64(n)).
		Dialect(2).
		Build()
	resp := r.cli.Do(ctx, cmd)
	_, docs, err := resp.AsFtSearch()
	if err != nil {
		return nil, err
	}
	return toKbSearchReturnDocuments(docs), nil
}

func (r *RedisVectorStore) ListChunks(ctx context.Context, fileId int64, paged bool, page, pageSize int) ([]*model.KbSearchReturnDocument, int, error) {
	offset, num := int64((page-1)*pageSize), int64(pageSize)
	c := r.cli.B().FtSearch().Index(indexName(r.kbId)).
		Query(fmt.Sprintf(`@fileId:[%d]`, fileId)).
		Sortby("order").
		Asc()
	var cmd rueidis.Completed
	if paged {
		cmd = c.Limit().OffsetNum(offset, num).Dialect(2).Build()
	} else {
		cmd = c.Dialect(2).Build()
	}
	total, docs, err := r.cli.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, 0, err
	}
	return toKbSearchReturnDocuments(docs), int(total), nil
}

func (r *RedisVectorStore) getChunkKeys(ctx context.Context, fileId int64) ([]string, error) {
	cmd := r.cli.B().FtSearch().Index(indexName(r.kbId)).
		Query(fmt.Sprintf(`@fileId:[%d]`, fileId)).
		Return("1").Identifier("id").
		Limit().
		OffsetNum(0, 10000). // TODO 优化获取所有分片key的方法
		Dialect(2).
		Build()
	_, resp, err := r.cli.Do(ctx, cmd).AsFtSearch()
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(resp))
	for i, doc := range resp {
		ids[i] = doc.Key
	}
	return ids, nil
}

func (r *RedisVectorStore) Delete(ctx context.Context, fileId int64) error {
	keys, err := r.getChunkKeys(ctx, fileId)
	if err != nil {
		return err
	}
	cmd := r.cli.B().Del().Key(keys...).Build()
	return r.cli.Do(ctx, cmd).Error()
}

func toKbSearchReturnDocuments(docs []rueidis.FtSearchDoc) []*model.KbSearchReturnDocument {
	res := make([]*model.KbSearchReturnDocument, len(docs))
	for i, doc := range docs {
		res[i] = &model.KbSearchReturnDocument{
			Content: doc.Doc["content"],
			Score:   float32(1 - doc.Score),
			ChunkId: doc.Doc["id"],
			FileId:  doc.Doc["fileId"],
		}
	}
	return res
}

func (r *RedisVectorStore) Close() {
	r.cli.Close()
}
