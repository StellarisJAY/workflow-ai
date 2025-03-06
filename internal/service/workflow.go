package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/workflow"
	"strings"
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

func (w *WorkflowService) ListWorkflowInstance(ctx context.Context, query model.WorkflowInstanceQuery) ([]*model.WorkflowInstanceListDTO, int, error) {
	instanceList, total, err := w.instanceRepo.ListWorkflowInstance(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	for _, instance := range instanceList {
		instance.StatusName = instance.Status.String()
		instance.Duration = instance.CompleteTime.Sub(instance.AddTime).String()
	}
	return instanceList, total, nil
}

func (w *WorkflowService) GetWorkflowInstanceDetail(ctx context.Context, workflowId int64) (*model.WorkflowInstanceDetailDTO, error) {
	instance, err := w.instanceRepo.GetWorkflowInstanceDetail(ctx, workflowId)
	if err != nil {
		return nil, err
	}
	instance.StatusName = instance.Status.String()
	instance.Duration = instance.CompleteTime.Sub(instance.AddTime).String()
	var definition model.WorkflowDefinition
	_ = json.Unmarshal([]byte(instance.Data), &definition)
	nodeStatusList, err := w.instanceRepo.ListNodeInstanceStatus(ctx, workflowId)
	if err != nil {
		return nil, err
	}
	branches, err := w.instanceRepo.GetConditionNodeBranch(ctx, workflowId)
	if err != nil {
		return nil, err
	}
	for _, nodeStatus := range nodeStatusList {
		nodeStatus.StatusName = nodeStatus.Status.String()
	}
	for i := range branches {
		branches[i].Branch = strings.Trim(branches[i].Branch, "\"")
	}
	instance.NodeStatusList = nodeStatusList
	instance.PassedEdgesList = workflow.GetPassedEdges(&definition, nodeStatusList, branches)
	instance.SuccessBranchList = branches
	return instance, nil
}

func (w *WorkflowService) GetNodeInstance(ctx context.Context, workflowId int64, nodeId string) (*model.NodeInstanceDetailDTO, error) {
	instance, err := w.instanceRepo.GetNodeInstanceByNodeId(ctx, workflowId, nodeId)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, errors.New("node instance not found")
	}
	data, err := w.instanceRepo.GetWorkflowDefinition(ctx, workflowId)
	if err != nil {
		return nil, errors.New("workflow definition not found")
	}
	var definition model.WorkflowDefinition
	if err := json.Unmarshal([]byte(data), &definition); err != nil {
		return nil, errors.New("invalid workflow definition")
	}
	// 获取输出变量类型
	node := workflow.FindNodeById(&definition, nodeId)
	outputVarTypes := make(map[string]model.VariableType)
	if node != nil {
		outputVars := workflow.GetNodeOutputVariables(node)
		for _, outputVar := range outputVars {
			outputVarTypes[outputVar.Name] = outputVar.Type
		}
	}
	return &model.NodeInstanceDetailDTO{
		Id:                  instance.Id,
		NodeId:              instance.NodeId,
		Type:                instance.Type,
		AddTime:             instance.AddTime,
		CompleteTime:        instance.CompleteTime,
		Status:              instance.Status,
		StatusName:          instance.Status.String(),
		Output:              instance.Output,
		Error:               instance.Error,
		OutputVariableTypes: outputVarTypes,
	}, nil
}

func (w *WorkflowService) GetWorkflowTimeline(ctx context.Context, workflowId int64) ([]*model.WorkflowInstanceTimelineDTO, error) {
	timeline, err := w.instanceRepo.GetWorkflowTimeline(ctx, workflowId)
	if err != nil {
		return nil, err
	}
	for _, timeline := range timeline {
		timeline.StatusName = timeline.Status.String()
		timeline.Duration = timeline.CompleteTime.Sub(timeline.AddTime).String()
	}
	return timeline, nil
}
