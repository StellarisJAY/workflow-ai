package gte

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"github.com/StellrisJAY/workflow-ai/internal/model"
	"io"
	"net/http"
	"slices"
)

type Reranker struct {
	apiKey string
}

type requestBody struct {
	Model      string `json:"model"`
	Parameters struct {
		ReturnDocuments bool `json:"return_documents"`
		TopN            int  `json:"top_n"`
	}
	Input struct {
		Query     string   `json:"query"`
		Documents []string `json:"documents"`
	} `json:"input"`
}

type responseBody struct {
	Output struct {
		Results []result `json:"results"`
	} `json:"output"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

type result struct {
	Index          int     `json:"index"`
	RelevanceScore float32 `json:"relevance_score"`
}

func NewReranker(apiKey string) *Reranker {
	return &Reranker{apiKey: apiKey}
}

func (g *Reranker) Rerank(ctx context.Context, query string, documents []*model.KbSearchReturnDocument) ([]*model.KbSearchReturnDocument, error) {
	url := "https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank"
	method := "POST"

	docs := make([]string, len(documents))
	for i, doc := range documents {
		docs[i] = doc.Content
	}

	body := requestBody{
		Model: "gte-rerank",
		Parameters: struct {
			ReturnDocuments bool `json:"return_documents"`
			TopN            int  `json:"top_n"`
		}{false, len(documents)},
		Input: struct {
			Query     string   `json:"query"`
			Documents []string `json:"documents"`
		}{Query: query, Documents: docs},
	}

	payload, _ := json.Marshal(body)
	reader := bytes.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, reader)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer sk-0cb35933caeb43169eac4d93b8b27395")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var r responseBody
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	for _, doc := range r.Output.Results {
		documents[doc.Index].Score = doc.RelevanceScore
	}
	slices.SortFunc(documents, func(a, b *model.KbSearchReturnDocument) int {
		return cmp.Compare(b.Score, a.Score)
	})
	return documents, nil
}
