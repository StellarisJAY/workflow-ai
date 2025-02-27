package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"log"
	"strconv"
	"time"
)

func (e *Engine) executeKnowledgeRetrieveNode(ctx context.Context, node *model.Node,
	nodeData *model.RetrieveKnowledgeBaseNodeData, nodeInstance *model.NodeInstance) {
	defer func() {
		if err := recover(); err != nil {
			nodeInstance.Status = model.NodeInstanceStatusFailed
			nodeInstance.Error = err.(error).Error()
			if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
				log.Println("update kb node instance error:", err)
			}
		}
	}()
	inputMap, err := e.LookupInputVariables(ctx, nodeData.InputVariables, nodeInstance.WorkflowId)
	if err != nil {
		panic(err)
	}
	query, ok := inputMap["query"]
	if !ok {
		panic(errors.New("missing query variable"))
	}
	var result []*model.KbSearchReturnDocument
	switch nodeData.SearchType {
	case model.KbSearchTypeSimilarity:
		result, err = e.rag.SimilaritySearch(ctx, nodeData.KbId, query.(string), nodeData.SimilarityThreshold, nodeData.Count)
		// 相似度检索
	case model.KbSearchTypeFulltext:
		// 全文检索
		result, err = e.rag.FulltextSearch(ctx, nodeData.KbId, query.(string), nodeData.Count)
	default:
		panic(errors.New("unknown search type"))
	}
	if err != nil {
		panic(err)
	}
	// TODO 输出格式转换
	output := make(map[string]any)
	output["total"] = strconv.Itoa(len(result))
	documents := make([]string, len(result))
	for i, document := range result {
		documents[i] = document.Content
	}
	output["documents"] = documents
	data, _ := json.Marshal(output)
	nodeInstance.Output = string(data)
	nodeInstance.Status = model.NodeInstanceStatusCompleted
	nodeInstance.CompleteTime = time.Now()
	if err := e.instanceRepo.UpdateNodeInstance(ctx, nodeInstance); err != nil {
		panic(err)
	}
}
