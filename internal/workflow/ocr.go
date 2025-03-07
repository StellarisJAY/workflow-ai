package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"strconv"
	"time"
)

func (e *Engine) executeOCRNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	nodeData *model.OCRNodeData,
	inputMap map[string]any) {
	img, ok := inputMap["image"]
	if !ok {
		panic(errors.New("缺少image参数"))
	}
	fileId, err := strconv.ParseInt(img.(string), 10, 64)
	if err != nil {
		panic(errors.New("image参数错误"))
	}

	llm, _ := e.modelRepo.GetDetail(ctx, nodeData.ModelId)
	if llm == nil || llm.ModelType != model.ModelTypeImageUnderstanding {
		panic(errors.New("模型不存在"))
	}

	output, err := e.doImageUnderstandingTask(ctx, fileId, model.OCRPrompt, "TEXT", llm)
	out := map[string]string{
		"document": output,
	}
	data, _ := json.Marshal(out)
	output = string(data)
	nodeInstance.Output = output
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	nodeInstance.CompleteTime = time.Now()
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		panic(err)
	}
}
