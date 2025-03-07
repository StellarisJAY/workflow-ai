package workflow

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/tmc/langchaingo/llms"
	"strconv"
	"strings"
	"time"
)

func (e *Engine) executeImageUnderstandingNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	nodeData *model.ImageUnderstandingNodeData,
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

	output, err := e.doImageUnderstandingTask(ctx, fileId, nodeData.Prompt, nodeData.OutputFormat, llm)
	if nodeData.OutputFormat == "JSON" {
		output = strings.TrimPrefix(output, "```json")
		output = strings.TrimSuffix(output, "```")
		output = strings.TrimSpace(output)
	} else {
		for _, variable := range nodeData.OutputVariables {
			if variable.Type == model.VariableTypeString {
				out := map[string]string{
					variable.Name: output,
				}
				data, _ := json.Marshal(out)
				output = string(data)
				break
			}
		}
	}
	nodeInstance.Output = output
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	nodeInstance.CompleteTime = time.Now()
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		panic(err)
	}
}

func (e *Engine) doImageUnderstandingTask(ctx context.Context, fileId int64, prompt string, outputFormat string,
	llm *model.ModelDetailDTO) (string, error) {
	file, err := e.fileRepo.Get(ctx, fileId)
	if err != nil {
		return "", err
	}
	data, err := e.fileStore.Download(ctx, file.Url)
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBuffer([]byte{})
	enc := base64.NewEncoder(base64.StdEncoding, buffer)
	defer enc.Close()
	if _, err = enc.Write(data); err != nil {
		return "", err
	}

	api, err := makeModelAPI(llm, outputFormat)
	if err != nil {
		return "", err
	}
	imageURL := fmt.Sprintf("data:image/%s;base64,%s", file.Type, buffer.String())
	messages := []llms.MessageContent{
		{Role: llms.ChatMessageTypeHuman, Parts: []llms.ContentPart{
			llms.TextContent{Text: prompt},
			llms.ImageURLContent{URL: imageURL},
		}},
	}
	content, err := api.GenerateContent(ctx, messages, llms.WithTemperature(0.2))
	if err != nil {
		return "", err
	}
	choice := content.Choices[0]
	return choice.Content, nil
}
