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
			nodeOutput.NodeName = node.Data.Name
			nodeOutput.Type = node.Type
		}
	}
	return nodeOutputs, nil
}

func (w *WorkflowService) ListWorkflowInstance(ctx context.Context) ([]*model.WorkflowInstanceListDTO, int, error) {
	instanceList, err := w.instanceRepo.ListWorkflowInstance(ctx)
	if err != nil {
		return nil, 0, err
	}
	for _, instance := range instanceList {
		instance.StatusName = instance.Status.String()
		instance.Duration = instance.CompleteTime.Sub(instance.AddTime).String()
	}
	return instanceList, len(instanceList), nil
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
	for i, branch := range branches {
		branches[i] = strings.Trim(branch, "\"")
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
	return &model.NodeInstanceDetailDTO{
		Id:           instance.Id,
		NodeId:       instance.NodeId,
		Type:         instance.Type,
		AddTime:      instance.AddTime,
		CompleteTime: instance.CompleteTime,
		Status:       instance.Status,
		StatusName:   instance.Status.String(),
		Output:       instance.Output,
		Error:        instance.Error,
	}, nil
}
