package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/repo"
	"github.com/bwmarrin/snowflake"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

type Engine struct {
	instanceRepo *repo.InstanceRepo
	tm           *repo.TransactionManager
	snowflake    *snowflake.Node
	llmRepo      *repo.LLMRepo
}

func NewEngine(instanceRepo *repo.InstanceRepo, llmRepo *repo.LLMRepo, snowflake *snowflake.Node,
	tm *repo.TransactionManager) *Engine {
	return &Engine{
		instanceRepo: instanceRepo,
		tm:           tm,
		snowflake:    snowflake,
		llmRepo:      llmRepo,
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
	e.createNextNode(ctx, startNode, startNodeInstance)
	return instance.Id, nil
}

//func (e *Engine) getInputMap(node *model.Node, workflowInstanceId int64) (map[string]any, error) {
//	var inputVariables map[string]string
//	switch model.NodeType(node.Type) {
//	case model.NodeTypeLLM:
//		inputVariables = node.Data.LLMNodeData.InputVariables
//	}
//	return e.findVariables(inputVariables, workflowInstanceId)
//}
//
//func (e *Engine) findVariables(variables map[string]string, workflowInstanceId int64) (map[string]any, error) {
//	result := make(map[string]any)
//	// 匹配{{来源变量名}}
//	regex := regexp.MustCompile("{{[^}]+}}")
//	for name, origin := range variables {
//		// 直接输入值的变量
//		if !regex.MatchString(origin) {
//			result[name] = origin
//			continue
//		}
//		// 需要从来源获取值的变量
//		origin = strings.TrimPrefix(origin, "{{")
//		origin = strings.TrimSuffix(origin, "}}")
//		parts := strings.Split(origin, ".")
//		if len(parts) != 2 {
//			return nil, errors.New("变量来源格式错误")
//		}
//		originNodeId, originVarName := parts[0], parts[1]
//		originNode, _ := e.instanceRepo.GetNodeInstanceByNodeId(context.Background(), workflowInstanceId, originNodeId)
//		if originNode == nil {
//			continue
//		}
//		originInputVars := make(map[string]any)
//		originOutputVars := make(map[string]any)
//		_ = json.Unmarshal([]byte(originNode.Input), &originInputVars)
//		_ = json.Unmarshal([]byte(originNode.Output), &originOutputVars)
//		if value, ok := originInputVars[originVarName]; ok {
//			result[name] = value
//			continue
//		}
//		if value, ok := originOutputVars[originVarName]; ok {
//			result[name] = value
//		}
//	}
//	return result, nil
//}

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
			e.executeLLMNode(ctx, node, nodeInstance, llmNodeData, inputMap)
			if nodeInstance.Status != model.NodeInstanceStatusFailed {
				e.createNextNode(ctx, node, nodeInstance)
			}
		}()
	case string(model.NodeTypeEnd):
		endNodeData := node.Data.EndNodeData
		if endNodeData == nil {
			return fmt.Errorf("invalid end node data")
		}
		e.executeEndNode(ctx, node, nodeInstance, endNodeData)

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
			e.executeCrawlerNode(ctx, node, nodeInstance, crawlerNodeData, inputMap)
			if nodeInstance.Status != model.NodeInstanceStatusFailed {
				e.createNextNode(ctx, node, nodeInstance)
			}
		}()
	}
	return nil
}

func (e *Engine) createNextNode(ctx context.Context, currNode *model.Node, currNodeInstance *model.NodeInstance) {
	instance, err := e.instanceRepo.GetWorkflowInstance(ctx, currNodeInstance.WorkflowId)
	if err != nil || instance == nil {
		log.Println("can't find flow instance", err, currNodeInstance.WorkflowId)
		return
	}
	var definition model.WorkflowDefinition
	_ = json.Unmarshal([]byte(instance.Data), &definition)
	nextNodes := GetNextNodes(&definition, currNode)
	for _, next := range nextNodes {
		nodes := GetPrevNodes(&definition, next)
		ids := make([]string, len(nodes))
		for i, node := range nodes {
			ids[i] = node.Id
		}
		count, err := e.instanceRepo.CountRunningNodeInstancesWithNodeIds(ctx, currNodeInstance.WorkflowId, ids)
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
			WorkflowId:   currNodeInstance.WorkflowId,
			AddTime:      time.Now(),
			CompleteTime: time.Now(),
			Type:         next.Type,
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

func (e *Engine) executeLLMNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	llmNodeData *model.LLMNodeData, inputMap map[string]any) {
	detail, _ := e.llmRepo.GetDetail(ctx, llmNodeData.ModelId)
	// 执行过程出错，将节点实例改为失败状态
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
	if detail == nil {
		panic(errors.New("无法找到节点需要的大模型"))
	}
	// 创建大模型接口
	var llm llms.Model
	var err error
	switch model.ApiType(detail.ApiType) {
	case model.ApiTypeOpenAI:
		llm, err = openai.New(openai.WithModel(detail.Code),
			openai.WithBaseURL(detail.BaseUrl),
			openai.WithToken(detail.ApiKey),
			openai.WithResponseFormat(openai.ResponseFormatJSON))
	case model.ApiTypeOllama:
		llm, err = ollama.New(ollama.WithServerURL(detail.BaseUrl),
			ollama.WithModel(detail.Code), ollama.WithFormat("json"))
	default:
		panic(errors.New("不支持的大模型类型"))
	}

	if err != nil {
		log.Println("create llm error:", err)
		panic(errors.New("创建大模型失败"))
	}
	// 创建提示词模板
	inputVariables := make([]string, 0, len(llmNodeData.InputVariables))
	for _, variable := range llmNodeData.InputVariables {
		inputVariables = append(inputVariables, variable.Name)
	}
	prompt := prompts.NewPromptTemplate(llmNodeData.UserPrompt, inputVariables)
	// 创建langchain，调用大模型API
	chain := chains.NewLLMChain(llm, prompt)
	response, err := chain.Call(context.TODO(), inputMap)
	if err != nil {
		panic(err)
	}
	// 将大模型输出结果写入节点实例，修改节点实例为完成状态
	output := response[chain.OutputKey].(string)
	nodeInstance.Output = output
	nodeInstance.CompleteTime = time.Now()
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		log.Println("update node instance error:", err)
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

func (e *Engine) executeCrawlerNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	crawlerData *model.CrawlerNodeData, inputMap map[string]any) {
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
	urlStr, ok := inputMap["url"]
	if !ok {
		panic(errors.New("url参数不存在"))
	}
	u, err := url.Parse(urlStr.(string))
	if err != nil {
		panic(err)
	}
	// 执行爬虫
	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		panic(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	if response.StatusCode != 200 {
		panic(errors.New("爬虫请求失败"))
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// 保存爬虫结果
	result := map[string]string{
		"content": string(bytes),
	}
	outputs, _ := json.Marshal(result)
	nodeInstance.Output = string(outputs)
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	nodeInstance.CompleteTime = time.Now()
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		log.Println("update node instance error:", err)
	}
}
