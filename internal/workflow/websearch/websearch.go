package websearch

import (
	"errors"
	"github.com/StellrisJAY/workflow-ai/internal/config"
)

type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type SearchProvider interface {
	Search(query string, topN int) ([]SearchResult, error)
}

func CreateSearchProvider(provider string, config *config.Config) SearchProvider {
	switch provider {
	case "bocha":
		return newBochaWebSearch(config)
	default:
		panic(errors.New("Unsupported provider: " + provider))
	}
}
