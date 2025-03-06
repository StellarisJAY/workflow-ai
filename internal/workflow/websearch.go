package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"github.com/StellrisJAY/workflow-ai/internal/workflow/websearch"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func (e *Engine) executeWebSearchNode(ctx context.Context, node *model.Node, nodeInstance *model.NodeInstance,
	nodeData *model.WebSearchNodeData, inputMap map[string]any) {
	defer func() {
		if err := recover(); err != nil {
			nodeInstance.Status = model.NodeInstanceStatusFailed
			nodeInstance.CompleteTime = time.Now()
			nodeInstance.Error = err.(error).Error()
			if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
				log.Println("update node instance error:", err)
			}
		}
	}()
	var query string
	if q, ok := inputMap["query"]; !ok {
		panic(errors.New("缺少query参数"))
	} else {
		query = q.(string)
	}
	count := nodeData.TopN
	if count <= 0 {
		panic(errors.New("topN参数错误"))
	}
	searchProvider := websearch.CreateSearchProvider("bocha", e.conf)
	result, err := searchProvider.Search(query, count)
	if err != nil {
		panic(err)
	}
	contents := make([]string, len(result))
	urls := make([]string, len(result))
	for i, res := range result {
		urls[i] = res.URL
		response, err := httpRequest(res.URL, http.MethodGet)
		if err != nil {
			panic(err)
		}
		result := make(map[string]any)
		// 请求失败
		if response.StatusCode != 200 {
			contents[i] = ""
			continue
		}

		contentType := response.Header.Get("Content-Type")
		contentType = strings.Split(contentType, ";")[0]
		result["contentType"] = contentType
		switch contentType {
		case "application/json", "text/plain":
			bytes, err := io.ReadAll(response.Body)
			if err != nil {
				continue
			}
			contents[i] = string(bytes)
		case "text/html":
			content, err := parseHTMLContent(response)
			if err != nil {
				continue
			}
			contents[i] = content
		default:
			contents[i] = ""
		}
	}
	outputMap := make(map[string]any)
	outputMap["contents"] = contents
	outputMap["total"] = len(result)
	outputMap["urls"] = urls
	data, _ := json.Marshal(outputMap)
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	nodeInstance.CompleteTime = time.Now()
	nodeInstance.Output = string(data)
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		panic(err)
	}
}
