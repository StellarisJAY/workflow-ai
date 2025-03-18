package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/bwmarrin/snowflake"
	"time"
)

type ProviderService struct {
	repo      *repo.ProviderRepo
	snowflake *snowflake.Node
}

func NewProviderService(repo *repo.ProviderRepo, snowflake *snowflake.Node) *ProviderService {
	return &ProviderService{repo: repo, snowflake: snowflake}
}

func (s *ProviderService) CreateProvider(ctx context.Context, p *model.Provider) error {
	p.Id = s.snowflake.Generate().Int64()
	p.AddTime = time.Now()
	var err error
	switch p.Code {
	case model.ProviderCodeOpenAI:
		err = validateOpenAIProvider(p)
	case model.ProviderCodeOllama:
		err = validateOllamaProvider(p)
	case model.ProviderCodeTongyi:
		err = validateTongyiProvider(p)
	default:
		return errors.New("unsupported provider")
	}
	if err != nil {
		return err
	}
	err = s.repo.InsertProvider(ctx, p)
	return err
}

func (s *ProviderService) CreateProviderModel(ctx context.Context, pm *model.ProviderModel) error {
	pm.Id = s.snowflake.Generate().Int64()
	pm.AddTime = time.Now()
	provider, err := s.repo.GetProvider(ctx, pm.ProviderId)
	if err != nil {
		return err
	}
	if provider == nil {
		return errors.New("provider not found")
	}
	if pm.Credentials == "" {
		pm.Credentials = provider.Credentials
	}
	return s.repo.InsertProviderModel(ctx, pm)
}

func (s *ProviderService) ListProviderModels(ctx context.Context, query *model.ProviderModelQuery) ([]*model.ProviderModelListDTO,
	int, error) {
	return s.repo.GetProviderModelList(ctx, query)
}

func (s *ProviderService) ListProviders(ctx context.Context) ([]*model.ProviderListDTO, error) {
	return s.repo.GetProviderList(ctx)
}

func validateOpenAIProvider(p *model.Provider) error {
	var credentials model.OpenAICredentials
	if err := json.Unmarshal([]byte(p.Credentials), &credentials); err != nil {
		return err
	}
	if credentials.ApiKey == "" {
		return errors.New("api key is empty")
	}
	if credentials.BaseUrl == "" {
		return errors.New("base url is empty")
	}
	return nil
}

func validateOllamaProvider(p *model.Provider) error {
	return nil
}

func validateTongyiProvider(p *model.Provider) error {
	var credentials model.TongyiCredentials
	if err := json.Unmarshal([]byte(p.Credentials), &credentials); err != nil {
		return err
	}
	if credentials.ApiKey == "" {
		return errors.New("api key is empty")
	}
	return nil
}
