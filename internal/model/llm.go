package model

import "time"

type LLMType string

const (
	LLMTypeChat      LLMType = "chat"
	LLMTypeEmbedding LLMType = "embedding"
	LLMlTypeReason   LLMType = "reason"
)

type ApiType string

const (
	ApiTypeOpenAI ApiType = "openai"
	ApiTypeOllama ApiType = "ollama"
)

type LLM struct {
	Id        int64     `json:"id,string"`
	Name      string    `json:"name"`
	ApiType   string    `json:"apiType"`
	BaseUrl   string    `json:"baseUrl"`
	ApiKey    string    `json:"apiKey"`
	Code      string    `json:"code"`
	AddUser   int64     `json:"addUser"`
	AddTime   time.Time `json:"addTime"`
	ModelType LLMType   `json:"modelType"`
}

type LLMDetailDTO struct {
	Id          int64     `json:"id,string"`
	Name        string    `json:"name"`
	ApiType     string    `json:"apiType"`
	BaseUrl     string    `json:"baseUrl"`
	ApiKey      string    `json:"apiKey"`
	Code        string    `json:"code"`
	AddUser     int64     `json:"addUser"`
	AddTime     time.Time `json:"addTime"`
	AddUsername string    `json:"addUsername"`
	ModelType   LLMType   `json:"modelType"`
}

type LLMQuery struct {
	Name      string  `json:"name"`
	ApiType   string  `json:"apiType"`
	Code      string  `json:"code"`
	ModelType LLMType `json:"modelType"`
}

type LLMListDTO struct {
	Id          int64     `json:"id,string"`
	Name        string    `json:"name"`
	ApiType     string    `json:"apiType"`
	Code        string    `json:"code"`
	AddUser     int64     `json:"addUser"`
	AddTime     time.Time `json:"addTime"`
	AddUsername string    `json:"addUsername"`
	ModelType   LLMType   `json:"modelType"`
}
