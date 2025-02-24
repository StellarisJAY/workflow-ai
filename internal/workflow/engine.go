package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/rag"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/bwmarrin/snowflake"
	"log"
	"slices"
	"strings"
	"time"
)

type Engine struct {
	instanceRepo *repo.InstanceRepo
	tm           *repo.TransactionManager
	snowflake    *snowflake.Node
	llmRepo      *repo.LLMRepo
	kbRepo       *repo.KnowledgeBaseRepo
	rag          *rag.DocumentProcessor
}

func NewEngine(instanceRepo *repo.InstanceRepo, llmRepo *repo.LLMRepo, snowflake *snowflake.Node,
	tm *repo.TransactionManager, kbRepo *repo.KnowledgeBaseRepo, rag *rag.DocumentProcessor) *Engine {
	return &Engine{
		instanceRepo: instanceRepo,
		tm:           tm,
		snowflake:    snowflake,
		llmRepo:      llmRepo,
		kbRepo:       kbRepo,
		rag:          rag,
	}
}

func (e *Engine) Start(ctx context.Context, defJSON string, templateId int64, addUser int64,
	input map[string]any) (int64, error) {
	var definition model.WorkflowDefinition
	if err := json.Unmarshal([]byte(defJSON), &definition); err != nil {
		return 0, fmt.Errorf("invalid workflow definition")
	}
	idx := slices.IndexFunc(definition.Nodes, func(n *model.Node) bool { return n.Type == string(model.NodeTypeStart) })
	if idx == -1 {
		return 0, fmt.Errorf("missing start node")
	}
	startNode := definition.Nodes[idx]
	startNodeData := startNode.Data.StartNodeData
	if startNodeData == nil {
		return 0, fmt.Errorf("invalid start node")
	}
	// 检查输入变量是否全部存在
	inputVariables := startNodeData.InputVariables
	for _, variable := range inputVariables {
		if _, ok := input[variable.Name]; !ok {
			return 0, fmt.Errorf("missing input variable: %s", variable.Name)
		}
	}
	instance := &model.WorkflowInstance{
		Id:           e.snowflake.Generate().Int64(),
		TemplateId:   templateId,
		Data:         defJSON,
		Status:       model.WorkflowInstanceStatusRunning,
		AddTime:      time.Now(),
		AddUser:      addUser,
		CompleteTime: time.Now(),
	}
	// 创建开始节点实例，把传入开始节点的参数用json保存
	inputJSON, _ := json.Marshal(input)
	startNodeInstance := &model.NodeInstance{
		Id:           e.snowflake.Generate().Int64(),
		NodeId:       startNode.Id,
		Status:       model.NodeInstanceStatusCompleted, // 开始节点实例始终为完成状态
		Output:       string(inputJSON),
		WorkflowId:   instance.Id,
		AddTime:      time.Now(),
		CompleteTime: time.Now(),
		Type:         startNode.Type,
	}

	err := e.tm.Tx(ctx, func(ctx context.Context) error {
		// 创建流程实例
		if err := e.instanceRepo.InsertWorkflowInstance(ctx, instance); err != nil {
			return fmt.Errorf("创建流程实例失败")
		}
		if err := e.instanceRepo.InsertNodeInstance(ctx, startNodeInstance); err != nil {
			return fmt.Errorf("创建开始节点实例失败")
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	e.stepWorkflow(ctx, startNode, instance.Id)
	return instance.Id, nil
}

func (e *Engine) LookupInputVariables(ctx context.Context, variableDef []*model.Variable, workflowId int64) (map[string]any, error) {
	result := make(map[string]any)
	nodeInstancesCache := make(map[string]*model.NodeInstance)
	for _, variable := range variableDef {
		if variable.Type == string(model.VariableTypeRef) {
			parts := strings.Split(variable.Value, ".")
			if len(parts) != 2 {
				return nil, errors.New("变量来源格式错误")
			}
			originNodeId, originVarName := parts[0], parts[1]
			var originNodeInstance *model.NodeInstance
			if n, ok := nodeInstancesCache[originNodeId]; ok {
				originNodeInstance = n
			} else {
				originNodeInstance, _ = e.instanceRepo.GetNodeInstanceByNodeId(ctx, workflowId, originNodeId)
				if originNodeInstance == nil {
					continue
				}
				nodeInstancesCache[originNodeId] = originNodeInstance
			}
			var inputMap map[string]any
			_ = json.Unmarshal([]byte(originNodeInstance.Output), &inputMap)
			if value, ok := inputMap[originVarName]; ok {
				result[variable.Name] = value
			}
		} else if variable.Type == string(model.VariableTypeString) {
			result[variable.Name] = variable.Value
		}
	}
	return result, nil
}

func (e *Engine) executeNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance) error {
	switch node.Type {
	case string(model.NodeTypeLLM):
		llmNodeData := node.Data.LLMNodeData
		if llmNodeData == nil {
			return fmt.Errorf("invalid LLM node data")
		}
		inputMap, err := e.LookupInputVariables(ctx, llmNodeData.InputVariables, nodeInstance.WorkflowId)
		if err != nil {
			return err
		}
		go func() {
			e.executeLLMNode(context.TODO(), node, nodeInstance, llmNodeData, inputMap)
			if nodeInstance.Status != model.NodeInstanceStatusFailed {
				e.stepWorkflow(context.TODO(), node, nodeInstance.WorkflowId)
			} else {
				e.UpdateWorkflowFailed(context.TODO(), nodeInstance.WorkflowId)
			}
		}()
	case string(model.NodeTypeEnd):
		endNodeData := node.Data.EndNodeData
		if endNodeData == nil {
			return fmt.Errorf("invalid end node data")
		}
		e.executeEndNode(context.TODO(), node, nodeInstance, endNodeData)
	case string(model.NodeTypeCrawler):
		crawlerNodeData := node.Data.CrawlerNodeData
		if crawlerNodeData == nil {
			return fmt.Errorf("invalid crawler node data")
		}
		inputMap, err := e.LookupInputVariables(ctx, crawlerNodeData.InputVariables, nodeInstance.WorkflowId)
		if err != nil {
			return err
		}
		go func() {
			e.executeCrawlerNode(context.TODO(), node, nodeInstance, crawlerNodeData, inputMap)
			if nodeInstance.Status != model.NodeInstanceStatusFailed {
				e.stepWorkflow(context.TODO(), node, nodeInstance.WorkflowId)
			} else {
				e.UpdateWorkflowFailed(context.TODO(), nodeInstance.WorkflowId)
			}
		}()
	case string(model.NodeTypeCondition):
		conditionNodeData := node.Data.ConditionNodeData
		if conditionNodeData == nil {
			return fmt.Errorf("invalid condition node data")
		}
		if err := e.executeConditionNode(context.TODO(), node, conditionNodeData, nodeInstance); err != nil {
			nodeInstance.Status = model.NodeInstanceStatusFailed
			nodeInstance.CompleteTime = time.Now()
			nodeInstance.Error = err.Error()
			if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
				log.Println("update node instance error:", err)
			}
			e.UpdateWorkflowFailed(ctx, nodeInstance.WorkflowId)
			return err
		}
	case string(model.NodeTypeKnowledgeRetrieval):
		kbRetrievalNodeData := node.Data.RetrieveKnowledgeBaseNodeData
		if kbRetrievalNodeData == nil {
			return fmt.Errorf("invalid knowledge base node data")
		}
		go func() {
			e.executeKnowledgeRetrieveNode(context.TODO(), node, kbRetrievalNodeData, nodeInstance)
			if nodeInstance.Status != model.NodeInstanceStatusFailed {
				e.stepWorkflow(context.TODO(), node, nodeInstance.WorkflowId)
			} else {
				e.UpdateWorkflowFailed(context.TODO(), nodeInstance.WorkflowId)
			}
		}()

	}
	return nil
}

func (e *Engine) UpdateWorkflowFailed(ctx context.Context, workflowId int64) {
	if err := e.instanceRepo.UpdateWorkflowInstance(ctx, &model.WorkflowInstance{
		Id:           workflowId,
		Status:       model.WorkflowInstanceStatusFailed,
		CompleteTime: time.Now(),
	}); err != nil {
		log.Println("update workflow instance error:", err)
	}
}

func (e *Engine) stepWorkflow(ctx context.Context, currNode *model.Node, workflowId int64) {
	instance, err := e.instanceRepo.GetWorkflowInstance(ctx, workflowId)
	if err != nil || instance == nil {
		log.Println("can't find flow instance", err, workflowId)
		return
	}
	var definition model.WorkflowDefinition
	_ = json.Unmarshal([]byte(instance.Data), &definition)
	nextNodes := GetNextNodes(&definition, currNode)
	e.executeNextNodes(ctx, nextNodes, &definition, workflowId)
}

func (e *Engine) executeNextNodes(ctx context.Context, nextNodes []*model.Node, definition *model.WorkflowDefinition,
	workflowId int64) {
	for _, next := range nextNodes {
		nodes := GetPrevNodes(definition, next)
		ids := make([]string, len(nodes))
		for i, node := range nodes {
			ids[i] = node.Id
		}
		count, err := e.instanceRepo.CountRunningNodeInstancesWithNodeIds(ctx, workflowId, ids)
		if err != nil {
			log.Println("count running node instance error", err)
			continue
		}
		if count > 0 {
			continue
		}
		nodeInstance := &model.NodeInstance{
			Id:           e.snowflake.Generate().Int64(),
			NodeId:       next.Id,
			Status:       model.NodeInstanceStatusRunning,
			WorkflowId:   workflowId,
			AddTime:      time.Now(),
			CompleteTime: time.Now(),
			Type:         next.Type,
			Output:       "{}",
			Error:        "",
		}
		if err := e.instanceRepo.InsertNodeInstance(ctx, nodeInstance); err != nil {
			log.Println("insert node instance error", err)
			continue
		}
		// 创建节点任务
		if err := e.executeNode(context.TODO(), next, nodeInstance); err != nil {
			log.Println("execute node instance error", err)
		}
	}
}

func (e *Engine) executeEndNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	endNodeData *model.EndNodeData) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			nodeInstance.Status = model.NodeInstanceStatusFailed
			nodeInstance.CompleteTime = time.Now()
			nodeInstance.Error = err.Error()
			if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
				log.Println("update node instance error:", err)
			}
		}
	}()
	outputVars := endNodeData.OutputVariables
	outputMap, err := e.LookupInputVariables(ctx, outputVars, nodeInstance.WorkflowId)
	if err != nil {
		log.Println("find output map error", err)
		panic(errors.New("无法获取output所需的变量"))
	}
	outputs, _ := json.Marshal(outputMap)
	nodeInstance.Output = string(outputs)
	nodeInstance.CompleteTime = time.Now()
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		log.Println("update node instance error:", err)
	}
	if err := e.instanceRepo.UpdateWorkflowInstance(ctx, &model.WorkflowInstance{
		Id:           nodeInstance.WorkflowId,
		Status:       model.WorkflowInstanceStatusCompleted,
		CompleteTime: time.Now(),
	}); err != nil {
		log.Println("update workflow instance error:", err)
	}
}
