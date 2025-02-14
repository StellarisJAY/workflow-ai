package repo

import (
	"context"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"gorm.io/gorm"
)

type InstanceRepo struct {
	*Repository
}

func NewInstanceRepo(repo *Repository) *InstanceRepo {
	return &InstanceRepo{repo}
}

var (
	NodeInstanceTableName     = "wf_node_instance"
	WorkflowInstanceTableName = "wf_workflow_instance"
)

func (i *InstanceRepo) InsertNodeInstance(ctx context.Context, nodeInstance *model.NodeInstance) error {
	return i.DB(ctx).Table(NodeInstanceTableName).WithContext(ctx).Create(nodeInstance).Error
}

func (i *InstanceRepo) InsertWorkflowInstance(ctx context.Context, workflowInstance *model.WorkflowInstance) error {
	return i.DB(ctx).Table(WorkflowInstanceTableName).WithContext(ctx).Create(workflowInstance).Error
}

func (i *InstanceRepo) GetNodeInstance(ctx context.Context, id int64) (*model.NodeInstance, error) {
	var nodeInstance *model.NodeInstance
	err := i.DB(ctx).Table(NodeInstanceTableName).WithContext(ctx).Where("id = ?", id).Scan(&nodeInstance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nodeInstance, err
}

func (i *InstanceRepo) GetWorkflowInstance(ctx context.Context, id int64) (*model.WorkflowInstance, error) {
	var workflowInstance *model.WorkflowInstance
	err := i.DB(ctx).Table(WorkflowInstanceTableName).WithContext(ctx).Where("id =?", id).Scan(&workflowInstance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return workflowInstance, err
}

func (i *InstanceRepo) GetNodeInstanceByNodeId(ctx context.Context, workflowId int64, nodeId string) (*model.NodeInstance, error) {
	var nodeInstance *model.NodeInstance
	err := i.DB(ctx).Table(NodeInstanceTableName).
		WithContext(ctx).
		Where("workflow_id =? and node_id =?", workflowId, nodeId).
		Scan(&nodeInstance).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nodeInstance, err
}

func (i *InstanceRepo) UpdateNodeInstance(ctx context.Context, nodeInstance *model.NodeInstance) error {
	return i.DB(ctx).Table(NodeInstanceTableName).WithContext(ctx).
		Where("id=?", nodeInstance.Id).
		UpdateColumns(map[string]interface{}{
			"status":        nodeInstance.Status,
			"complete_time": nodeInstance.CompleteTime,
			"output":        nodeInstance.Output,
			"error":         nodeInstance.Error,
		}).Error
}

func (i *InstanceRepo) CountRunningNodeInstancesWithNodeIds(ctx context.Context, workflowId int64, nodeIds []string) (int64, error) {
	var count int64
	err := i.DB(ctx).Table(NodeInstanceTableName).
		WithContext(ctx).
		Where("workflow_id =?", workflowId).
		Where("node_id IN (?)", nodeIds).
		Where("status = ?", model.NodeInstanceStatusRunning).
		Count(&count).Error
	return count, err
}

func (i *InstanceRepo) UpdateWorkflowInstance(ctx context.Context, instance *model.WorkflowInstance) error {
	return i.DB(ctx).Table(WorkflowInstanceTableName).WithContext(ctx).
		Where("id=?", instance.Id).
		UpdateColumns(map[string]interface{}{
			"status":        instance.Status,
			"complete_time": instance.CompleteTime,
		}).Error
}

func (i *InstanceRepo) GetNodeInstanceOutputs(ctx context.Context, workflowId int64) ([]*model.NodeInstanceOutputDTO, error) {
	var outputs []*model.NodeInstanceOutputDTO
	err := i.DB(ctx).Table(NodeInstanceTableName).
		Select("id, node_id, output, error, status, add_time, complete_time").
		WithContext(ctx).
		Where("workflow_id =?", workflowId).
		Where("status != ?", model.NodeInstanceStatusRunning).
		Scan(&outputs).Error
	return outputs, err
}

func (i *InstanceRepo) GetWorkflowDefinition(ctx context.Context, workflowId int64) (string, error) {
	var definition string
	err := i.DB(ctx).Table(WorkflowInstanceTableName).
		Select("data").
		WithContext(ctx).
		Where("id =?", workflowId).
		Scan(&definition).Error
	return definition, err
}
