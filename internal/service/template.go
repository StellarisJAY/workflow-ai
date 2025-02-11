package service

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/bwmarrin/snowflake"
	"time"
)

type TemplateService struct {
	repo      *repo.TemplateRepo
	snowflake *snowflake.Node
}

func NewTemplateService(repo *repo.TemplateRepo, snowflake *snowflake.Node) *TemplateService {
	return &TemplateService{repo: repo, snowflake: snowflake}
}

func (t *TemplateService) Insert(ctx context.Context, template *model.Template) error {
	template.Id = t.snowflake.Generate().Int64()
	template.AddTime = time.Now()
	template.AddUser = 1
	return t.repo.Insert(ctx, template)
}

func (t *TemplateService) Get(ctx context.Context, id int64) (*model.TemplateDetailDTO, error) {
	detail, err := t.repo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, errors.New("not found")
	}
	return detail, nil
}

func (t *TemplateService) List(ctx context.Context, query *model.TemplateQuery) ([]*model.TemplateListDTO, error) {
	return t.repo.List(ctx, query)
}
