package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/rag"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/StellrisJAY/workflow-ai/internal/repo/fs"
	"github.com/bwmarrin/snowflake"
	"log"
	"slices"
	"time"
)

type Engine struct {
	instanceRepo *repo.InstanceRepo
	tm           *repo.TransactionManager
	snowflake    *snowflake.Node
	modelRepo    *repo.ModelRepo
	kbRepo       *repo.KnowledgeBaseRepo
	rag          *rag.DocumentProcessor
	conf         *config.Config
	fileRepo     *repo.FileRepo
	fileStore    fs.FileStore
}

func NewEngine(instanceRepo *repo.InstanceRepo, modelRepo *repo.ModelRepo, snowflake *snowflake.Node,
	tm *repo.TransactionManager, kbRepo *repo.KnowledgeBaseRepo, rag *rag.DocumentProcessor, conf *config.Config,
	fileRepo *repo.FileRepo, fileStore fs.FileStore) *Engine {
	return &Engine{
		instanceRepo: instanceRepo,
		tm:           tm,
		snowflake:    snowflake,
		modelRepo:    modelRepo,
		kbRepo:       kbRepo,
		rag:          rag,
		conf:         conf,
		fileRepo:     fileRepo,
		fileStore:    fileStore,
	}
}

func (e *Engine) Start(ctx context.Context, defJSON string, templateId int64, addUser int64,
	input map[string]any) (int64, error) {
	var definition model.WorkflowDefinition
	if err := json.Unmarshal([]byte(defJSON), &definition); err != nil {
		return 0, fmt.Errorf("invalid workflow definition")
	}
	idx := slices.IndexFunc(definition.Nodes, func(n *model.Node) bool { return n.Type == model.NodeTypeStart })
	if idx == -1 {
		return 0, fmt.Errorf("missing start node")
	}
	startNode := definition.Nodes[idx]
	startNodeData := startNode.Data.StartNodeData
	if startNodeData == nil {
		return 0, fmt.Errorf("invalid start node")
	}
	// 检查输入变量是否全部存在
	inputVariables := startNode.Data.Input
	for _, variable := range inputVariables {
		if _, ok := input[variable.Name]; !ok && variable.Required {
			return 0, fmt.Errorf("缺少必填变量: %s", variable.Name)
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

func (e *Engine) LookupInputVariables(ctx context.Context, variableDef []model.Input, workflowId int64) (map[string]any, error) {
	result := make(map[string]any)
	nodeInstancesCache := make(map[string]*model.NodeInstance)
	for _, variable := range variableDef {
		if variable.Value.Type == model.VarValueTypeLiteral {
			result[variable.Name] = variable.Value.Content
			continue
		}
		originNodeId, originVarName := variable.Value.SourceNode, variable.Value.SourceName
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
	}
	return result, nil
}

func (e *Engine) executeNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			nodeInstance.Status = model.NodeInstanceStatusFailed
			nodeInstance.CompleteTime = time.Now()
			nodeInstance.Error = err.Error()
			nodeInstance.Output = "{}"
			if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
				log.Println("update node instance failed", err)
			}
			e.UpdateWorkflowFailed(ctx, nodeInstance.WorkflowId)
		}
	}()
	inputMap, err := e.LookupInputVariables(ctx, node.Data.Input, nodeInstance.WorkflowId)
	if err != nil {
		panic(err)
	}
	switch node.Type {
	case model.NodeTypeLLM:
		llmNodeData := node.Data.LLMNodeData
		if llmNodeData == nil {
			panic(errors.New("invalid Model node data"))
		}
		e.executeLLMNode(context.TODO(), node, nodeInstance, llmNodeData, inputMap)
	case model.NodeTypeEnd:
		endNodeData := node.Data.EndNodeData
		if endNodeData == nil {
			panic(errors.New("invalid end node data"))
		}
		e.executeEndNode(context.TODO(), node, nodeInstance)
	case model.NodeTypeCrawler:
		crawlerNodeData := node.Data.CrawlerNodeData
		if crawlerNodeData == nil {
			panic(errors.New("invalid crawler node data"))
		}
		e.executeCrawlerNode(context.TODO(), node, nodeInstance, crawlerNodeData, inputMap)
	case model.NodeTypeCondition:
		conditionNodeData := node.Data.ConditionNodeData
		if conditionNodeData == nil {
			panic(errors.New("invalid condition node data"))
		}
		if err := e.executeConditionNode(context.TODO(), node, conditionNodeData, nodeInstance); err != nil {
			panic(err)
		}
	case model.NodeTypeKnowledgeRetrieval:
		kbRetrievalNodeData := node.Data.RetrieveKnowledgeBaseNodeData
		if kbRetrievalNodeData == nil {
			panic(errors.New("invalid knowledge base node data"))
		}
		e.executeKnowledgeRetrieveNode(context.TODO(), node, kbRetrievalNodeData, nodeInstance)
	case model.NodeTypeWebSearch:
		webSearchNodeData := node.Data.WebSearchNodeData
		if webSearchNodeData == nil {
			panic(errors.New("invalid webSearch Node Data"))
		}
		e.executeWebSearchNode(context.TODO(), node, nodeInstance, webSearchNodeData, inputMap)
	case model.NodeTypeKeywordExtraction:
		keywordExtractionNodeData := node.Data.KeywordExtractionNodeData
		if keywordExtractionNodeData == nil {
			panic(errors.New("invalid keyword extraction node data"))
		}
		e.executeKeywordExtractionNode(context.TODO(), node, nodeInstance, keywordExtractionNodeData, inputMap)
	case model.NodeTypeQuestionOptimization:
		llmTaskNodeData := node.Data.QuestionOptimizationNodeData
		if llmTaskNodeData == nil {
			panic(errors.New("invalid llm task node data"))
		}
		e.executeQuestionOptimizeNode(context.TODO(), node, nodeInstance, llmTaskNodeData, inputMap)
	case model.NodeTypeImageUnderstanding:
		nodeData := node.Data.ImageUnderstandingNodeData
		if nodeData == nil {
			panic(errors.New("invalid image understanding node data"))
		}
		e.executeImageUnderstandingNode(context.TODO(), node, nodeInstance, nodeData, inputMap)
	case model.NodeTypeOCR:
		nodeData := node.Data.OCRNodeData
		if nodeData == nil {
			panic(errors.New("invalid ocr node data"))
		}
		e.executeOCRNode(context.TODO(), node, nodeInstance, nodeData, inputMap)
	}
	if nodeInstance.Status != model.NodeInstanceStatusFailed {
		// 条件节点已经推进了流程，不需要再执行后续节点
		if nodeInstance.Type != model.NodeTypeCondition {
			e.stepWorkflow(context.TODO(), node, nodeInstance.WorkflowId)
		}
	} else {
		e.UpdateWorkflowFailed(context.TODO(), nodeInstance.WorkflowId)
	}
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
	if status, err := e.instanceRepo.GetWorkflowInstanceStatus(ctx, workflowId); err != nil {
		log.Println("get workflow instance status error:", err)
		return
	} else if status == model.WorkflowInstanceStatusCompleted || status == model.WorkflowInstanceStatusFailed {
		return
	}
	var executableNodes []struct {
		NodeInstance *model.NodeInstance
		Node         *model.Node
	}
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
		executableNodes = append(executableNodes, struct {
			NodeInstance *model.NodeInstance
			Node         *model.Node
		}{NodeInstance: nodeInstance, Node: next})
	}

	for _, executableNode := range executableNodes {
		go e.executeNode(context.Background(), executableNode.Node, executableNode.NodeInstance)
	}
}

func (e *Engine) executeEndNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance) {
	// 结束节点的输出与输入相同
	inputMap, err := e.LookupInputVariables(ctx, node.Data.Input, nodeInstance.WorkflowId)
	if err != nil {
		panic(err)
	}
	outputData, _ := json.Marshal(inputMap)
	nodeInstance.Output = string(outputData)
	nodeInstance.CompleteTime = time.Now()
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		panic(err)
	}
	if err := e.instanceRepo.UpdateWorkflowInstance(ctx, &model.WorkflowInstance{
		Id:           nodeInstance.WorkflowId,
		Status:       model.WorkflowInstanceStatusCompleted,
		CompleteTime: time.Now(),
	}); err != nil {
		panic(err)
	}
}
