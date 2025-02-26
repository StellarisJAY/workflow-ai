package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"strings"
	"time"
)

func (e *Engine) executeConditionNode(ctx context.Context, node *model.Node, nodeData *model.ConditionNodeData,
	nodeInstance *model.NodeInstance) error {
	branches := nodeData.Branches
	// 获取流程定义
	data, _ := e.instanceRepo.GetWorkflowDefinition(ctx, nodeInstance.WorkflowId)
	if data == "" {
		return errors.New("workflow definition not found")
	}
	var definition model.WorkflowDefinition
	if err := json.Unmarshal([]byte(data), &definition); err != nil {
		return errors.New("invalid workflow definition")
	}

	// 找到第一个满足条件的分支
	var targetBranch *model.ConditionNodeBranch
	for idx, branch := range branches {
		// else 分支
		if idx == len(branches)-1 {
			targetBranch = branch
			break
		}
		// if和else if分支
		ok, err := e.evaluateConditions(ctx, branch.Conditions, branch.Connector, nodeInstance.WorkflowId, &definition)
		if err != nil {
			return err
		}
		if ok {
			targetBranch = branch
			break
		}
	}
	if targetBranch == nil {
		return errors.New("no possible condition found")
	}

	// 节点实例记录条件节点选择的分支id
	output := model.ConditionNodeOutput{SuccessBranch: targetBranch.Handle}
	outputData, _ := json.Marshal(output)
	nodeInstance.Output = string(outputData)
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	nodeInstance.CompleteTime = time.Now()
	err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance)
	if err != nil {
		return err
	}

	// 找到目标分支的后续节点
	nextNodes := FindBranchNextNodes(&definition, node, targetBranch)
	// 执行分支后续节点
	e.executeNextNodes(ctx, nextNodes, &definition, nodeInstance.WorkflowId)
	return nil
}

func (e *Engine) evaluateConditions(ctx context.Context, conditions []*model.Condition, connector string, workflowId int64,
	definition *model.WorkflowDefinition) (bool, error) {
	isAnd := connector == "and"
	for _, condition := range conditions {
		// 从变量实例表中获取变量值
		value1, _, err := e.getConditionVariableValue(ctx, condition.Value1, workflowId, definition)
		if err != nil {
			return false, err
		}
		value2, _, err := e.getConditionVariableValue(ctx, condition.Value2, workflowId, definition)
		if err != nil {
			return false, err
		}
		//if value1Type != value2Type {
		//	return false, nil
		//}
		success := false
		switch condition.Op {
		case "==":
			success = value1 == value2
		case "!=":
			success = value1 != value2
		case ">": // TODO 大小比较
		case "<":
		case ">=":
		case "<=":
		}
		// 或 连接只需要一个条件满足
		if success && !isAnd {
			return true, nil
		}
		// 与 一个条件不满足
		if !success && isAnd {
			return false, nil
		}
	}
	return isAnd, nil
}

func (e *Engine) getConditionVariableValue(ctx context.Context, variable *model.Variable, workflowId int64,
	definition *model.WorkflowDefinition) (string,
	model.VariableType, error) {
	varType := variable.Type
	if !variable.IsRef {
		return variable.Value, varType, nil
	}
	parts := strings.Split(variable.Ref, ".")
	if len(parts) != 2 {
		return "", varType, errors.New("invalid condition variable")
	}
	nodeId, varName := parts[0], parts[1]
	value, err := e.instanceRepo.GetOutputVariableFromNodeInstance(ctx, nodeId, workflowId, varName)
	if err != nil {
		return "", varType, err
	}
	originNode := FindNodeById(definition, nodeId)
	originVar := FindNodeOutputVariable(originNode, varName)
	value = strings.Trim(value, "\"")
	return value, originVar.Type, nil
}
