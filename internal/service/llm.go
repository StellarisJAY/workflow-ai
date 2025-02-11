package service

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/bwmarrin/snowflake"
	"time"
)

type LLMService struct {
	repo      *repo.LLMRepo
	snowflake *snowflake.Node
}

func NewLLMService(repo *repo.LLMRepo, snowflake *snowflake.Node) *LLMService {
	return &LLMService{repo, snowflake}
}

func (ls *LLMService) Create(ctx context.Context, llm *model.LLM) error {
	llm.Id = ls.snowflake.Generate().Int64()
	llm.AddTime = time.Now()
	llm.AddUser = 1
	if err := ls.repo.Insert(ctx, llm); err != nil {
		return err
	}
	return nil
}

func (ls *LLMService) Get(ctx context.Context, id int64) (*model.LLMDetailDTO, error) {
	detail, err := ls.repo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, errors.New("not found")
	}
	return detail, nil
}

func (ls *LLMService) List(ctx context.Context, query *model.LLMQuery) ([]*model.LLMListDTO, error) {
	return ls.repo.List(ctx, query)
}
