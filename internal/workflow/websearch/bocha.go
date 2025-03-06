package websearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/config"
	"io"
	"net/http"
	"time"
)

type bochaAnswer struct {
	Type         string `json:"_type"`
	QueryContext struct {
		OriginalQuery string `json:"originalQuery"`
	} `json:"queryContext"`
	WebPages struct {
		WebSearchURL          string `json:"webSearchUrl"`
		TotalEstimatedMatches int    `json:"totalEstimatedMatches"`
		Value                 []struct {
			ID               string    `json:"id"`
			Name             string    `json:"name"`
			URL              string    `json:"url"`
			IsFamilyFriendly bool      `json:"isFamilyFriendly"`
			DisplayURL       string    `json:"displayUrl"`
			Snippet          string    `json:"snippet"`
			DateLastCrawled  time.Time `json:"dateLastCrawled"`
			SearchTags       []struct {
				Name    string `json:"name"`
				Content string `json:"content"`
			} `json:"searchTags,omitempty"`
			About []struct {
				Name string `json:"name"`
			} `json:"about,omitempty"`
		} `json:"value"`
	} `json:"webPages"`
	RelatedSearches struct {
		ID    string `json:"id"`
		Value []struct {
			Text         string `json:"text"`
			DisplayText  string `json:"displayText"`
			WebSearchURL string `json:"webSearchUrl"`
		} `json:"value"`
	} `json:"relatedSearches"`
	RankingResponse struct {
		Mainline struct {
			Items []struct {
				AnswerType  string `json:"answerType"`
				ResultIndex int    `json:"resultIndex"`
				Value       struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"mainline"`
		Sidebar struct {
			Items []struct {
				AnswerType string `json:"answerType"`
				Value      struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"sidebar"`
	} `json:"rankingResponse"`
}

type bochaHttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    bochaAnswer `json:"data"`
}

type bochaWebSearch struct {
	apiKey string
}

func newBochaWebSearch(config *config.Config) SearchProvider {
	return &bochaWebSearch{apiKey: config.BochaAPIKey}
}

func (b *bochaWebSearch) doSearch(query string, count int) (*bochaAnswer, error) {
	const endpoint = "https://api.bochaai.com/v1/web-search"
	body := struct {
		Query     string `json:"query"`
		Count     int    `json:"count"`
		Freshness string `json:"freshness"`
	}{
		Query:     query,
		Count:     count,
		Freshness: "nolimit",
	}
	data, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res bochaHttpResponse
	if err := json.Unmarshal(respData, &res); err != nil {
		return nil, err
	}
	return &res.Data, nil
}

func (b *bochaWebSearch) Search(query string, topN int) ([]SearchResult, error) {
	answer, err := b.doSearch(query, topN)
	if err != nil {
		return nil, err
	}
	results := make([]SearchResult, len(answer.WebPages.Value))
	for i, res := range answer.WebPages.Value {
		results[i] = SearchResult{
			Title:   res.Name,
			URL:     res.URL,
			Snippet: res.Snippet,
		}
	}
	return results, nil
}
