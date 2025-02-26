package workflow

import (
	"context"
	"errors"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"log"
	"strings"
	"time"
)

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
			openai.WithToken(detail.ApiKey))
	case model.ApiTypeOllama:
		llm, err = ollama.New(ollama.WithServerURL(detail.BaseUrl),
			ollama.WithModel(detail.Code),
			ollama.WithFormat("json"))
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
	prompt := prompts.NewPromptTemplate(llmNodeData.Prompt, inputVariables)
	// 创建langchain，调用大模型API
	chain := chains.NewLLMChain(llm, prompt)
	response, err := chain.Call(context.TODO(), inputMap,
		chains.WithTemperature(llmNodeData.Temperature),
		chains.WithTopP(llmNodeData.TopP))
	if err != nil {
		panic(err)
	}
	// 将大模型输出结果写入节点实例，修改节点实例为完成状态
	output := response[chain.OutputKey].(string)
	// llm可能输出markdown格式，需要去除代码块前缀后缀
	if llmNodeData.OutputFormat == "JSON" {
		output = strings.TrimPrefix(output, "```json")
		output = strings.TrimSuffix(output, "```")
		output = strings.TrimSpace(output)
	} else if llmNodeData.OutputFormat == "TEXT" {
		// 文本格式输出，需要转换成与输出变量表对于的JSON格式
		for _, variable := range llmNodeData.OutputVariables {
			if variable.Type == model.VariableTypeString {
				output = fmt.Sprintf("{\"%s\":\"%s\"}", variable.Name, output)
				break
			}
		}
	}
	nodeInstance.Output = output
	nodeInstance.CompleteTime = time.Now()
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		log.Println("update node instance error:", err)
	}
}
