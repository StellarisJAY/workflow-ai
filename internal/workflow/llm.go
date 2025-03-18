package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/ai"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/prompts"
	"log"
	"strings"
	"time"
)

func (e *Engine) executeLLMNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	llmNodeData *model.LLMNodeData, inputMap map[string]any) {
	detail, _ := e.modelRepo.GetProviderModelDetail(ctx, llmNodeData.ModelId)
	if detail == nil {
		panic(errors.New("无法找到节点需要的大模型"))
	}
	// 创建大模型接口
	llm, err := ai.MakeModelInterface(detail, llmNodeData.OutputFormat)
	if err != nil {
		log.Println("create llm error:", err)
		panic(errors.New("创建大模型失败"))
	}
	// 创建提示词模板
	inputVariables := make([]string, 0, len(node.Data.Input))
	for _, variable := range node.Data.Input {
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
	}
	outputMap := make(map[string]any)
	outputMap["text"] = output
	outData, _ := json.Marshal(outputMap)
	nodeInstance.Output = string(outData)
	nodeInstance.CompleteTime = time.Now()
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		panic(fmt.Errorf("llm output format error: %v", err))
	}
}

func executeLLMTask(detail *model.ProviderModelDetail, promptTemplate string, outputFormat string, inputMap map[string]any) (string, error) {
	modelAPI, err := ai.MakeModelInterface(detail, outputFormat)
	if err != nil {
		return "", err
	}
	inputVariables := make([]string, len(inputMap))
	for k := range inputMap {
		inputVariables = append(inputVariables, k)
	}
	prompt := prompts.NewPromptTemplate(promptTemplate, inputVariables)

	chain := chains.NewLLMChain(modelAPI, prompt)
	response, err := chain.Call(context.TODO(), inputMap,
		chains.WithTemperature(0.2))
	if err != nil {
		return "", err
	}
	output := response[chain.OutputKey].(string)
	if outputFormat == "JSON" {
		output = strings.TrimPrefix(output, "```json")
		output = strings.TrimSuffix(output, "```")
		output = strings.TrimSpace(output)
	}
	return output, nil
}
