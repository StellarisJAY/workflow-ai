package workflow

import (
	"context"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"time"
)

func (e *Engine) executeQuestionOptimizeNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	nodeData *model.QuestionOptimizationNodeData, inputMap map[string]any) {
	q, ok := inputMap["question"]
	if !ok {
		panic("缺少question参数")
	}
	question := q.(string)
	llm, err := e.modelRepo.GetDetail(ctx, nodeData.ModelId)
	if err != nil {
		panic(err)
	}
	output, err := executeLLMTask(llm, model.QuestionOptimizationPrompt, "TEXT", map[string]interface{}{
		"question": question,
	})
	if err != nil {
		panic(err)
	}
	nodeInstance.Output = fmt.Sprintf("{\"result\": \"%s\"}", output)
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	nodeInstance.CompleteTime = time.Now()
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		panic(err)
	}
}
