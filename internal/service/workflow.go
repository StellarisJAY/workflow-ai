package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/workflow"
)

type WorkflowService struct {
	templateRepo *repo.TemplateRepo
	engine       *workflow.Engine
	instanceRepo *repo.InstanceRepo
}

func NewWorkflowService(templateRepo *repo.TemplateRepo, engine *workflow.Engine,
	instanceRepo *repo.InstanceRepo) *WorkflowService {
	return &WorkflowService{
		templateRepo: templateRepo,
		engine:       engine,
		instanceRepo: instanceRepo,
	}
}

func (w *WorkflowService) Start(ctx context.Context, request *model.StartWorkflowRequest) (int64, error) {
	definition := request.Definition
	if request.Definition == "" {
		detail, err := w.templateRepo.GetDetail(ctx, request.TemplateId)
		if detail == nil {
			return 0, errors.New("template not found")
		}
		if err != nil {
			return 0, err
		}
		definition = detail.Data
	}
	return w.engine.Start(ctx, definition, request.TemplateId, 1, request.Inputs)
}

func (w *WorkflowService) Outputs(ctx context.Context, workflowId int64) ([]*model.NodeInstanceOutputDTO, error) {
	data, err := w.instanceRepo.GetWorkflowDefinition(ctx, workflowId)
	if err != nil {
		return nil, err
	}
	if data == "" {
		return nil, errors.New("workflow not found")
	}
	var def model.WorkflowDefinition
	_ = json.Unmarshal([]byte(data), &def)
	nodeOutputs, err := w.instanceRepo.GetNodeInstanceOutputs(ctx, workflowId)
	if err != nil {
		return nil, err
	}
	nodeMap := make(map[string]*model.Node)
	for _, node := range def.Nodes {
		nodeMap[node.Id] = node
	}
	for _, nodeOutput := range nodeOutputs {
		if node, ok := nodeMap[nodeOutput.NodeId]; ok {
			nodeOutput.NodeName = node.Name
			nodeOutput.Type = node.Type
		}
	}
	return nodeOutputs, nil
}
