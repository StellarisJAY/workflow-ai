package ai

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/ai/gte"
	"github.com/StellrisJAY/workflow-ai/internal/model"
)

type Reranker interface {
	Rerank(ctx context.Context, query string, documents []*model.KbSearchReturnDocument) ([]*model.KbSearchReturnDocument, error)
}

const GteReranker = "gte-rerank"

func MakeReranker(modelName string, providerType model.ProviderCode, credentialString string) (Reranker, error) {
	switch providerType {
	case model.ProviderCodeTongyi:
		var credentials model.TongyiCredentials
		if err := json.Unmarshal([]byte(credentialString), &credentials); err != nil {
			return nil, err
		}
		if modelName == GteReranker {
			return gte.NewReranker(credentials.ApiKey), nil
		}
		return nil, errors.New("暂不支持的排序模型")
	default:
		return nil, errors.New("暂不支持该类型的排序模型")
	}
}
