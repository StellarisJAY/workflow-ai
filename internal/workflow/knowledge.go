package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"strconv"
	"time"
)

func (e *Engine) executeKnowledgeRetrieveNode(ctx context.Context, node *model.Node,
	nodeData *model.RetrieveKnowledgeBaseNodeData, nodeInstance *model.NodeInstance) {
	inputMap, err := e.LookupInputVariables(ctx, node.Data.Input, nodeInstance.WorkflowId)
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
	case model.KbSearchTypeHybrid:
		// 混合检索
		result, err = e.rag.HybridSearch(ctx, nodeData.KbId, query.(string), nodeData.Count, nodeData.SimilarityThreshold,
			nodeData.DenseWeight, nodeData.SparseWeight)
	default:
		panic(errors.New("unknown search type"))
	}
	if err != nil {
		panic(err)
	}
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
