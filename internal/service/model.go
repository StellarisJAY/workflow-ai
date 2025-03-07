package service

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/bwmarrin/snowflake"
	"time"
)

type ModelService struct {
	repo      *repo.ModelRepo
	snowflake *snowflake.Node
}

func NewModelService(repo *repo.ModelRepo, snowflake *snowflake.Node) *ModelService {
	return &ModelService{repo, snowflake}
}

func (ms *ModelService) Create(ctx context.Context, m *model.Model) error {
	m.Id = ms.snowflake.Generate().Int64()
	m.AddTime = time.Now()
	m.AddUser = 1
	if err := ms.repo.Insert(ctx, m); err != nil {
		return err
	}
	return nil
}

func (ms *ModelService) Get(ctx context.Context, id int64) (*model.ModelDetailDTO, error) {
	detail, err := ms.repo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, errors.New("not found")
	}
	return detail, nil
}

func (ms *ModelService) List(ctx context.Context, query *model.ModelQuery) ([]*model.ModelListDTO, int, error) {
	return ms.repo.List(ctx, query)
}
