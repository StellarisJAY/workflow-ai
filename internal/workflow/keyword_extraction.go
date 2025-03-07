package workflow

import (
	"context"
	"encoding/json"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/prompts"
	"strings"
	"time"
)

func (e *Engine) executeKeywordExtractionNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	nodeData *model.KeywordExtractionNodeData, inputMap map[string]any) {
	q, ok := inputMap["question"]
	if !ok {
		panic("缺少question参数")
	}
	question := q.(string)
	llm, err := e.modelRepo.GetDetail(ctx, nodeData.ModelId)
	if err != nil {
		panic(err)
	}
	modelAPI, err := makeModelAPI(llm, "JSON")
	if err != nil {
		panic(err)
	}
	prompt := prompts.NewPromptTemplate(model.KeywordExtractionPrompt, []string{"question"})

	chain := chains.NewLLMChain(modelAPI, prompt)
	response, err := chain.Call(context.TODO(), map[string]any{
		"question": question,
	},
		chains.WithTemperature(0.2))
	if err != nil {
		panic(err)
	}
	// 将大模型输出结果写入节点实例，修改节点实例为完成状态
	output := response[chain.OutputKey].(string)
	output = strings.TrimPrefix(output, "```json")
	output = strings.TrimSuffix(output, "```")
	output = strings.TrimSpace(output)
	keywords := make([]string, 0)
	if err := json.Unmarshal([]byte(output), &keywords); err != nil {
		panic("大模型输出格式错误")
	}
	result := make(map[string]any)
	result["keywords"] = keywords
	result["total"] = len(keywords)
	data, _ := json.Marshal(result)
	nodeInstance.Output = string(data)
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	nodeInstance.CompleteTime = time.Now()
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		panic(err)
	}
}
