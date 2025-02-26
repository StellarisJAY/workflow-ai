package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/bwmarrin/snowflake"
	"slices"
	"time"
)

type TemplateService struct {
	repo      *repo.TemplateRepo
	snowflake *snowflake.Node
}

func NewTemplateService(repo *repo.TemplateRepo, snowflake *snowflake.Node) *TemplateService {
	return &TemplateService{repo: repo, snowflake: snowflake}
}

func (t *TemplateService) Insert(ctx context.Context, template *model.Template) (int64, error) {
	template.Id = t.snowflake.Generate().Int64()
	template.AddTime = time.Now()
	template.AddUser = 1
	if err := t.repo.Insert(ctx, template); err != nil {
		return 0, err
	}
	return template.Id, nil
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

func (t *TemplateService) List(ctx context.Context, query *model.TemplateQuery) ([]*model.TemplateListDTO, int, error) {
	return t.repo.List(ctx, query)
}

func (t *TemplateService) Delete(ctx context.Context, id int64) error {
	return t.repo.Delete(ctx, id)
}

func (t *TemplateService) Update(ctx context.Context, template *model.Template) error {
	return t.repo.Update(ctx, template)
}

func (t *TemplateService) GetStartInputVariables(ctx context.Context, id int64) ([]*model.Variable, error) {
	detail, err := t.repo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, errors.New("not found")
	}
	var definition model.WorkflowDefinition
	_ = json.Unmarshal([]byte(detail.Data), &definition)
	idx := slices.IndexFunc(definition.Nodes, func(node *model.Node) bool {
		return node.Type == string(model.NodeTypeStart)
	})
	if idx == -1 {
		return nil, errors.New("无效的流程定义")
	}
	return definition.Nodes[idx].Data.StartNodeData.InputVariables, nil
}

func (t *TemplateService) GetNodePrototype(_ context.Context, nodeType model.NodeType) (string, error) {
	var prototype *model.Node
	switch nodeType {
	case model.NodeTypeCondition:
		prototype = model.ConditionNodePrototype
	case model.NodeTypeLLM:
		prototype = model.LLMNodePrototype
	case model.NodeTypeCrawler:
		prototype = model.CrawlerNodePrototype
	case model.NodeTypeKnowledgeRetrieval:
		prototype = model.KbRetrievalNodePrototype
	default:
		return "", errors.New("无效的节点类型")
	}
	data, _ := json.Marshal(prototype)
	return string(data), nil
}
