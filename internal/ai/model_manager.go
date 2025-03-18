package ai

import (
	"encoding/json"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const dashscopeBaseUrl = "https://dashscope.aliyuncs.com/compatible-mode/v1"

func MakeModelInterface(detail *model.ProviderModelDetail, outputFormat string) (llms.Model, error) {
	var m llms.Model
	var err error
	switch detail.ProviderCode {
	case model.ProviderCodeOpenAI:
		var credentials model.OpenAICredentials
		if err := json.Unmarshal([]byte(detail.ProviderCredentials), &credentials); err != nil {
			return nil, err
		}
		m, err = openai.New(openai.WithModel(detail.ModelName),
			openai.WithBaseURL(credentials.BaseUrl),
			openai.WithToken(credentials.ApiKey))
	case model.ProviderCodeTongyi:
		var credentials model.TongyiCredentials
		if err := json.Unmarshal([]byte(detail.ProviderCredentials), &credentials); err != nil {
			return nil, err
		}
		m, err = openai.New(openai.WithModel(detail.ModelName),
			openai.WithToken(credentials.ApiKey),
			openai.WithBaseURL(dashscopeBaseUrl))
	default:
		return nil, fmt.Errorf("unsupported provider code: %s", detail.ProviderCode)
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func MakeEmbeddingModel(detail *model.ProviderModelDetail) (embeddings.EmbedderClient, error) {
	if detail.ModelType != model.ProviderModelTypeEmbedding {
		return nil, fmt.Errorf("model type not support embedding")
	}
	var embeddingModel embeddings.EmbedderClient
	var err error
	switch detail.ProviderCode {
	case model.ProviderCodeOpenAI:
		var credentials model.OpenAICredentials
		if err := json.Unmarshal([]byte(detail.ProviderCredentials), &credentials); err != nil {
			return nil, err
		}
		embeddingModel, err = openai.New(openai.WithEmbeddingModel(detail.ModelName),
			openai.WithBaseURL(credentials.BaseUrl),
			openai.WithToken(credentials.ApiKey))
	case model.ProviderCodeTongyi:
		var credentials model.TongyiCredentials
		if err := json.Unmarshal([]byte(detail.ProviderCredentials), &credentials); err != nil {
			return nil, err
		}
		embeddingModel, err = openai.New(openai.WithEmbeddingModel(detail.ModelName),
			openai.WithBaseURL(dashscopeBaseUrl),
			openai.WithToken(credentials.ApiKey))
	default:
		return nil, fmt.Errorf("unsupported provider code: %s", detail.ProviderCode)
	}
	if err != nil {
		return nil, err
	}
	return embeddingModel, nil
}
