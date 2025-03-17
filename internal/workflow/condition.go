package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"log"
	"strconv"
	"strings"
	"time"
)

func (e *Engine) executeConditionNode(ctx context.Context, node *model.Node, nodeData *model.ConditionNodeData,
	nodeInstance *model.NodeInstance) error {
	defer func() {
		if err := recover(); err != nil {
			nodeInstance.Status = model.NodeInstanceStatusFailed
			nodeInstance.Error = err.(error).Error()
			if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
				log.Println("update kb node instance error:", err)
			}
		}
	}()
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
		value1, value1Type, err := e.getConditionVariableValue(ctx, condition.Value1, workflowId, definition)
		if err != nil {
			return false, err
		}
		value2, value2Type, err := e.getConditionVariableValue(ctx, condition.Value2, workflowId, definition)
		if err != nil {
			return false, err
		}
		if value1Type != value2Type {
			return false, fmt.Errorf("condition variable value type mismatch %s, %s", value1Type, value2Type)
		}
		success := false
		switch value1Type {
		case model.VariableTypeString:
			success, err = compareString(condition.Op, value1, value2)
		case model.VariableTypeNumber:
			success, err = compareNumber(condition.Op, value1, value2)
		case model.VariableTypeStringArray, model.VariableTypeNumberArray:
			success, err = compareArray(condition.Op, value1, value1Type)
		default:
			return false, fmt.Errorf("value type not supported: %s", value1Type)
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

func (e *Engine) getConditionVariableValue(ctx context.Context, variable *model.Input, workflowId int64,
	definition *model.WorkflowDefinition) (string,
	model.VariableType, error) {
	varType := variable.Type
	if variable.Value.Type == model.VarValueTypeLiteral {
		return variable.Value.Content, varType, nil
	}

	nodeId, varName := variable.Value.SourceNode, variable.Value.SourceName
	value, err := e.instanceRepo.GetOutputVariableFromNodeInstance(ctx, nodeId, workflowId, varName)
	if err != nil {
		return "", varType, err
	}
	originNode := FindNodeById(definition, nodeId)
	originVar := FindNodeOutputVariable(originNode, varName)
	value = strings.Trim(value, "\"")
	return value, originVar.Type, nil
}

func compareNumber(op string, value1, value2 string) (bool, error) {
	val1, err := strconv.ParseFloat(value1, 64)
	if err != nil {
		return false, fmt.Errorf("invalid number: %s", value1)
	}
	val2, err := strconv.ParseFloat(value2, 64)
	if err != nil {
		return false, fmt.Errorf("invalid number: %s", value2)
	}
	switch op {
	case "==":
		return val1 == val2, nil
	case "!=":
		return val1 != val2, nil
	case ">":
		return val1 > val2, nil
	case "<":
		return val1 < val2, nil
	case ">=":
		return val1 >= val2, nil
	case "<=":
		return val1 <= val2, nil
	default:
		return false, errors.New("invalid operator")
	}
}

func compareString(op string, value1, value2 string) (bool, error) {
	switch op {
	case "==":
		return value1 == value2, nil
	case "!=":
		return value1 != value2, nil
	case "contains":
		return strings.Contains(value1, value2), nil
	case "!contains":
		return !strings.Contains(value1, value2), nil
	case "empty":
		return len(value1) == 0, nil
	case "!empty":
		return len(value1) != 0, nil
	default:
		return false, errors.New("invalid operator")
	}
}

func compareArray(op string, value1 string, valType model.VariableType) (bool, error) {
	var array []any
	if err := json.Unmarshal([]byte(value1), &array); err != nil {
		return false, errors.New("invalid array")
	}
	switch op {
	case "empty":
		return len(array) == 0, nil
	case "!empty":
		return len(array) != 0, nil
	default:
		return false, errors.New("invalid operator")
	}
}
