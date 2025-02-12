package service

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/workflow"
)

type WorkflowService struct {
	templateRepo *repo.TemplateRepo
	engine       *workflow.Engine
}

func NewWorkflowService(templateRepo *repo.TemplateRepo, engine *workflow.Engine) *WorkflowService {
	return &WorkflowService{
		templateRepo: templateRepo,
		engine:       engine,
	}
}

func (w *WorkflowService) Start(ctx context.Context, request *model.StartWorkflowRequest) error {
	detail, err := w.templateRepo.GetDetail(ctx, request.TemplateId)
	if detail == nil {
		return errors.New("template not found")
	}
	if err != nil {
		return err
	}
	return w.engine.Start(ctx, detail, 1, request.Inputs)
}
